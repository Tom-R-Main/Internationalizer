package formats

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
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
	case string:
		out[prefix] = v
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

func (f *JSONFormat) RemoveEntries(original []byte, keys map[string]struct{}) ([]byte, error) {
	var raw interface{}
	dec := json.NewDecoder(bytes.NewReader(original))
	dec.UseNumber()
	if err := dec.Decode(&raw); err != nil {
		return nil, fmt.Errorf("json parse original: %w", err)
	}
	removeJSONEntries("", raw, keys)
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetIndent("", "  ")
	enc.SetEscapeHTML(false)
	if err := enc.Encode(raw); err != nil {
		return nil, err
	}
	return bytes.TrimRight(buf.Bytes(), "\n"), nil
}

func removeJSONEntries(prefix string, value interface{}, keys map[string]struct{}) {
	switch node := value.(type) {
	case map[string]interface{}:
		for key, child := range node {
			path := key
			if prefix != "" {
				path = prefix + "." + key
			}
			if _, remove := keys[path]; remove {
				delete(node, key)
				continue
			}
			removeJSONEntries(path, child, keys)
		}
	case []interface{}:
		for index, child := range node {
			path := fmt.Sprintf("%s.%d", prefix, index)
			removeJSONEntries(path, child, keys)
		}
	}
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
	replaced := make(map[string]struct{}, len(entries))
	replaceLeaves("", raw, entries, replaced)
	for key, value := range entries {
		if _, ok := replaced[key]; ok {
			continue
		}
		if err := setPath(&raw, strings.Split(key, "."), value); err != nil {
			return nil, fmt.Errorf("json set path %q: %w", key, err)
		}
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

func replaceLeaves(prefix string, val interface{}, entries map[string]string, replaced map[string]struct{}) {
	switch v := val.(type) {
	case map[string]interface{}:
		for key, child := range v {
			p := key
			if prefix != "" {
				p = prefix + "." + key
			}
			switch child.(type) {
			case map[string]interface{}, []interface{}:
				replaceLeaves(p, child, entries, replaced)
			default:
				if replacement, ok := entries[p]; ok {
					v[key] = replacement
					replaced[p] = struct{}{}
				}
			}
		}
	case []interface{}:
		for i, child := range v {
			p := fmt.Sprintf("%s.%d", prefix, i)
			switch child.(type) {
			case map[string]interface{}, []interface{}:
				replaceLeaves(p, child, entries, replaced)
			default:
				if replacement, ok := entries[p]; ok {
					v[i] = replacement
					replaced[p] = struct{}{}
				}
			}
		}
	}
}

func setPath(target *interface{}, parts []string, value string) error {
	if len(parts) == 0 {
		*target = value
		return nil
	}

	if idx, err := strconv.Atoi(parts[0]); err == nil {
		if arr, ok := (*target).([]interface{}); ok {
			if len(arr) <= idx {
				expanded := make([]interface{}, idx+1)
				copy(expanded, arr)
				arr = expanded
			}
			child := arr[idx]
			if err := setPath(&child, parts[1:], value); err != nil {
				return err
			}
			arr[idx] = child
			*target = arr
			return nil
		}
		if *target == nil {
			return fmt.Errorf("numeric segment %q has no source structure", parts[0])
		}
	}

	var obj map[string]interface{}
	switch current := (*target).(type) {
	case nil:
		obj = make(map[string]interface{})
	case map[string]interface{}:
		obj = current
	default:
		return fmt.Errorf("object segment %q conflicts with existing array or scalar", parts[0])
	}
	child := obj[parts[0]]
	if err := setPath(&child, parts[1:], value); err != nil {
		return err
	}
	obj[parts[0]] = child
	*target = obj
	return nil
}

func serializeFromScratch(entries map[string]string) ([]byte, error) {
	var root interface{} = make(map[string]interface{})
	// Sort keys for deterministic output.
	keys := make([]string, 0, len(entries))
	for k := range entries {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, key := range keys {
		if err := setPath(&root, strings.Split(key, "."), entries[key]); err != nil {
			return nil, fmt.Errorf("json set path %q: %w", key, err)
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
