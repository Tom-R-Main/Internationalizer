// Package message parses and compares ICU MessageFormat messages.
package message

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	localeid "github.com/Tom-R-Main/Internationalizer/internal/locale"
)

// Code is a stable machine-readable message validation outcome.
type Code string

const (
	CodeSyntax                Code = "syntax"
	CodeArgumentMismatch      Code = "argument_mismatch"
	CodeArgumentTypeMismatch  Code = "argument_type_mismatch"
	CodeArgumentStyleMismatch Code = "argument_style_mismatch"
	CodeSelectorMismatch      Code = "selector_mismatch"
	CodeInvalidPluralCategory Code = "invalid_plural_category"
)

const maxNestingDepth = 100

// Issue describes one structural incompatibility between source and target.
type Issue struct {
	Code     Code
	Argument string
	Message  string
}

// ArgumentType identifies the ICU operation performed by an argument.
type ArgumentType string

const (
	ArgumentSimple        ArgumentType = ""
	ArgumentSelect        ArgumentType = "select"
	ArgumentPlural        ArgumentType = "plural"
	ArgumentSelectOrdinal ArgumentType = "selectordinal"
	ArgumentNumber        ArgumentType = "number"
	ArgumentDate          ArgumentType = "date"
	ArgumentTime          ArgumentType = "time"
)

type elementKind uint8

const (
	elementText elementKind = iota
	elementArgument
	elementPound
)

// Message is a parsed ICU message. Its String method emits a deterministic
// representation suitable for structural hashing.
type Message struct {
	elements []element
}

type element struct {
	kind     elementKind
	text     string
	argument *Argument
}

// Argument is one parsed ICU argument.
type Argument struct {
	Name    string
	Type    ArgumentType
	Style   string
	Offset  int
	Options []Option
}

// Option is one select or plural branch.
type Option struct {
	Selector string
	Message  *Message
}

// Parse parses an ICU MessageFormat message.
func Parse(input string) (*Message, error) {
	parser := parser{input: input}
	message, err := parser.parseMessage(false, false)
	if err != nil {
		return nil, err
	}
	if parser.position != len(input) {
		return nil, parser.errorf("unexpected trailing input")
	}
	return message, nil
}

// LooksLike reports whether input contains an ICU-shaped argument. It excludes
// the other interpolation syntaxes Internationalizer already supports.
func LooksLike(input string) bool {
	for index := 0; index < len(input); index++ {
		if input[index] != '{' || (index > 0 && (input[index-1] == '{' || input[index-1] == '%')) || (index+1 < len(input) && input[index+1] == '{') {
			continue
		}
		cursor := index + 1
		for cursor < len(input) {
			width := spaceWidth(input, cursor)
			if width == 0 {
				break
			}
			cursor += width
		}
		start := cursor
		for cursor < len(input) {
			width := identifierWidth(input, cursor)
			if width == 0 {
				break
			}
			cursor += width
		}
		if cursor == start {
			continue
		}
		for cursor < len(input) {
			width := spaceWidth(input, cursor)
			if width == 0 {
				break
			}
			cursor += width
		}
		if cursor >= len(input) || input[cursor] == '}' || input[cursor] == ',' {
			return true
		}
	}
	return false
}

// Compare reports structural incompatibilities between ICU-shaped source and
// target messages. Plain text and non-ICU interpolation are ignored.
func Compare(source, target, targetLocale string) []Issue {
	if !LooksLike(source) && !LooksLike(target) {
		return nil
	}
	sourceMessage, err := Parse(source)
	if err != nil {
		return []Issue{{Code: CodeSyntax, Message: fmt.Sprintf("source ICU message: %v", err)}}
	}
	targetMessage, err := Parse(target)
	if err != nil {
		return []Issue{{Code: CodeSyntax, Message: fmt.Sprintf("target ICU message: %v", err)}}
	}

	issues := compareMessages(sourceMessage, targetMessage, targetLocale, "")
	sort.Slice(issues, func(left, right int) bool {
		if issues[left].Argument != issues[right].Argument {
			return issues[left].Argument < issues[right].Argument
		}
		if issues[left].Code != issues[right].Code {
			return issues[left].Code < issues[right].Code
		}
		return issues[left].Message < issues[right].Message
	})
	return issues
}

func compareMessages(source, target *Message, targetLocale, parentArgument string) []Issue {
	sourceArguments := immediateArguments(source)
	targetArguments := immediateArguments(target)
	var issues []Issue
	if sourcePounds, targetPounds := poundCount(source), poundCount(target); sourcePounds != targetPounds {
		issues = append(issues, Issue{Code: CodeSelectorMismatch, Argument: parentArgument, Message: fmt.Sprintf("ICU plural branch contains %d pound placeholders in the source and %d in the target", sourcePounds, targetPounds)})
	}
	for name, sourceOccurrences := range sourceArguments {
		targetOccurrences, ok := targetArguments[name]
		if !ok {
			issues = append(issues, Issue{Code: CodeArgumentMismatch, Argument: name, Message: fmt.Sprintf("ICU argument %q is missing from the corresponding target branch", name)})
			continue
		}
		if len(sourceOccurrences) != len(targetOccurrences) {
			issues = append(issues, Issue{Code: CodeArgumentMismatch, Argument: name, Message: fmt.Sprintf("ICU argument %q occurs %d times in the source branch and %d times in the target branch", name, len(sourceOccurrences), len(targetOccurrences))})
		}
		if len(sourceOccurrences) == len(targetOccurrences) && hasPerfectArgumentMatching(sourceOccurrences, targetOccurrences, targetLocale) {
			continue
		}
		sourceOccurrences = sortedArguments(sourceOccurrences)
		targetOccurrences = sortedArguments(targetOccurrences)
		occurrences := min(len(sourceOccurrences), len(targetOccurrences))
		for index := range occurrences {
			issues = append(issues, compareArguments(sourceOccurrences[index], targetOccurrences[index], targetLocale)...)
		}
	}
	for name := range targetArguments {
		if _, ok := sourceArguments[name]; !ok {
			issues = append(issues, Issue{Code: CodeArgumentMismatch, Argument: name, Message: fmt.Sprintf("ICU argument %q is absent from the corresponding source branch", name)})
		}
	}
	return issues
}

func hasPerfectArgumentMatching(source, target []*Argument, targetLocale string) bool {
	matchedSource := make([]int, len(target))
	for index := range matchedSource {
		matchedSource[index] = -1
	}
	var match func(int, []bool) bool
	match = func(sourceIndex int, visited []bool) bool {
		for targetIndex := range target {
			if visited[targetIndex] || len(compareArguments(source[sourceIndex], target[targetIndex], targetLocale)) != 0 {
				continue
			}
			visited[targetIndex] = true
			if matchedSource[targetIndex] == -1 || match(matchedSource[targetIndex], visited) {
				matchedSource[targetIndex] = sourceIndex
				return true
			}
		}
		return false
	}
	for sourceIndex := range source {
		if !match(sourceIndex, make([]bool, len(target))) {
			return false
		}
	}
	return true
}

func sortedArguments(arguments []*Argument) []*Argument {
	result := append([]*Argument(nil), arguments...)
	sort.SliceStable(result, func(left, right int) bool {
		return argumentSortKey(result[left]) < argumentSortKey(result[right])
	})
	return result
}

func argumentSortKey(argument *Argument) string {
	options := make([]string, 0, len(argument.Options))
	for _, option := range argument.Options {
		options = append(options, option.Selector+"{"+messageSortKey(option.Message)+"}")
	}
	sort.Strings(options)
	return fmt.Sprintf("%s\x00%s\x00%020d\x00%s", argument.Type, argument.Style, argument.Offset, strings.Join(options, "\x00"))
}

func messageSortKey(message *Message) string {
	parts := make([]string, 0, len(message.elements))
	for _, element := range message.elements {
		switch element.kind {
		case elementArgument:
			parts = append(parts, element.argument.Name+"{"+argumentSortKey(element.argument)+"}")
		case elementPound:
			parts = append(parts, "#")
		}
	}
	sort.Strings(parts)
	return strings.Join(parts, "\x00")
}

func poundCount(message *Message) int {
	count := 0
	for _, element := range message.elements {
		if element.kind == elementPound {
			count++
		}
	}
	return count
}

func compareArguments(source, target *Argument, targetLocale string) []Issue {
	name := source.Name
	if source.Type != target.Type {
		return []Issue{{Code: CodeArgumentTypeMismatch, Argument: name, Message: fmt.Sprintf("ICU argument %q types differ: source %q, target %q", name, displayType(source.Type), displayType(target.Type))}}
	}
	var issues []Issue
	switch source.Type {
	case ArgumentNumber, ArgumentDate, ArgumentTime:
		if source.Style != target.Style {
			issues = append(issues, Issue{Code: CodeArgumentStyleMismatch, Argument: name, Message: fmt.Sprintf("ICU argument %q formatter styles differ", name)})
		}
	case ArgumentSelect:
		sourceOptions := optionsBySelector(source)
		targetOptions := optionsBySelector(target)
		if !sameOptionSet(sourceOptions, targetOptions) {
			issues = append(issues, Issue{Code: CodeSelectorMismatch, Argument: name, Message: fmt.Sprintf("ICU select argument %q selectors differ", name)})
		}
		for selector, sourceBranch := range sourceOptions {
			if targetBranch, ok := targetOptions[selector]; ok {
				issues = append(issues, compareMessages(sourceBranch, targetBranch, targetLocale, name)...)
			}
		}
	case ArgumentPlural, ArgumentSelectOrdinal:
		if source.Offset != target.Offset {
			issues = append(issues, Issue{Code: CodeSelectorMismatch, Argument: name, Message: fmt.Sprintf("ICU argument %q plural offsets differ", name)})
		}
		issues = append(issues, comparePluralSelectors(source, target, targetLocale)...)
		sourceOptions := optionsBySelector(source)
		targetOptions := optionsBySelector(target)
		for selector, sourceBranch := range sourceOptions {
			if targetBranch, ok := targetOptions[selector]; ok {
				issues = append(issues, compareMessages(sourceBranch, targetBranch, targetLocale, name)...)
			}
		}
		sourceOther := sourceOptions["other"]
		for selector, targetBranch := range targetOptions {
			if _, sourceHasSelector := sourceOptions[selector]; sourceHasSelector {
				continue
			}
			if isExactSelector(selector) {
				issues = append(issues, compareTargetArgumentSubset(sourceOther, targetBranch, targetLocale)...)
			} else {
				issues = append(issues, compareMessages(sourceOther, targetBranch, targetLocale, name)...)
			}
		}
	}
	return issues
}

func compareTargetArgumentSubset(source, target *Message, targetLocale string) []Issue {
	sourceArguments := immediateArguments(source)
	targetArguments := immediateArguments(target)
	var issues []Issue
	for name, targetOccurrences := range targetArguments {
		sourceOccurrences, ok := sourceArguments[name]
		if !ok {
			issues = append(issues, Issue{Code: CodeArgumentMismatch, Argument: name, Message: fmt.Sprintf("ICU exact-selector branch introduces argument %q that is absent from the source fallback branch", name)})
			continue
		}
		if len(targetOccurrences) > len(sourceOccurrences) || !hasArgumentSubsetMatching(sourceOccurrences, targetOccurrences, targetLocale) {
			issues = append(issues, Issue{Code: CodeArgumentMismatch, Argument: name, Message: fmt.Sprintf("ICU exact-selector branch introduces an incompatible occurrence of argument %q", name)})
		}
	}
	return issues
}

func hasArgumentSubsetMatching(source, target []*Argument, targetLocale string) bool {
	matchedTarget := make([]int, len(source))
	for index := range matchedTarget {
		matchedTarget[index] = -1
	}
	var match func(int, []bool) bool
	match = func(targetIndex int, visited []bool) bool {
		for sourceIndex := range source {
			if visited[sourceIndex] || len(compareArguments(source[sourceIndex], target[targetIndex], targetLocale)) != 0 {
				continue
			}
			visited[sourceIndex] = true
			if matchedTarget[sourceIndex] == -1 || match(matchedTarget[sourceIndex], visited) {
				matchedTarget[sourceIndex] = targetIndex
				return true
			}
		}
		return false
	}
	for targetIndex := range target {
		if !match(targetIndex, make([]bool, len(source))) {
			return false
		}
	}
	return true
}

func comparePluralSelectors(source, target *Argument, targetLocale string) []Issue {
	name := source.Name
	categories, err := localeid.CardinalCategories(targetLocale)
	if source.Type == ArgumentSelectOrdinal {
		categories, err = localeid.OrdinalCategories(targetLocale)
	}
	if err != nil {
		return []Issue{{Code: CodeInvalidPluralCategory, Argument: name, Message: fmt.Sprintf("cannot resolve plural categories for locale %q: %v", targetLocale, err)}}
	}
	allowed := make(map[string]struct{}, len(categories))
	for _, category := range categories {
		allowed[category] = struct{}{}
	}
	sourceSelectors := selectorSet(source)
	targetSelectors := selectorSet(target)
	var issues []Issue
	for selector := range targetSelectors {
		if isExactSelector(selector) {
			continue
		}
		if _, ok := allowed[selector]; !ok {
			issues = append(issues, Issue{Code: CodeInvalidPluralCategory, Argument: name, Message: fmt.Sprintf("ICU argument %q plural category %q is not used by locale %q", name, selector, targetLocale)})
		}
	}
	for selector := range sourceSelectors {
		_, targetHasSelector := targetSelectors[selector]
		if targetHasSelector {
			continue
		}
		_, applicableCategory := allowed[selector]
		if isExactSelector(selector) || applicableCategory {
			issues = append(issues, Issue{Code: CodeSelectorMismatch, Argument: name, Message: fmt.Sprintf("ICU argument %q target is missing selector %q", name, selector)})
		}
	}
	return issues
}

func immediateArguments(message *Message) map[string][]*Argument {
	result := make(map[string][]*Argument)
	for _, element := range message.elements {
		if element.kind == elementArgument {
			argument := element.argument
			result[argument.Name] = append(result[argument.Name], argument)
		}
	}
	return result
}

func optionsBySelector(argument *Argument) map[string]*Message {
	result := make(map[string]*Message, len(argument.Options))
	for _, option := range argument.Options {
		result[option.Selector] = option.Message
	}
	return result
}

func selectorSet(argument *Argument) map[string]struct{} {
	result := make(map[string]struct{}, len(argument.Options))
	for _, option := range argument.Options {
		result[option.Selector] = struct{}{}
	}
	return result
}

func sameOptionSet(left, right map[string]*Message) bool {
	if len(left) != len(right) {
		return false
	}
	for value := range left {
		if _, ok := right[value]; !ok {
			return false
		}
	}
	return true
}

func displayType(argumentType ArgumentType) string {
	if argumentType == ArgumentSimple {
		return "argument"
	}
	return string(argumentType)
}

func isExactSelector(selector string) bool {
	_, ok := normalizeExactSelector(selector)
	return ok
}

func normalizeExactSelector(selector string) (string, bool) {
	if !strings.HasPrefix(selector, "=") || len(selector) == 1 {
		return "", false
	}
	number := selector[1:]
	negative := false
	if number[0] == '-' {
		negative = true
		number = number[1:]
	}
	if number == "" {
		return "", false
	}
	dotSeen := false
	digitSeenAfterDot := false
	for _, character := range number {
		switch {
		case character >= '0' && character <= '9':
			if dotSeen {
				digitSeenAfterDot = true
			}
		case character == '.' && !dotSeen:
			dotSeen = true
		default:
			return "", false
		}
	}
	if dotSeen && !digitSeenAfterDot {
		return "", false
	}
	integer, fraction, _ := strings.Cut(number, ".")
	integer = strings.TrimLeft(integer, "0")
	if integer == "" {
		integer = "0"
	}
	fraction = strings.TrimRight(fraction, "0")
	if integer == "0" && fraction == "" {
		negative = false
	}
	canonical := "="
	if negative {
		canonical += "-"
	}
	canonical += integer
	if fraction != "" {
		canonical += "." + fraction
	}
	return canonical, true
}

type parser struct {
	input    string
	position int
	depth    int
}

func (p *parser) parseBranch(pluralContext bool) (*Message, error) {
	if p.depth >= maxNestingDepth {
		return nil, p.errorf("ICU message nesting exceeds %d levels", maxNestingDepth)
	}
	p.depth++
	defer func() { p.depth-- }()
	return p.parseMessage(true, pluralContext)
}

func (p *parser) parseMessage(nested, pluralContext bool) (*Message, error) {
	message := &Message{}
	var literal strings.Builder
	flush := func() {
		if literal.Len() == 0 {
			return
		}
		message.elements = append(message.elements, element{kind: elementText, text: literal.String()})
		literal.Reset()
	}

	for p.position < len(p.input) {
		switch current := p.input[p.position]; current {
		case '{':
			flush()
			argument, err := p.parseArgument(pluralContext)
			if err != nil {
				return nil, err
			}
			message.elements = append(message.elements, element{kind: elementArgument, argument: argument})
		case '}':
			if !nested {
				return nil, p.errorf("unexpected closing brace")
			}
			flush()
			return message, nil
		case '#':
			if pluralContext {
				flush()
				message.elements = append(message.elements, element{kind: elementPound})
				p.position++
				continue
			}
			literal.WriteByte(current)
			p.position++
		case '\'':
			quoted, err := p.parseApostrophe(pluralContext)
			if err != nil {
				return nil, err
			}
			literal.WriteString(quoted)
		default:
			literal.WriteByte(current)
			p.position++
		}
	}
	if nested {
		return nil, p.errorf("unterminated message branch")
	}
	flush()
	return message, nil
}

func (p *parser) parseApostrophe(pluralContext bool) (string, error) {
	p.position++
	if p.position < len(p.input) && p.input[p.position] == '\'' {
		p.position++
		return "'", nil
	}
	if p.position >= len(p.input) || !p.apostropheStartsQuote(pluralContext) {
		return "'", nil
	}
	var quoted strings.Builder
	for p.position < len(p.input) {
		if p.input[p.position] != '\'' {
			quoted.WriteByte(p.input[p.position])
			p.position++
			continue
		}
		if p.position+1 < len(p.input) && p.input[p.position+1] == '\'' {
			quoted.WriteByte('\'')
			p.position += 2
			continue
		}
		p.position++
		return quoted.String(), nil
	}
	// ICU MessageFormat treats the end of a pattern as the implicit end of an
	// apostrophe-quoted syntax span.
	return quoted.String(), nil
}

func (p *parser) apostropheStartsQuote(pluralContext bool) bool {
	switch p.input[p.position] {
	case '{', '}':
		return true
	case '#':
		return pluralContext
	default:
		return false
	}
}

func (p *parser) parseArgument(pluralContext bool) (*Argument, error) {
	p.position++
	p.skipSpace()
	name := p.readIdentifier()
	if name == "" {
		return nil, p.errorf("argument name is required")
	}
	p.skipSpace()
	if p.consume('}') {
		return &Argument{Name: name}, nil
	}
	if !p.consume(',') {
		return nil, p.errorf("expected comma or closing brace after argument %q", name)
	}
	p.skipSpace()
	argumentType := ArgumentType(strings.ToLower(p.readIdentifier()))
	switch argumentType {
	case ArgumentSelect, ArgumentPlural, ArgumentSelectOrdinal, ArgumentNumber, ArgumentDate, ArgumentTime:
	default:
		return nil, p.errorf("unsupported argument type %q", argumentType)
	}
	argument := &Argument{Name: name, Type: argumentType}
	p.skipSpace()
	if argumentType == ArgumentNumber || argumentType == ArgumentDate || argumentType == ArgumentTime {
		if p.consume(',') {
			styleStart := p.position
			for p.position < len(p.input) && p.input[p.position] != '}' {
				if p.input[p.position] == '{' {
					return nil, p.errorf("unexpected opening brace in %s style", argumentType)
				}
				p.position++
			}
			argument.Style = strings.TrimSpace(p.input[styleStart:p.position])
			if argument.Style == "" {
				return nil, p.errorf("%s style must not be empty", argumentType)
			}
		}
		if !p.consume('}') {
			return nil, p.errorf("unterminated %s argument %q", argumentType, name)
		}
		return argument, nil
	}

	if !p.consume(',') {
		return nil, p.errorf("expected comma before selectors for argument %q", name)
	}
	p.skipSpace()
	if argumentType == ArgumentPlural || argumentType == ArgumentSelectOrdinal {
		if strings.HasPrefix(p.input[p.position:], "offset:") {
			p.position += len("offset:")
			p.skipSpace()
			start := p.position
			for p.position < len(p.input) && p.input[p.position] >= '0' && p.input[p.position] <= '9' {
				p.position++
			}
			if start == p.position {
				return nil, p.errorf("plural offset must be a non-negative integer")
			}
			offset, err := strconv.Atoi(p.input[start:p.position])
			if err != nil {
				return nil, p.errorf("invalid plural offset")
			}
			argument.Offset = offset
			p.skipSpace()
		}
	}

	selectors := make(map[string]struct{})
	for p.position < len(p.input) && p.input[p.position] != '}' {
		selector := p.readSelector()
		if selector == "" {
			return nil, p.errorf("selector is required for argument %q", name)
		}
		if argumentType != ArgumentSelect {
			if canonical, ok := normalizeExactSelector(selector); ok {
				selector = canonical
			}
		}
		if _, duplicate := selectors[selector]; duplicate {
			return nil, p.errorf("duplicate selector %q for argument %q", selector, name)
		}
		if argumentType == ArgumentSelect {
			if !isIdentifier(selector) {
				return nil, p.errorf("invalid select keyword %q", selector)
			}
		} else if !isPluralSelector(selector) {
			return nil, p.errorf("invalid plural selector %q", selector)
		}
		selectors[selector] = struct{}{}
		p.skipSpace()
		if !p.consume('{') {
			return nil, p.errorf("expected message branch for selector %q", selector)
		}
		branchContext := pluralContext || argumentType == ArgumentPlural || argumentType == ArgumentSelectOrdinal
		branch, err := p.parseBranch(branchContext)
		if err != nil {
			return nil, err
		}
		if !p.consume('}') {
			return nil, p.errorf("unterminated selector %q", selector)
		}
		argument.Options = append(argument.Options, Option{Selector: selector, Message: branch})
		p.skipSpace()
	}
	if !p.consume('}') {
		return nil, p.errorf("unterminated argument %q", name)
	}
	if _, ok := selectors["other"]; !ok {
		return nil, p.errorf("argument %q requires an \"other\" selector", name)
	}
	return argument, nil
}

func (p *parser) readIdentifier() string {
	start := p.position
	for p.position < len(p.input) {
		width := identifierWidth(p.input, p.position)
		if width == 0 {
			break
		}
		p.position += width
	}
	return p.input[start:p.position]
}

func (p *parser) readSelector() string {
	p.skipSpace()
	start := p.position
	for p.position < len(p.input) && p.input[p.position] != '{' && p.input[p.position] != '}' {
		if spaceWidth(p.input, p.position) != 0 {
			break
		}
		_, width := utf8.DecodeRuneInString(p.input[p.position:])
		p.position += width
	}
	return p.input[start:p.position]
}

func (p *parser) skipSpace() {
	for p.position < len(p.input) {
		width := spaceWidth(p.input, p.position)
		if width == 0 {
			break
		}
		p.position += width
	}
}

func (p *parser) consume(expected byte) bool {
	if p.position >= len(p.input) || p.input[p.position] != expected {
		return false
	}
	p.position++
	return true
}

func (p *parser) errorf(format string, arguments ...any) error {
	return fmt.Errorf("at byte %d: %s", p.position, fmt.Sprintf(format, arguments...))
}

func identifierWidth(input string, position int) int {
	character, width := utf8.DecodeRuneInString(input[position:])
	if character == utf8.RuneError && width == 1 {
		return 0
	}
	if character < utf8.RuneSelf {
		value := byte(character)
		if value == '_' || value == '-' || value == '.' || value >= '0' && value <= '9' || value >= 'A' && value <= 'Z' || value >= 'a' && value <= 'z' {
			return width
		}
		return 0
	}
	if unicode.Is(unicode.Pattern_Syntax, character) || unicode.Is(unicode.Pattern_White_Space, character) {
		return 0
	}
	return width
}

func isIdentifier(input string) bool {
	if input == "" {
		return false
	}
	for position := 0; position < len(input); {
		width := identifierWidth(input, position)
		if width == 0 {
			return false
		}
		position += width
	}
	return true
}

func spaceWidth(input string, position int) int {
	character, width := utf8.DecodeRuneInString(input[position:])
	if character == utf8.RuneError && width == 1 || !unicode.IsSpace(character) {
		return 0
	}
	return width
}

func isPluralSelector(selector string) bool {
	if isExactSelector(selector) {
		return true
	}
	switch selector {
	case "zero", "one", "two", "few", "many", "other":
		return true
	default:
		return false
	}
}

// String emits a deterministic ICU representation.
func (message *Message) String() string {
	var output strings.Builder
	message.writeTo(&output, false)
	return output.String()
}

func (message *Message) writeTo(output *strings.Builder, pluralContext bool) {
	for _, element := range message.elements {
		switch element.kind {
		case elementText:
			writeLiteral(output, element.text, pluralContext)
		case elementPound:
			output.WriteByte('#')
		case elementArgument:
			element.argument.writeTo(output, pluralContext)
		}
	}
}

func (argument *Argument) writeTo(output *strings.Builder, pluralContext bool) {
	output.WriteByte('{')
	output.WriteString(argument.Name)
	if argument.Type == ArgumentSimple {
		output.WriteByte('}')
		return
	}
	output.WriteString(", ")
	output.WriteString(string(argument.Type))
	if argument.Type == ArgumentNumber || argument.Type == ArgumentDate || argument.Type == ArgumentTime {
		if argument.Style != "" {
			output.WriteString(", ")
			output.WriteString(argument.Style)
		}
		output.WriteByte('}')
		return
	}
	output.WriteString(", ")
	if argument.Offset != 0 {
		_, _ = fmt.Fprintf(output, "offset:%d ", argument.Offset)
	}
	options := append([]Option(nil), argument.Options...)
	sort.SliceStable(options, func(left, right int) bool { return options[left].Selector < options[right].Selector })
	for index, option := range options {
		if index > 0 {
			output.WriteByte(' ')
		}
		output.WriteString(option.Selector)
		output.WriteString(" {")
		branchContext := pluralContext || argument.Type == ArgumentPlural || argument.Type == ArgumentSelectOrdinal
		option.Message.writeTo(output, branchContext)
		output.WriteByte('}')
	}
	output.WriteByte('}')
}

func writeLiteral(output *strings.Builder, literal string, pluralContext bool) {
	characters := []rune(literal)
	firstQuoted := -1
	lastQuoted := -1
	for index, character := range characters {
		if literalNeedsQuote(character, pluralContext) {
			if firstQuoted < 0 {
				firstQuoted = index
			}
			lastQuoted = index
		}
	}
	if firstQuoted < 0 {
		writeUnquotedLiteral(output, characters)
		return
	}
	writeUnquotedLiteral(output, characters[:firstQuoted])
	output.WriteByte('\'')
	for _, character := range characters[firstQuoted : lastQuoted+1] {
		if character == '\'' {
			output.WriteString("''")
		} else {
			output.WriteRune(character)
		}
	}
	output.WriteByte('\'')
	writeUnquotedLiteral(output, characters[lastQuoted+1:])
}

func literalNeedsQuote(character rune, pluralContext bool) bool {
	return character == '{' || character == '}' || pluralContext && character == '#'
}

func writeUnquotedLiteral(output *strings.Builder, characters []rune) {
	for _, character := range characters {
		if character == '\'' {
			output.WriteString("''")
		} else {
			output.WriteRune(character)
		}
	}
}
