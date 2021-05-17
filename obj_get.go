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

	yamlv3 "gopkg.in/yaml.v3"
)

// GetPath is a convenience function for Get, which parses the given path
// and then delegates to the Get function
func GetPath(node *yamlv3.Node, pathString string) (*yamlv3.Node, error) {
	path, err := ParsePathString(pathString)
	if err != nil {
		return nil, err
	}

	switch node.Kind {
	case yamlv3.DocumentNode:
		return Get(node.Content[0], *path)

	default:
		return Get(node, *path)
	}
}

// Get retrieves the value at the provided Path in the given node
func Get(node *yamlv3.Node, path Path) (*yamlv3.Node, error) {
	if node.Kind == yamlv3.DocumentNode {
		node = node.Content[0]
	}

	pointer := node
	pointerPath := Path{docIdx: path.docIdx}

	for _, element := range path.sections {
		switch element := element.(type) {
		// Key/Value Map, where the element name is the key for the map
		case mappingNameSection:
			if pointer.Kind != yamlv3.MappingNode {
				return nil,
					fmt.Errorf("failed to traverse tree, expected %s but found type %s at %s",
						typeMap,
						GetType(pointer),
						pointerPath.GoPatchStyle(),
					)
			}

			entry, err := getValueByKey(pointer, element.name)
			if err != nil {
				return nil, err
			}

			pointer = entry

		// Complex List, where each list entry is a Key/Value map and the entry is
		// identified by name using an identifier (e.g. name, key, or id)
		case listNamedSection:
			if pointer.Kind != yamlv3.SequenceNode {
				return nil,
					fmt.Errorf("failed to traverse tree, expected %s but found type %s at %s",
						typeComplexList,
						GetType(pointer),
						pointerPath.GoPatchStyle(),
					)
			}

			entry, err := getEntryByIdentifierAndName(pointer, element.id, element.name)
			if err != nil {
				return nil, err
			}

			pointer = entry

		// Simple List (identified by index)
		case listIdxSection:
			if pointer.Kind != yamlv3.SequenceNode {
				return nil,
					fmt.Errorf("failed to traverse tree, expected %s but found type %s at %s",
						typeSimpleList,
						GetType(pointer),
						pointerPath.GoPatchStyle(),
					)
			}

			if element.idx == -1 {
				element.idx = len(pointer.Content) - 1
			}

			if element.idx < 0 || element.idx >= len(pointer.Content) {
				return nil,
					fmt.Errorf(
						"failed to traverse tree, provided %s index %d is not in range: 0..%d",
						typeSimpleList,
						element.idx,
						len(pointer.Content)-1,
					)
			}

			pointer = pointer.Content[element.idx]

		case undefSection:
			panic("not implemented yet")

		default:
			return nil, fmt.Errorf(
				"failed to traverse tree, the provided path %s seems to be invalid",
				path,
			)
		}

		// Update the path that the current pointer to keep track of the traversing
		pointerPath.sections = append(pointerPath.sections, element)
	}

	return pointer, nil
}
