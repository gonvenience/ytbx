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

package ytbx_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	. "github.com/gonvenience/bunt"
	. "github.com/gonvenience/ytbx"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/types"

	yamlv3 "gopkg.in/yaml.v3"
)

var exampleTOML = `
required = ["gopkg.in/fsnotify.v1"]

[prune]
  go-tests = true
  unused-packages = true
  non-go = true

[[constraint]]
  name = "gopkg.in/fsnotify.v1"
  source = "https://github.com/fsnotify/fsnotify.git"

[[constraint]]
  name = "k8s.io/helm"
  branch = "release-2.10"

[[override]]
  name = "gopkg.in/yaml.v2"
  revision = "670d4cfef0544295bc27a114dbac37980d83185a"

[[override]]
  branch = "release-1.10"
  name = "k8s.io/api"

[[override]]
  branch = "release-1.10"
  name = "k8s.io/apimachinery"


[[override]]
  branch = "release-7.0"
  name = "k8s.io/client-go"
`

func TestYtbx(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "ytbx suite")
}

var _ = BeforeSuite(func() {
	SetColorSettings(OFF, OFF)
})

var _ = AfterSuite(func() {
	SetColorSettings(AUTO, AUTO)
})

func assets(pathElement ...string) string {
	targetPath := filepath.Join(append(
		[]string{"assets"},
		pathElement...,
	)...)

	abs, err := filepath.Abs(targetPath)
	Expect(err).ToNot(HaveOccurred())

	return abs
}

func parseGoPatch(path string) Path {
	result, err := ParseGoPatchStylePathString(path)
	Expect(err).ToNot(HaveOccurred())

	return *result
}

func loadFile(filename string) *File {
	file, err := LoadFile(filename)
	Expect(err).ToNot(HaveOccurred())
	return &file
}

func yml(input string) *yamlv3.Node {
	// If input is a file location, load this as YAML
	if _, err := os.Open(input); err == nil {
		var content File
		var err error
		if content, err = LoadFile(input); err != nil {
			Fail(fmt.Sprintf("Failed to load YAML MapSlice from '%s': %s", input, err.Error()))
		}

		if len(content.Documents) > 1 {
			Fail(fmt.Sprintf("Failed to load YAML MapSlice from '%s': Provided file contains more than one document", input))
		}

		return content.Documents[0]
	}

	// Load YAML by parsing the actual string as YAML if it was not a file location
	document := singleDoc(input)
	return document.Content[0]
}

func list(input string) *yamlv3.Node {
	document := singleDoc(input)
	return document.Content[0]
}

func singleDoc(input string) *yamlv3.Node {
	docs, err := LoadYAMLDocuments([]byte(input))
	if err != nil {
		Fail(fmt.Sprintf("Failed to parse as YAML:\n%s\n\n%v", input, err))
	}

	if len(docs) > 1 {
		Fail(fmt.Sprintf("Failed to use YAML, because it contains multiple documents:\n%s\n", input))
	}

	return docs[0]
}

func BeAsNode(expected *yamlv3.Node) GomegaMatcher {
	return &nodeMatcher{
		expected: expected,
	}
}

type nodeMatcher struct {
	expected *yamlv3.Node
}

func (matcher *nodeMatcher) Match(actual interface{}) (success bool, err error) {
	actualNodePtr, ok := actual.(*yamlv3.Node)
	if !ok {
		return false, fmt.Errorf("BeAsNode matcher expected a Go YAML v3 Node, not %T", actual)
	}

	return isSameNode(actualNodePtr, matcher.expected)
}

func (matcher *nodeMatcher) FailureMessage(actual interface{}) string {
	return fmt.Sprintf("Expected\n\t%#v\n"+"to be same as\n\t%#v",
		actual,
		matcher.expected)
}

func (matcher *nodeMatcher) NegatedFailureMessage(actual interface{}) string {
	return fmt.Sprintf("Expected\n\t%#v\nnot to be same as\n\t%#v",
		actual,
		matcher.expected,
	)
}

func isSameNode(a, b *yamlv3.Node) (bool, error) {
	if a == nil && b == nil {
		return true, nil
	}

	if (a == nil && b != nil) || (a != nil && b == nil) {
		return false, nil
	}

	if a.Kind != b.Kind {
		return false, nil
	}

	if a.Tag != b.Tag {
		return false, nil
	}

	if a.Value != b.Value {
		return false, nil
	}

	if len(a.Content) != len(b.Content) {
		return false, nil
	}

	for i := range a.Content {
		if same, err := isSameNode(a.Content[i], b.Content[i]); !same {
			return same, err
		}
	}

	return true, nil
}
