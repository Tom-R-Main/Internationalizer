package formats

// MarkdownFormat treats a whole markdown file as a single translation unit.
// The key is always "_content" and the value is the full document text.
type MarkdownFormat struct{}

func (f *MarkdownFormat) Name() string         { return "markdown" }
func (f *MarkdownFormat) Extensions() []string { return []string{".md", ".mdx"} }

func (f *MarkdownFormat) Parse(data []byte) (map[string]string, error) {
	return map[string]string{"_content": string(data)}, nil
}

func (f *MarkdownFormat) Serialize(entries map[string]string, _ []byte) ([]byte, error) {
	if content, ok := entries["_content"]; ok {
		return []byte(content), nil
	}
	return nil, nil
}
