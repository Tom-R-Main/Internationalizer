package formats

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"
)

const markdownPreambleKey = "markdown:preamble"

var (
	markdownFencePattern  = regexp.MustCompile(`^[ \t]{0,3}(` + "`{3,}" + `|~{3,})`)
	markdownH2Pattern     = regexp.MustCompile(`^[ \t]{0,3}##[ \t]+(.+?)[ \t]*#*[ \t]*(?:\r?\n)?$`)
	markdownMarkerPattern = regexp.MustCompile(`^[ \t]*<!--[ \t]*internationalizer:unit[ \t]+([a-z0-9][a-z0-9:._-]*)[ \t]*-->[ \t]*(?:\r?\n)?$`)
)

// MarkdownFormat translates the document preamble and each level-two section
// independently. Paired serialization writes stable, invisible unit markers
// into targets so later source-section insertions do not remap existing prose.
type MarkdownFormat struct{}

func (f *MarkdownFormat) Name() string         { return "markdown" }
func (f *MarkdownFormat) Extensions() []string { return []string{".md", ".mdx"} }

func (f *MarkdownFormat) Parse(data []byte) (map[string]string, error) {
	sections, err := parseMarkdownSections(data)
	if err != nil {
		return nil, err
	}
	assignSourceKeys(sections)
	return markdownEntries(sections), nil
}

// ParseTarget maps a translated document to the source's stable section keys.
// Existing unmarked documents are accepted only when their section count still
// matches, which provides a safe one-time migration path.
func (f *MarkdownFormat) ParseTarget(source, target []byte) (map[string]string, error) {
	sourceSections, err := parseMarkdownSections(source)
	if err != nil {
		return nil, fmt.Errorf("parsing markdown source: %w", err)
	}
	assignSourceKeys(sourceSections)
	targetSections, err := parseMarkdownSections(target)
	if err != nil {
		return nil, fmt.Errorf("parsing markdown target: %w", err)
	}

	marked := 0
	for _, section := range targetSections[1:] {
		if section.marker != "" {
			marked++
		}
	}
	if marked != 0 && marked != len(targetSections)-1 {
		return nil, fmt.Errorf("markdown target mixes marked and unmarked sections")
	}
	if marked == 0 {
		if len(targetSections) != len(sourceSections) {
			return nil, fmt.Errorf("unmarked markdown target has %d sections; source has %d", len(targetSections)-1, len(sourceSections)-1)
		}
		for i := range targetSections {
			targetSections[i].key = sourceSections[i].key
		}
	} else {
		targetSections[0].key = markdownPreambleKey
		seen := map[string]struct{}{markdownPreambleKey: {}}
		for _, section := range targetSections[1:] {
			if _, duplicate := seen[section.marker]; duplicate {
				return nil, fmt.Errorf("duplicate markdown unit marker %q", section.marker)
			}
			seen[section.marker] = struct{}{}
			section.key = section.marker
		}
	}
	return markdownEntries(targetSections), nil
}

func (f *MarkdownFormat) Serialize(entries map[string]string, original []byte) ([]byte, error) {
	sections, err := parseMarkdownSections(original)
	if err != nil {
		return nil, err
	}
	assignSourceKeys(sections)
	return serializeMarkdownSections(entries, sections, false)
}

// SerializeTarget follows source-section order and emits target-side markers.
func (f *MarkdownFormat) SerializeTarget(entries map[string]string, source, _ []byte) ([]byte, error) {
	sections, err := parseMarkdownSections(source)
	if err != nil {
		return nil, fmt.Errorf("parsing markdown source: %w", err)
	}
	assignSourceKeys(sections)
	return serializeMarkdownSections(entries, sections, true)
}

type markdownSection struct {
	key     string
	marker  string
	heading string
	value   string
}

func parseMarkdownSections(data []byte) ([]*markdownSection, error) {
	sections := []*markdownSection{{key: markdownPreambleKey}}
	current := sections[0]
	var value strings.Builder
	var pendingMarker string
	var fence byte
	var fenceLength int

	flush := func() {
		current.value = value.String()
		value.Reset()
	}
	for _, line := range splitLinesAfter(string(data)) {
		if match := markdownFencePattern.FindStringSubmatch(line); match != nil {
			marker := match[1]
			if fence == 0 {
				fence = marker[0]
				fenceLength = len(marker)
			} else if marker[0] == fence && len(marker) >= fenceLength {
				fence = 0
				fenceLength = 0
			}
			value.WriteString(line)
			continue
		}
		if fence == 0 {
			if match := markdownMarkerPattern.FindStringSubmatch(line); match != nil {
				if pendingMarker != "" {
					return nil, fmt.Errorf("markdown unit marker %q is not followed by a level-two heading", pendingMarker)
				}
				pendingMarker = match[1]
				continue
			}
			if match := markdownH2Pattern.FindStringSubmatch(line); match != nil {
				flush()
				current = &markdownSection{marker: pendingMarker, heading: strings.TrimSpace(match[1])}
				sections = append(sections, current)
				pendingMarker = ""
				value.WriteString(line)
				continue
			}
			if pendingMarker != "" {
				return nil, fmt.Errorf("markdown unit marker %q is not immediately followed by a level-two heading", pendingMarker)
			}
		}
		value.WriteString(line)
	}
	if pendingMarker != "" {
		return nil, fmt.Errorf("markdown unit marker %q has no section", pendingMarker)
	}
	flush()
	return sections, nil
}

func assignSourceKeys(sections []*markdownSection) {
	seen := make(map[string]int)
	sections[0].key = markdownPreambleKey
	for index, section := range sections[1:] {
		if section.marker != "" {
			section.key = section.marker
			continue
		}
		slug := markdownSlug(section.heading)
		if slug == "" {
			slug = fmt.Sprintf("section-%d", index+1)
		}
		seen[slug]++
		if seen[slug] > 1 {
			slug = fmt.Sprintf("%s-%d", slug, seen[slug])
		}
		section.key = "markdown:" + slug
	}
}

func markdownEntries(sections []*markdownSection) map[string]string {
	entries := make(map[string]string, len(sections))
	for _, section := range sections {
		entries[section.key] = section.value
	}
	return entries
}

func serializeMarkdownSections(entries map[string]string, order []*markdownSection, markers bool) ([]byte, error) {
	var output strings.Builder
	for index, section := range order {
		value, ok := entries[section.key]
		if !ok {
			return nil, fmt.Errorf("markdown entry %q is missing", section.key)
		}
		if index > 0 && markers {
			if output.Len() > 0 && !strings.HasSuffix(output.String(), "\n") {
				output.WriteByte('\n')
			}
			fmt.Fprintf(&output, "<!-- internationalizer:unit %s -->\n", section.key)
		}
		output.WriteString(value)
	}
	return []byte(output.String()), nil
}

func markdownSlug(heading string) string {
	var slug strings.Builder
	separator := false
	for _, r := range strings.ToLower(heading) {
		switch {
		case unicode.IsLetter(r) || unicode.IsNumber(r):
			if separator && slug.Len() > 0 {
				slug.WriteByte('-')
			}
			separator = false
			slug.WriteRune(r)
		default:
			separator = true
		}
	}
	return slug.String()
}

func splitLinesAfter(value string) []string {
	if value == "" {
		return nil
	}
	lines := strings.SplitAfter(value, "\n")
	if lines[len(lines)-1] == "" {
		return lines[:len(lines)-1]
	}
	return lines
}
