package formats

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/Tom-R-Main/Internationalizer/internal/jsonintegrity"
)

type JSONFormat struct{}

func (f *JSONFormat) Name() string         { return "json" }
func (f *JSONFormat) Extensions() []string { return []string{".json"} }

func (f *JSONFormat) Parse(data []byte) (map[string]string, error) {
	raw, err := jsonintegrity.Decode(data)
	if err != nil {
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
	raw, err := jsonintegrity.Decode(original)
	if err != nil {
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
	if _, err := jsonintegrity.Decode(buf.Bytes()); err != nil {
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
				if _, isString := child.(string); isString {
					delete(node, key)
					continue
				}
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
	raw, err := jsonintegrity.Decode(original)
	if err != nil {
		return nil, fmt.Errorf("json parse original: %w", err)
	}
	if raw == nil && len(entries) > 0 {
		return nil, fmt.Errorf("json insertion would discard root null metadata")
	}
	replaced := make(map[string]struct{}, len(entries))
	if _, rootString := raw.(string); rootString {
		if replacement, ok := entries[""]; ok {
			raw = replacement
			replaced[""] = struct{}{}
		}
	}
	replaceLeaves("", raw, entries, replaced)
	keys := make([]string, 0, len(entries))
	for key := range entries {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		if _, ok := replaced[key]; ok {
			continue
		}
		if strings.Count(key, ".") >= jsonintegrity.MaxDepth {
			return nil, &jsonintegrity.Error{Code: "json_nesting_limit"}
		}
		if err := setPath(&raw, strings.Split(key, "."), entries[key]); err != nil {
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
	if err := checkSerializedEntries(buf.Bytes(), entries); err != nil {
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
			case string:
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
			case string:
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
		if *target != nil {
			return fmt.Errorf("replacement would discard existing content")
		}
		*target = value
		return nil
	}

	if idx, err := strconv.Atoi(parts[0]); err == nil {
		if arr, ok := (*target).([]interface{}); ok {
			if idx < 0 || idx > len(arr) {
				return fmt.Errorf("array index %q is outside existing structure", parts[0])
			}
			if idx < len(arr) && (len(parts) == 1 || arr[idx] == nil) {
				return fmt.Errorf("replacement would discard existing array content")
			}
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
	child, exists := obj[parts[0]]
	if exists && (len(parts) == 1 || child == nil) {
		return fmt.Errorf("replacement would discard existing content")
	}
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
		if strings.Count(key, ".") >= jsonintegrity.MaxDepth {
			return nil, &jsonintegrity.Error{Code: "json_nesting_limit"}
		}
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
	if err := checkSerializedEntries(buf.Bytes(), entries); err != nil {
		return nil, err
	}
	return bytes.TrimRight(buf.Bytes(), "\n"), nil
}

func checkSerializedEntries(data []byte, entries map[string]string) error {
	parsed, err := (&JSONFormat{}).Parse(data)
	if err != nil {
		return err
	}
	keys := make([]string, 0, len(entries))
	for key := range entries {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		if value, exists := parsed[key]; !exists || value != entries[key] {
			return fmt.Errorf("json set path %q cannot preserve the requested identity", key)
		}
	}
	return nil
}
