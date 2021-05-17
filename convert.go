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
	"strconv"
	"time"

	yamlv3 "gopkg.in/yaml.v3"
)

func asNode(obj interface{}) (*yamlv3.Node, error) {
	switch obj := obj.(type) {
	case string:
		return &yamlv3.Node{
			Kind:  yamlv3.ScalarNode,
			Tag:   "!!str",
			Value: obj,
		}, nil

	default:
		tmp, err := yamlv3.Marshal(obj)
		if err != nil {
			return nil, err
		}

		var node yamlv3.Node
		if err := yamlv3.Unmarshal(tmp, &node); err != nil {
			return nil, err
		}

		return &node, nil
	}
}

func asNodeP(obj interface{}) *yamlv3.Node {
	val, err := asNode(obj)
	if err != nil {
		panic(err)
	}

	return val
}

func asType(node *yamlv3.Node) (interface{}, error) {
	switch node.Kind {
	case yamlv3.DocumentNode:
		return asType(node.Content[0])

	case yamlv3.MappingNode:
		var tmp = map[interface{}]interface{}{}
		for i := 0; i < len(node.Content); i += 2 {
			key, err := asType(node.Content[i])
			if err != nil {
				return nil, err
			}

			val, err := asType(node.Content[i+1])
			if err != nil {
				return nil, err
			}

			tmp[key] = val
		}

		return tmp, nil

	case yamlv3.SequenceNode:
		var tmp = make([]interface{}, len(node.Content))
		for i := range node.Content {
			val, err := asType(node.Content[i])
			if err != nil {
				return nil, err
			}

			tmp[i] = val
		}

		return tmp, nil

	case yamlv3.ScalarNode:
		switch node.Tag {
		case "!!str":
			return node.Value, nil

		case "!!timestamp":
			return time.Parse(time.RFC3339, node.Value)

		case "!!int":
			return strconv.Atoi(node.Value)

		case "!!float":
			return strconv.ParseFloat(node.Value, 64)

		case "!!bool":
			return strconv.ParseBool(node.Value)

		case "!!null":
			return nil, nil

		default:
			return nil, fmt.Errorf("unknown YAML node tag %s", node.Tag)
		}
	}

	return nil, fmt.Errorf("failed to translate node (kind=%v, tag=%s) into specific type", node.Kind, node.Tag)
}
