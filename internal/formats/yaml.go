package formats

import (
	"fmt"
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
		out[prefix] = node.Value
	}
}

func (f *YAMLFormat) Serialize(entries map[string]string, original []byte) ([]byte, error) {
	if len(original) > 0 {
		return serializeYAMLPreserving(entries, original)
	}
	return serializeYAMLFromScratch(entries)
}

func serializeYAMLPreserving(entries map[string]string, original []byte) ([]byte, error) {
	var doc yaml.Node
	if err := yaml.Unmarshal(original, &doc); err != nil {
		return nil, fmt.Errorf("yaml parse original: %w", err)
	}
	if doc.Kind == yaml.DocumentNode && len(doc.Content) > 0 {
		replaceYAMLLeaves("", doc.Content[0], entries)
	}
	out, err := yaml.Marshal(&doc)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func replaceYAMLLeaves(prefix string, node *yaml.Node, entries map[string]string) {
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
				}
			} else {
				replaceYAMLLeaves(p, val, entries)
			}
		}
	case yaml.SequenceNode:
		for i, child := range node.Content {
			p := fmt.Sprintf("%s.%d", prefix, i)
			if child.Kind == yaml.ScalarNode {
				if replacement, ok := entries[p]; ok {
					child.Value = replacement
				}
			} else {
				replaceYAMLLeaves(p, child, entries)
			}
		}
	}
}

func serializeYAMLFromScratch(entries map[string]string) ([]byte, error) {
	root := &yaml.Node{Kind: yaml.MappingNode}
	for key, value := range entries {
		parts := strings.Split(key, ".")
		setYAMLPath(root, parts, value)
	}
	doc := &yaml.Node{Kind: yaml.DocumentNode, Content: []*yaml.Node{root}}
	return yaml.Marshal(doc)
}

func setYAMLPath(node *yaml.Node, parts []string, value string) {
	if len(parts) == 0 {
		return
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
			valNode = &yaml.Node{Kind: yaml.ScalarNode, Value: value}
		} else {
			valNode = &yaml.Node{Kind: yaml.MappingNode}
		}
		node.Content = append(node.Content, keyNode, valNode)
	}
	if len(parts) == 1 {
		valNode.Value = value
		valNode.Kind = yaml.ScalarNode
		return
	}
	if valNode.Kind != yaml.MappingNode {
		valNode.Kind = yaml.MappingNode
		valNode.Content = nil
	}
	setYAMLPath(valNode, parts[1:], value)
}
