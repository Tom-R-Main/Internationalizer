package formats

import (
	"bytes"
	"fmt"
	"sort"
	"strings"
	"unicode"

	"github.com/Tom-R-Main/Internationalizer/internal/fluentpattern"
)

// FluentFormat is a lossless adapter for Fluent Translation List resources.
// It exposes message values, terms, and attributes as independent semantic
// units while retaining comments and resource ordering in the source file.
type FluentFormat struct{}

func (f *FluentFormat) Name() string         { return "fluent" }
func (f *FluentFormat) Extensions() []string { return []string{".ftl"} }

func (f *FluentFormat) Parse(data []byte) (map[string]string, error) {
	units, err := f.ParseUnits(data)
	if err != nil {
		return nil, err
	}
	return UnitValues(units), nil
}

func (f *FluentFormat) Serialize(entries map[string]string, original []byte) ([]byte, error) {
	units := make([]Unit, 0, len(entries))
	keys := make([]string, 0, len(entries))
	for key := range entries {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		kind := UnitMessage
		if strings.HasPrefix(key, "-") {
			kind = UnitTerm
		}
		if strings.Contains(key, ".") {
			kind = UnitAttribute
		}
		units = append(units, Unit{ID: key, Value: entries[key], Kind: kind})
	}
	return f.SerializeUnits(units, original)
}

func (f *FluentFormat) ParseUnits(data []byte) ([]Unit, error) {
	document, err := parseFluentDocument(data)
	if err != nil {
		return nil, err
	}
	units := make([]Unit, 0, len(document.assignments))
	for _, assignment := range document.assignments {
		if assignment.value == "" {
			continue
		}
		if _, err := fluentpattern.Analyze(assignment.value); err != nil {
			return nil, fmt.Errorf("fluent parse: pattern %q: %w", assignment.id, err)
		}
		units = append(units, Unit{
			ID:        assignment.id,
			Value:     assignment.value,
			Kind:      assignment.kind,
			Context:   assignment.context,
			Structure: "fluent-pattern-v1",
		})
	}
	return units, nil
}

func (f *FluentFormat) SerializeUnits(units []Unit, original []byte) ([]byte, error) {
	if err := ValidateUnits(units); err != nil {
		return nil, err
	}
	document, err := parseFluentDocument(original)
	if err != nil {
		return nil, err
	}

	edits := make([]fluentEdit, 0)
	missing := make([]Unit, 0)
	for _, unit := range units {
		assignment, ok := document.byID[unit.ID]
		if !ok {
			missing = append(missing, unit)
			continue
		}
		if assignment.value == unit.Value {
			continue
		}
		edits = append(edits, fluentEdit{
			start: assignment.start,
			end:   assignment.end,
			text:  renderFluentAssignment(assignment, unit.Value, document.newline),
		})
	}

	entryInsertions := make(map[string][]Unit)
	newEntryOrder := make([]string, 0)
	newEntries := make(map[string][]Unit)
	for _, unit := range missing {
		parent, _, attribute := fluentUnitIdentity(unit)
		if attribute {
			if _, exists := document.entries[parent]; exists {
				entryInsertions[parent] = append(entryInsertions[parent], unit)
				continue
			}
		} else if entry, exists := document.entries[parent]; exists {
			// The resource has an attribute-only message. Its empty value
			// assignment is addressable for serialization but not translation.
			assignment := entry.main
			edits = append(edits, fluentEdit{
				start: assignment.start,
				end:   assignment.end,
				text:  renderFluentAssignment(assignment, unit.Value, document.newline),
			})
			continue
		}
		if _, exists := newEntries[parent]; !exists {
			newEntryOrder = append(newEntryOrder, parent)
		}
		newEntries[parent] = append(newEntries[parent], unit)
	}

	for parent, inserted := range entryInsertions {
		entry := document.entries[parent]
		var text strings.Builder
		if entry.end > 0 && original[entry.end-1] != '\n' {
			text.WriteString(document.newline)
		}
		for _, unit := range inserted {
			_, attribute, _ := fluentUnitIdentity(unit)
			text.WriteString(renderNewFluentAssignment(entry.attributeIndent+"."+attribute+" =", unit.Value, entry.attributeIndent+"    ", document.newline))
		}
		edits = append(edits, fluentEdit{start: entry.end, end: entry.end, text: text.String()})
	}

	if len(newEntryOrder) > 0 {
		var appended strings.Builder
		if len(original) > 0 {
			if original[len(original)-1] != '\n' {
				appended.WriteString(document.newline)
			}
			if !bytes.HasSuffix(original, []byte(document.newline+document.newline)) {
				appended.WriteString(document.newline)
			}
		}
		for index, parent := range newEntryOrder {
			if index > 0 {
				appended.WriteString(document.newline)
			}
			group := newEntries[parent]
			context := ""
			for _, unit := range group {
				if unit.Context != "" {
					context = unit.Context
					break
				}
			}
			if context != "" {
				appended.WriteString(strings.TrimRight(context, "\r\n"))
				appended.WriteString(document.newline)
			}
			var main *Unit
			attributes := make([]Unit, 0)
			for unitIndex := range group {
				_, _, attribute := fluentUnitIdentity(group[unitIndex])
				if attribute {
					attributes = append(attributes, group[unitIndex])
				} else {
					main = &group[unitIndex]
				}
			}
			if main != nil {
				appended.WriteString(renderNewFluentAssignment(parent+" =", main.Value, "    ", document.newline))
			} else {
				appended.WriteString(parent + " =" + document.newline)
			}
			for _, attributeUnit := range attributes {
				_, attribute, _ := fluentUnitIdentity(attributeUnit)
				appended.WriteString(renderNewFluentAssignment("    ."+attribute+" =", attributeUnit.Value, "        ", document.newline))
			}
		}
		edits = append(edits, fluentEdit{start: len(original), end: len(original), text: appended.String()})
	}

	sort.SliceStable(edits, func(left, right int) bool {
		if edits[left].start != edits[right].start {
			return edits[left].start > edits[right].start
		}
		return edits[left].end > edits[right].end
	})
	result := append([]byte(nil), original...)
	for _, edit := range edits {
		result = append(result[:edit.start], append([]byte(edit.text), result[edit.end:]...)...)
	}
	return result, nil
}

type fluentDocument struct {
	assignments []fluentAssignment
	byID        map[string]fluentAssignment
	entries     map[string]fluentEntry
	newline     string
}

type fluentEntry struct {
	id              string
	end             int
	attributeIndent string
	main            fluentAssignment
}

type fluentAssignment struct {
	id        string
	parent    string
	attribute string
	kind      UnitKind
	context   string
	value     string
	header    string
	indent    string
	start     int
	end       int
}

type fluentLine struct {
	start   int
	end     int
	content string
}

type fluentEdit struct {
	start int
	end   int
	text  string
}

func parseFluentDocument(data []byte) (fluentDocument, error) {
	document := fluentDocument{
		byID:    make(map[string]fluentAssignment),
		entries: make(map[string]fluentEntry),
		newline: "\n",
	}
	if bytes.Contains(data, []byte("\r\n")) {
		document.newline = "\r\n"
	}
	lines := splitFluentLines(data)
	for index := 0; index < len(lines); {
		parent, equals, ok := parseFluentMessageHeader(lines[index].content)
		if !ok {
			index++
			continue
		}
		if _, duplicate := document.entries[parent]; duplicate {
			return fluentDocument{}, fmt.Errorf("fluent parse: duplicate message or term %q", parent)
		}
		entryEndLine := index + 1
		for entryEndLine < len(lines) && isFluentContinuation(lines[entryEndLine].content) {
			entryEndLine++
		}
		context := fluentCommentContext(lines, index, document.newline)
		attributeStarts := make([]int, 0)
		for candidate := index + 1; candidate < entryEndLine; candidate++ {
			if _, _, _, attribute := parseFluentAttributeHeader(lines[candidate].content); attribute {
				attributeStarts = append(attributeStarts, candidate)
			}
		}
		boundaries := append([]int{index}, attributeStarts...)
		boundaries = append(boundaries, entryEndLine)
		entry := fluentEntry{id: parent, end: lines[entryEndLine-1].end, attributeIndent: "    "}
		for boundary := 0; boundary < len(boundaries)-1; boundary++ {
			startLine := boundaries[boundary]
			endLine := boundaries[boundary+1]
			assignmentParent := parent
			attribute := ""
			assignmentEquals := equals
			kind := UnitMessage
			if startLine != index {
				indent, name, attrEquals, _ := parseFluentAttributeHeader(lines[startLine].content)
				attribute = name
				assignmentEquals = attrEquals
				entry.attributeIndent = indent
				kind = UnitAttribute
			}
			if strings.HasPrefix(parent, "-") && attribute == "" {
				kind = UnitTerm
			}
			id := parent
			if attribute != "" {
				id += "." + attribute
			}
			if _, duplicate := document.byID[id]; duplicate {
				return fluentDocument{}, fmt.Errorf("fluent parse: duplicate translation unit %q", id)
			}
			assignment := buildFluentAssignment(lines, startLine, endLine, assignmentEquals, id, assignmentParent, attribute, kind, context)
			document.assignments = append(document.assignments, assignment)
			document.byID[id] = assignment
			if attribute == "" {
				entry.main = assignment
			}
		}
		document.entries[parent] = entry
		index = entryEndLine
	}
	return document, nil
}

func buildFluentAssignment(lines []fluentLine, startLine, endLine, equals int, id, parent, attribute string, kind UnitKind, context string) fluentAssignment {
	line := lines[startLine]
	header := line.content[:equals+1]
	first := strings.TrimLeftFunc(line.content[equals+1:], unicode.IsSpace)
	continuations := make([]string, 0, endLine-startLine-1)
	minimumIndent := -1
	for index := startLine + 1; index < endLine; index++ {
		content := lines[index].content
		if strings.TrimSpace(content) == "" {
			continuations = append(continuations, "")
			continue
		}
		indent := leadingWhitespaceLength(content)
		if indent > 0 && (minimumIndent < 0 || indent < minimumIndent) {
			minimumIndent = indent
		}
		continuations = append(continuations, content)
	}
	if minimumIndent < 0 {
		minimumIndent = 4
		if attribute != "" {
			minimumIndent = leadingWhitespaceLength(line.content) + 4
		}
	}
	for index, continuation := range continuations {
		if leadingWhitespaceLength(continuation) >= minimumIndent {
			continuations[index] = continuation[minimumIndent:]
		}
	}
	parts := continuations
	if first != "" {
		parts = append([]string{first}, continuations...)
	}
	value := strings.TrimRight(strings.Join(parts, "\n"), "\n")
	indent := strings.Repeat(" ", minimumIndent)
	if startLine+1 < endLine {
		candidate := lines[startLine+1].content
		if len(candidate) >= minimumIndent {
			indent = candidate[:minimumIndent]
		}
	}
	return fluentAssignment{
		id: id, parent: parent, attribute: attribute, kind: kind, context: context,
		value: value, header: header, indent: indent, start: line.start,
		end: lines[endLine-1].end,
	}
}

func renderFluentAssignment(assignment fluentAssignment, value, newline string) string {
	return renderNewFluentAssignment(assignment.header, value, assignment.indent, newline)
}

func renderNewFluentAssignment(header, value, indent, newline string) string {
	if !strings.Contains(value, "\n") && value != "" {
		return header + " " + value + newline
	}
	var result strings.Builder
	result.WriteString(header)
	result.WriteString(newline)
	for _, line := range strings.Split(value, "\n") {
		result.WriteString(indent)
		result.WriteString(line)
		result.WriteString(newline)
	}
	return result.String()
}

func splitFluentLines(data []byte) []fluentLine {
	if len(data) == 0 {
		return nil
	}
	lines := make([]fluentLine, 0, bytes.Count(data, []byte{'\n'})+1)
	for start := 0; start < len(data); {
		end := bytes.IndexByte(data[start:], '\n')
		if end < 0 {
			end = len(data)
		} else {
			end += start + 1
		}
		contentEnd := end
		if contentEnd > start && data[contentEnd-1] == '\n' {
			contentEnd--
		}
		if contentEnd > start && data[contentEnd-1] == '\r' {
			contentEnd--
		}
		lines = append(lines, fluentLine{start: start, end: end, content: string(data[start:contentEnd])})
		start = end
	}
	return lines
}

func parseFluentMessageHeader(line string) (string, int, bool) {
	if line == "" || unicode.IsSpace(rune(line[0])) {
		return "", 0, false
	}
	equals := strings.IndexByte(line, '=')
	if equals < 0 {
		return "", 0, false
	}
	name := strings.TrimSpace(line[:equals])
	if !validFluentIdentifier(name, true) {
		return "", 0, false
	}
	return name, equals, true
}

func parseFluentAttributeHeader(line string) (string, string, int, bool) {
	indentLength := leadingWhitespaceLength(line)
	if indentLength == 0 || indentLength >= len(line) || line[indentLength] != '.' {
		return "", "", 0, false
	}
	equals := strings.IndexByte(line[indentLength+1:], '=')
	if equals < 0 {
		return "", "", 0, false
	}
	equals += indentLength + 1
	name := strings.TrimSpace(line[indentLength+1 : equals])
	if !validFluentIdentifier(name, false) {
		return "", "", 0, false
	}
	return line[:indentLength], name, equals, true
}

func validFluentIdentifier(value string, allowTerm bool) bool {
	if allowTerm && strings.HasPrefix(value, "-") {
		value = value[1:]
	}
	if value == "" || !isASCIILetter(value[0]) {
		return false
	}
	for index := 1; index < len(value); index++ {
		if !isASCIILetter(value[index]) && !isASCIIDigit(value[index]) && value[index] != '-' && value[index] != '_' {
			return false
		}
	}
	return true
}

func isFluentContinuation(line string) bool {
	return line != "" && (unicode.IsSpace(rune(line[0])) || strings.TrimSpace(line) == "}")
}

func leadingWhitespaceLength(value string) int {
	for index, character := range value {
		if !unicode.IsSpace(character) {
			return index
		}
	}
	return len(value)
}

func fluentCommentContext(lines []fluentLine, entryIndex int, newline string) string {
	start := entryIndex
	for start > 0 {
		candidate := strings.TrimSpace(lines[start-1].content)
		if !strings.HasPrefix(candidate, "#") {
			break
		}
		start--
	}
	if start == entryIndex {
		return ""
	}
	comments := make([]string, 0, entryIndex-start)
	for index := start; index < entryIndex; index++ {
		comments = append(comments, lines[index].content)
	}
	return strings.Join(comments, newline)
}

func fluentUnitIdentity(unit Unit) (parent, attribute string, isAttribute bool) {
	if unit.Kind != UnitAttribute {
		return unit.ID, "", false
	}
	separator := strings.LastIndexByte(unit.ID, '.')
	if separator < 1 || separator == len(unit.ID)-1 {
		return unit.ID, "", false
	}
	return unit.ID[:separator], unit.ID[separator+1:], true
}

func isASCIILetter(value byte) bool {
	return value >= 'A' && value <= 'Z' || value >= 'a' && value <= 'z'
}
func isASCIIDigit(value byte) bool { return value >= '0' && value <= '9' }
