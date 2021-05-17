// Copyright © 2021 The Homeport Team
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

// SetPath is a convenience function for Set, which parses the given path
// and then delegates to the Set function
func SetPath(node *yamlv3.Node, pathString string, value interface{}) error {
	path, err := ParsePathString(pathString)
	if err != nil {
		return err
	}

	return Set(node, *path, value)
}

// Set creates or updates the value at the provided Path in the given Node
func Set(node *yamlv3.Node, path Path, value interface{}) error {
	var pointer = node
	var pathPtr = Path{root: path.root, docIdx: path.docIdx}

	var create = func(idx int) (*yamlv3.Node, error) {
		switch {
		case idx == len(path.sections):
			return asNode(value)

		case idx > -1 && idx < len(path.sections):
			switch path.sections[idx].sectionType() {
			case mappingEntry:
				return &yamlv3.Node{
					Kind: yamlv3.MappingNode,
				}, nil

			case namedListEntry:
				return &yamlv3.Node{
					Kind: yamlv3.SequenceNode,
					Content: []*yamlv3.Node{{
						Kind: yamlv3.MappingNode,
						Content: []*yamlv3.Node{
							asNodeP(path.sections[idx].(listNamedSection).id),
							asNodeP(path.sections[idx].(listNamedSection).name),
						},
					}},
				}, nil

			case indexedListEntry:
				return &yamlv3.Node{
					Kind: yamlv3.SequenceNode,
				}, nil

			default:
				return nil, fmt.Errorf("failed to create new entry of unknown type %v", path.sections[idx].sectionType())
			}

		default:
			return nil, fmt.Errorf("failed to create new entry, section index %d is out of bounds [0..%d]", idx, len(path.sections)-1)
		}
	}

	for i, section := range path.sections {
		switch section := section.(type) {
		case mappingNameSection:
			if err := typeCheck(pathPtr, pointer, yamlv3.MappingNode); err != nil {
				return err
			}

			entry, err := getValueByKey(pointer, section.name)
			if err != nil {
				if entry, err = create(i + 1); err != nil {
					return err
				}

				pointer.Content = append(pointer.Content, asNodeP(section.name), entry)
			}

			pointer = entry

		case listNamedSection:
			if err := typeCheck(pathPtr, pointer, yamlv3.SequenceNode); err != nil {
				return err
			}

			entry, err := getEntryByIdentifierAndName(pointer, section.id, section.name)
			if err != nil {
				if entry, err = create(i + 1); err != nil {
					return err
				}

				pointer.Content = append(pointer.Content, entry)
			}

			pointer = entry

		case listIdxSection:
			if err := typeCheck(pathPtr, pointer, yamlv3.SequenceNode); err != nil {
				return err
			}

			if section.idx == -1 {
				entry, err := create(i + 1)
				if err != nil {
					return err
				}

				pointer.Content = append(pointer.Content, entry)
				section.idx = len(pointer.Content) - 1
			}

			if section.idx < 0 || section.idx >= len(pointer.Content) {
				panic("not implemented yet")
			}

			pointer = pointer.Content[section.idx]

		case undefSection:
			panic("not implemented yet")

		default:
			panic("not implemented yet")
		}
	}

	return nil
}
