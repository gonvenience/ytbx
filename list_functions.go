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

// GetIdentifierFromNamedList returns the identifier key used in the provided
// list, or an empty string if there is none.
// The identifier key is either 'name', 'key', or 'id'.
func GetIdentifierFromNamedList(sequenceNode *yamlv3.Node) string {
	counters := map[string]int{}

	for _, mappingNode := range sequenceNode.Content {
		for i := 0; i < len(mappingNode.Content); i += 2 {
			k := mappingNode.Content[i]

			if _, ok := counters[k.Value]; !ok {
				counters[k.Value] = 0
			}

			counters[k.Value]++
		}
	}

	listLength := len(sequenceNode.Content)
	for _, identifier := range []string{"name", "key", "id"} {
		if count, ok := counters[identifier]; ok && count == listLength {
			return identifier
		}
	}

	return ""
}

func getEntryByIdentifierAndName(sequenceNode *yamlv3.Node, identifier string, name string) (*yamlv3.Node, error) {
	idx, err := getIndexByIdentifierAndName(sequenceNode, identifier, name)
	if err != nil {
		return nil, err
	}

	return sequenceNode.Content[idx], nil
}

func getIndexByIdentifierAndName(sequenceNode *yamlv3.Node, identifier string, name string) (int, error) {
	for idx, mappingNode := range sequenceNode.Content {
		for i := 0; i < len(mappingNode.Content); i += 2 {
			k, v := mappingNode.Content[i], mappingNode.Content[i+1]
			if k.Value == identifier && v.Value == name {
				return idx, nil
			}
		}
	}

	return -1,
		fmt.Errorf("there is no entry %s=%v in the list",
			identifier,
			name,
		)
}
