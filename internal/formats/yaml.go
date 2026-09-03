package formats

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

type YAMLFormat struct{}

func (f *YAMLFormat) Name() string         { return "yaml" }
func (f *YAMLFormat) Extensions() []string { return []string{".yml", ".yaml"} }

func (f *YAMLFormat) Parse(data []byte) (map[string]string, error) {
	var doc yaml.Node
	if err := yaml.Unmarshal(data, &doc); err != nil {
		return nil, fmt.Errorf("yaml parse: %w", err)
	}
	result := make(map[string]string)
	if doc.Kind == yaml.DocumentNode && len(doc.Content) > 0 {
		flattenYAMLNode("", doc.Content[0], result)
	}
	return result, nil
}

func flattenYAMLNode(prefix string, node *yaml.Node, out map[string]string) {
	switch node.Kind {
	case yaml.MappingNode:
		for i := 0; i+1 < len(node.Content); i += 2 {
			key := node.Content[i].Value
			val := node.Content[i+1]
			p := key
			if prefix != "" {
				p = prefix + "." + key
			}
			flattenYAMLNode(p, val, out)
		}
	case yaml.SequenceNode:
		for i, child := range node.Content {
			p := fmt.Sprintf("%s.%d", prefix, i)
			flattenYAMLNode(p, child, out)
		}
	case yaml.ScalarNode:
		if node.Tag == "!!str" {
			out[prefix] = node.Value
		}
	}
}

func (f *YAMLFormat) Serialize(entries map[string]string, original []byte) ([]byte, error) {
	if len(original) > 0 {
		return serializeYAMLPreserving(entries, original)
	}
	return serializeYAMLFromScratch(entries)
}

func (f *YAMLFormat) RemoveEntries(original []byte, keys map[string]struct{}) ([]byte, error) {
	var doc yaml.Node
	if err := yaml.Unmarshal(original, &doc); err != nil {
		return nil, fmt.Errorf("yaml parse original: %w", err)
	}
	if doc.Kind == yaml.DocumentNode && len(doc.Content) > 0 {
		removeYAMLEntries("", doc.Content[0], keys)
	}
	return yaml.Marshal(&doc)
}

func removeYAMLEntries(prefix string, node *yaml.Node, keys map[string]struct{}) {
	switch node.Kind {
	case yaml.MappingNode:
		kept := node.Content[:0]
		for index := 0; index+1 < len(node.Content); index += 2 {
			keyNode := node.Content[index]
			valueNode := node.Content[index+1]
			path := keyNode.Value
			if prefix != "" {
				path = prefix + "." + keyNode.Value
			}
			if _, remove := keys[path]; remove {
				continue
			}
			removeYAMLEntries(path, valueNode, keys)
			kept = append(kept, keyNode, valueNode)
		}
		node.Content = kept
	case yaml.SequenceNode:
		for index, child := range node.Content {
			path := fmt.Sprintf("%s.%d", prefix, index)
			removeYAMLEntries(path, child, keys)
		}
	}
}

func serializeYAMLPreserving(entries map[string]string, original []byte) ([]byte, error) {
	var doc yaml.Node
	if err := yaml.Unmarshal(original, &doc); err != nil {
		return nil, fmt.Errorf("yaml parse original: %w", err)
	}
	replaced := make(map[string]struct{}, len(entries))
	if doc.Kind == yaml.DocumentNode && len(doc.Content) > 0 {
		root := doc.Content[0]
		replaceYAMLLeaves("", root, entries, replaced)
		keys := make([]string, 0, len(entries))
		for key := range entries {
			if _, ok := replaced[key]; !ok {
				keys = append(keys, key)
			}
		}
		sort.Strings(keys)
		for _, key := range keys {
			if err := setYAMLPath(root, strings.Split(key, "."), entries[key]); err != nil {
				return nil, fmt.Errorf("yaml set path %q: %w", key, err)
			}
		}
	}
	out, err := yaml.Marshal(&doc)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func replaceYAMLLeaves(prefix string, node *yaml.Node, entries map[string]string, replaced map[string]struct{}) {
	switch node.Kind {
	case yaml.MappingNode:
		for i := 0; i+1 < len(node.Content); i += 2 {
			key := node.Content[i].Value
			val := node.Content[i+1]
			p := key
			if prefix != "" {
				p = prefix + "." + key
			}
			if val.Kind == yaml.ScalarNode {
				if replacement, ok := entries[p]; ok {
					val.Value = replacement
					val.Tag = "!!str"
					replaced[p] = struct{}{}
				}
			} else {
				replaceYAMLLeaves(p, val, entries, replaced)
			}
		}
	case yaml.SequenceNode:
		for i, child := range node.Content {
			p := fmt.Sprintf("%s.%d", prefix, i)
			if child.Kind == yaml.ScalarNode {
				if replacement, ok := entries[p]; ok {
					child.Value = replacement
					child.Tag = "!!str"
					replaced[p] = struct{}{}
				}
			} else {
				replaceYAMLLeaves(p, child, entries, replaced)
			}
		}
	}
}

func serializeYAMLFromScratch(entries map[string]string) ([]byte, error) {
	root := &yaml.Node{Kind: yaml.MappingNode}
	for key, value := range entries {
		parts := strings.Split(key, ".")
		if err := setYAMLPath(root, parts, value); err != nil {
			return nil, fmt.Errorf("yaml set path %q: %w", key, err)
		}
	}
	doc := &yaml.Node{Kind: yaml.DocumentNode, Content: []*yaml.Node{root}}
	return yaml.Marshal(doc)
}

func setYAMLPath(node *yaml.Node, parts []string, value string) error {
	if len(parts) == 0 {
		return nil
	}
	if node.Kind == yaml.SequenceNode {
		index, err := strconv.Atoi(parts[0])
		if err != nil || index < 0 {
			return fmt.Errorf("sequence requires a non-negative numeric segment, got %q", parts[0])
		}
		for len(node.Content) <= index {
			node.Content = append(node.Content, &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!null", Value: "null"})
		}
		child := node.Content[index]
		if len(parts) == 1 {
			child.Kind = yaml.ScalarNode
			child.Tag = "!!str"
			child.Value = value
			child.Content = nil
			return nil
		}
		ensureYAMLContainer(child, parts[1])
		return setYAMLPath(child, parts[1:], value)
	}
	if node.Kind != yaml.MappingNode {
		node.Kind = yaml.MappingNode
		node.Tag = "!!map"
		node.Value = ""
		node.Content = nil
	}
	// Find or create the key node.
	var valNode *yaml.Node
	for i := 0; i+1 < len(node.Content); i += 2 {
		if node.Content[i].Value == parts[0] {
			valNode = node.Content[i+1]
			break
		}
	}
	if valNode == nil {
		keyNode := &yaml.Node{Kind: yaml.ScalarNode, Value: parts[0]}
		if len(parts) == 1 {
			valNode = &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!str", Value: value}
		} else {
			valNode = &yaml.Node{}
			ensureYAMLContainer(valNode, parts[1])
		}
		node.Content = append(node.Content, keyNode, valNode)
	}
	if len(parts) == 1 {
		valNode.Value = value
		valNode.Kind = yaml.ScalarNode
		valNode.Tag = "!!str"
		valNode.Content = nil
		return nil
	}
	ensureYAMLContainer(valNode, parts[1])
	return setYAMLPath(valNode, parts[1:], value)
}

func ensureYAMLContainer(node *yaml.Node, nextPart string) {
	if node.Kind == yaml.MappingNode || node.Kind == yaml.SequenceNode {
		return
	}
	if _, err := strconv.Atoi(nextPart); err == nil {
		node.Kind = yaml.SequenceNode
		node.Tag = "!!seq"
		node.Value = ""
		node.Content = nil
		return
	}
	node.Kind = yaml.MappingNode
	node.Tag = "!!map"
	node.Value = ""
	node.Content = nil
}
