// Package jsonintegrity decodes JSON without discarding duplicate members or
// ambiguous dotted identities used by localization catalogs.
package jsonintegrity

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// MaxDepth bounds recursive catalog processing, including downstream walkers.
const MaxDepth = 256

// Error describes an input integrity failure without including catalog values.
// Paths are RFC 6901 JSON pointers. Offsets distinguish duplicate members whose
// decoded pointers are necessarily identical.
type Error struct {
	Code        string `json:"code"`
	Key         string `json:"key"`
	Path        string `json:"path"`
	OtherPath   string `json:"other_path"`
	Offset      int64  `json:"offset"`
	OtherOffset int64  `json:"other_offset"`
}

func (e *Error) JSONCode() string { return e.Code }

func (e *Error) Error() string {
	switch e.Code {
	case "json_duplicate_member":
		return fmt.Sprintf("duplicate JSON member at %q (byte offsets %d and %d)", e.Path, e.OtherOffset, e.Offset)
	case "json_flattened_key_collision":
		return fmt.Sprintf("flattened JSON key %q collides between %q and %q", e.Key, e.OtherPath, e.Path)
	case "json_nesting_limit":
		return fmt.Sprintf("JSON nesting exceeds %d levels", MaxDepth)
	default:
		return "unexpected content after JSON document"
	}
}

type location struct {
	path   string
	offset int64
}

// Decode returns maps, slices, strings, bools, nil and json.Number values. It
// checks leaf identities, including metadata and empty containers, before a map
// insertion could erase information. The first error follows source order.
func Decode(data []byte) (any, error) {
	d := json.NewDecoder(bytes.NewReader(data))
	d.UseNumber()
	seen := make(map[string]location)
	value, err := decodeValue(d, "", "", 0, seen)
	if err != nil {
		return nil, err
	}
	if _, err := d.Token(); err != io.EOF {
		return nil, &Error{Code: "json_trailing_content", Offset: d.InputOffset()}
	}
	return value, nil
}

func decodeValue(d *json.Decoder, key, path string, depth int, seen map[string]location) (any, error) {
	if depth > MaxDepth {
		return nil, &Error{Code: "json_nesting_limit", Path: path}
	}
	offset := d.InputOffset()
	token, err := d.Token()
	if err != nil {
		return nil, err
	}
	switch token {
	case json.Delim('{'):
		object := make(map[string]any)
		members := make(map[string]int64)
		if !d.More() && depth > 0 {
			if err := registerLeaf(key, path, offset, seen); err != nil {
				return nil, err
			}
		}
		for d.More() {
			token, err := d.Token()
			if err != nil {
				return nil, err
			}
			member, ok := token.(string)
			if !ok {
				return nil, fmt.Errorf("expected JSON object member")
			}
			childPath := path + "/" + strings.ReplaceAll(strings.ReplaceAll(member, "~", "~0"), "/", "~1")
			childKey := member
			if key != "" {
				childKey = key + "." + member
			}
			if offset, exists := members[member]; exists {
				return nil, &Error{Code: "json_duplicate_member", Key: childKey, Path: childPath, OtherPath: childPath, Offset: d.InputOffset(), OtherOffset: offset}
			}
			members[member] = d.InputOffset()
			child, err := decodeValue(d, childKey, childPath, depth+1, seen)
			if err != nil {
				return nil, err
			}
			object[member] = child
		}
		if _, err := d.Token(); err != nil {
			return nil, err
		}
		return object, nil
	case json.Delim('['):
		array := make([]any, 0)
		if !d.More() && depth > 0 {
			if err := registerLeaf(key, path, offset, seen); err != nil {
				return nil, err
			}
		}
		for d.More() {
			index := strconv.Itoa(len(array))
			child, err := decodeValue(d, key+"."+index, path+"/"+index, depth+1, seen)
			if err != nil {
				return nil, err
			}
			array = append(array, child)
		}
		if _, err := d.Token(); err != nil {
			return nil, err
		}
		return array, nil
	default:
		if depth > 0 {
			if err := registerLeaf(key, path, offset, seen); err != nil {
				return nil, err
			}
		}
		return token, nil
	}
}

func registerLeaf(key, path string, offset int64, seen map[string]location) error {
	if previous, ok := seen[key]; ok {
		return &Error{Code: "json_flattened_key_collision", Key: key, Path: path, OtherPath: previous.path, Offset: offset, OtherOffset: previous.offset}
	}
	seen[key] = location{path, offset}
	return nil
}
