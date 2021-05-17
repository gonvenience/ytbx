// Copyright © 2018 The Homeport Team
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package ytbx

import (
	"fmt"
	"reflect"
	"strings"

	yamlv3 "gopkg.in/yaml.v3"
)

// Internal string constants for type names and type decisions
const (
	typeMap         = "map"
	typeSimpleList  = "list"
	typeComplexList = "complex-list"
)

// IsStdin checks whether the provided input location refers to the dash
// character which usually serves as the replacement to point to STDIN rather
// than a file.
func IsStdin(location string) bool {
	return strings.TrimSpace(location) == "-"
}

// GetType returns the type of the input value with a YAML specific view
func GetType(value interface{}) string {
	switch tobj := value.(type) {
	case *yamlv3.Node:
		switch tobj.Kind {
		case yamlv3.MappingNode:
			return typeMap

		case yamlv3.SequenceNode:
			if hasMappingNodes(tobj) {
				return typeComplexList
			}

			return typeSimpleList

		default:
			return reflect.TypeOf(tobj.Value).Kind().String()
		}

	default:
		return reflect.TypeOf(value).Kind().String()
	}
}

func hasMappingNodes(sequenceNode *yamlv3.Node) bool {
	counter := 0

	for _, entry := range sequenceNode.Content {
		if entry.Kind == yamlv3.MappingNode {
			counter++
		}
	}

	return counter == len(sequenceNode.Content)
}

func typeCheck(path Path, node *yamlv3.Node, kind yamlv3.Kind) error {
	if node.Kind != kind {
		return fmt.Errorf("unexpected element in tree, expected %v but found type %v at %v",
			kind,
			node.Kind,
			path,
		)
	}

	return nil
}

func traverseTree(path Path, parent *yamlv3.Node, node *yamlv3.Node, leafFunc func(path Path, parent *yamlv3.Node, leaf *yamlv3.Node)) {
	switch node.Kind {
	case yamlv3.DocumentNode:
		traverseTree(
			path,
			node,
			node.Content[0],
			leafFunc,
		)

	case yamlv3.SequenceNode:
		if identifier := GetIdentifierFromNamedList(node); identifier != "" {
			for _, mappingNode := range node.Content {
				name, _ := getValueByKey(mappingNode, identifier)
				tmpPath := NewPathWithNamedEntryListSection(path, identifier, name.Value)
				for i := 0; i < len(mappingNode.Content); i += 2 {
					k, v := mappingNode.Content[i], mappingNode.Content[i+1]
					if k.Value == identifier { // skip the identifier mapping entry
						continue
					}

					traverseTree(
						NewPathWithNamedEntrySection(tmpPath, k.Value),
						node,
						v,
						leafFunc,
					)
				}
			}

		} else {
			for idx, entry := range node.Content {
				traverseTree(
					NewPathWithIndexedEntrySection(path, idx),
					node,
					entry,
					leafFunc,
				)
			}
		}

	case yamlv3.MappingNode:
		for i := 0; i < len(node.Content); i += 2 {
			k, v := node.Content[i], node.Content[i+1]
			traverseTree(
				NewPathWithNamedEntrySection(path, k.Value),
				node,
				v,
				leafFunc,
			)
		}

	default:
		leafFunc(path, parent, node)
	}
}
