package llm

import (
	"encoding/json"
	"fmt"
	"strings"
)

// ParseTranslationResponse extracts a map[string]string from LLM response text.
// It handles raw JSON, JSON wrapped in markdown code blocks, and nested JSON.
func ParseTranslationResponse(text string) (map[string]string, error) {
	text = strings.TrimSpace(text)

	// Try raw JSON first.
	if result, err := tryParseJSON(text); err == nil {
		return result, nil
	}

	// Try extracting from markdown code block.
	if idx := strings.Index(text, "```"); idx >= 0 {
		start := strings.Index(text[idx+3:], "\n")
		if start >= 0 {
			rest := text[idx+3+start+1:]
			if end := strings.Index(rest, "```"); end >= 0 {
				if result, err := tryParseJSON(strings.TrimSpace(rest[:end])); err == nil {
					return result, nil
				}
			}
		}
	}

	// Try finding the first { to last } span.
	first := strings.Index(text, "{")
	last := strings.LastIndex(text, "}")
	if first >= 0 && last > first {
		if result, err := tryParseJSON(text[first : last+1]); err == nil {
			return result, nil
		}
	}

	return nil, fmt.Errorf("could not parse translation response as JSON: %.200s", text)
}

func tryParseJSON(text string) (map[string]string, error) {
	// Decode through interface values so null and other non-string leaves cannot
	// silently become empty strings during unmarshalling.
	var nested map[string]interface{}
	if err := json.Unmarshal([]byte(text), &nested); err != nil {
		return nil, err
	}

	result := make(map[string]string)
	if err := flattenResponse("", nested, result); err != nil {
		return nil, err
	}
	return result, nil
}

func flattenResponse(prefix string, val map[string]interface{}, out map[string]string) error {
	for key, v := range val {
		p := key
		if prefix != "" {
			p = prefix + "." + key
		}
		switch child := v.(type) {
		case string:
			out[p] = child
		case map[string]interface{}:
			if err := flattenResponse(p, child, out); err != nil {
				return err
			}
		default:
			return fmt.Errorf("translation %q must be a string, got %T", p, v)
		}
	}
	return nil
}
