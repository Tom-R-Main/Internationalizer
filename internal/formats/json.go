package formats

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"sort"
	"strings"
)

type JSONFormat struct{}

func (f *JSONFormat) Name() string         { return "json" }
func (f *JSONFormat) Extensions() []string { return []string{".json"} }

func (f *JSONFormat) Parse(data []byte) (map[string]string, error) {
	var raw interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("json parse: %w", err)
	}
	result := make(map[string]string)
	flatten("", raw, result)
	return result, nil
}

func flatten(prefix string, val interface{}, out map[string]string) {
	switch v := val.(type) {
	case map[string]interface{}:
		for key, child := range v {
			p := key
			if prefix != "" {
				p = prefix + "." + key
			}
			flatten(p, child, out)
		}
	case []interface{}:
		for i, child := range v {
			p := fmt.Sprintf("%s.%d", prefix, i)
			flatten(p, child, out)
		}
	default:
		out[prefix] = fmt.Sprintf("%v", v)
	}
}

func (f *JSONFormat) Serialize(entries map[string]string, original []byte) ([]byte, error) {
	// If we have original data, parse it to preserve key ordering.
	if len(original) > 0 {
		return serializePreservingOrder(entries, original)
	}
	// No original data — build a nested structure from scratch.
	return serializeFromScratch(entries)
}

// serializePreservingOrder walks the original JSON structure and replaces
// leaf values from the entries map, preserving key ordering.
func serializePreservingOrder(entries map[string]string, original []byte) ([]byte, error) {
	var raw interface{}
	dec := json.NewDecoder(bytes.NewReader(original))
	dec.UseNumber()
	if err := dec.Decode(&raw); err != nil {
		return nil, fmt.Errorf("json parse original: %w", err)
	}
	replaceLeaves("", raw, entries)
	for key, value := range entries {
		setPath(&raw, strings.Split(key, "."), value)
	}
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetIndent("", "  ")
	enc.SetEscapeHTML(false)
	if err := enc.Encode(raw); err != nil {
		return nil, err
	}
	// json.Encoder adds a trailing newline; trim then add exactly one.
	return bytes.TrimRight(buf.Bytes(), "\n"), nil
}

func replaceLeaves(prefix string, val interface{}, entries map[string]string) {
	switch v := val.(type) {
	case map[string]interface{}:
		for key, child := range v {
			p := key
			if prefix != "" {
				p = prefix + "." + key
			}
			switch child.(type) {
			case map[string]interface{}, []interface{}:
				replaceLeaves(p, child, entries)
			default:
				if replacement, ok := entries[p]; ok {
					v[key] = replacement
				}
			}
		}
	case []interface{}:
		for i, child := range v {
			p := fmt.Sprintf("%s.%d", prefix, i)
			switch child.(type) {
			case map[string]interface{}, []interface{}:
				replaceLeaves(p, child, entries)
			default:
				if replacement, ok := entries[p]; ok {
					v[i] = replacement
				}
			}
		}
	}
}

func setPath(target *interface{}, parts []string, value string) {
	if len(parts) == 0 {
		*target = value
		return
	}

	if idx, err := strconv.Atoi(parts[0]); err == nil {
		var arr []interface{}
		switch current := (*target).(type) {
		case nil:
			arr = make([]interface{}, idx+1)
		case []interface{}:
			arr = current
			if len(arr) <= idx {
				expanded := make([]interface{}, idx+1)
				copy(expanded, arr)
				arr = expanded
			}
		default:
			arr = make([]interface{}, idx+1)
		}
		child := arr[idx]
		setPath(&child, parts[1:], value)
		arr[idx] = child
		*target = arr
		return
	}

	var obj map[string]interface{}
	switch current := (*target).(type) {
	case nil:
		obj = make(map[string]interface{})
	case map[string]interface{}:
		obj = current
	default:
		obj = make(map[string]interface{})
	}
	child := obj[parts[0]]
	setPath(&child, parts[1:], value)
	obj[parts[0]] = child
	*target = obj
}

func serializeFromScratch(entries map[string]string) ([]byte, error) {
	root := make(map[string]interface{})
	// Sort keys for deterministic output.
	keys := make([]string, 0, len(entries))
	for k := range entries {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, key := range keys {
		parts := strings.Split(key, ".")
		current := root
		for i, part := range parts {
			if i == len(parts)-1 {
				current[part] = entries[key]
			} else {
				if _, ok := current[part]; !ok {
					current[part] = make(map[string]interface{})
				}
				current = current[part].(map[string]interface{})
			}
		}
	}

	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetIndent("", "  ")
	enc.SetEscapeHTML(false)
	if err := enc.Encode(root); err != nil {
		return nil, err
	}
	return bytes.TrimRight(buf.Bytes(), "\n"), nil
}
