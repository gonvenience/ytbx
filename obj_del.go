// Copyright © 2020 The Homeport Team
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

// DeletePath is a convenience function for Delete, which parses the given path
// and then delegates to the Delete function
func DeletePath(node *yamlv3.Node, pathString string) (*yamlv3.Node, error) {
	path, err := ParsePathString(pathString)
	if err != nil {
		return nil, err
	}

	switch node.Kind {
	case yamlv3.DocumentNode:
		return Delete(node.Content[0], *path)

	default:
		return Delete(node, *path)
	}
}

// Delete removes the provided Path from the given node
func Delete(node *yamlv3.Node, path Path) (*yamlv3.Node, error) {
	parentPath, err := path.Parent()
	if err != nil {
		return nil, err
	}

	parent, err := Get(node, parentPath)
	if err != nil {
		return nil, err
	}

	var (
		lastPathElement              = path.sections[len(path.sections)-1]
		deletedNode     *yamlv3.Node = nil
	)

	switch parent.Kind {
	case yamlv3.MappingNode:
		var deleteIdx int
		for i := 0; i < len(parent.Content); i += 2 {
			k, v := parent.Content[i], parent.Content[i+1]

			if k.Value == lastPathElement.(mappingNameSection).name {
				deleteIdx = i
				deletedNode = v
				break
			}
		}

		// delete the entry at delete index and the one after that as these two are
		// the key (first entry) and the value (second entry)
		parent.Content = append(
			parent.Content[:deleteIdx],
			parent.Content[deleteIdx+2:]...,
		)

		return deletedNode, nil

	case yamlv3.SequenceNode:
		var deleteIdx int
		switch obj := lastPathElement.(type) {
		case listIdxSection:
			deleteIdx = obj.idx

		case listNamedSection:
			deleteIdx, err = getIndexByIdentifierAndName(parent, obj.id, obj.name)
			if err != nil {
				return nil, err
			}

		default:
			return nil, fmt.Errorf("illegal type %T in sequence node delete", lastPathElement)
		}

		deletedNode = parent.Content[deleteIdx]

		// delete the entry that was identified by the deletion index, since it is a
		// sequence (list), only one entry needs to be deleted
		parent.Content = append(
			parent.Content[:deleteIdx],
			parent.Content[deleteIdx+1:]...,
		)

		return deletedNode, nil
	}

	return nil, fmt.Errorf("failed to delete path %s, because it could not be found", path)
}
