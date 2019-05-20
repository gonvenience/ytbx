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
	"reflect"
	"runtime"
	"testing"

	"github.com/homeport/gonvenience/pkg/v1/bunt"
	"github.com/homeport/gonvenience/pkg/v1/neat"
	"github.com/homeport/ytbx/pkg/v1/ytbx"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	yaml "gopkg.in/yaml.v2"
)

var assetsDirectory string

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
	bunt.ColorSetting = bunt.OFF

	_, file, _, ok := runtime.Caller(0)
	Expect(ok).To(BeTrue())

	dir, err := filepath.Abs(filepath.Dir(file) + "/../../../assets")
	Expect(err).To(BeNil())

	assetsDirectory = dir
})

func yml(input string) yaml.MapSlice {
	// If input is a file loacation, load this as YAML
	if _, err := os.Open(input); err == nil {
		var content ytbx.InputFile
		var err error
		if content, err = ytbx.LoadFile(input); err != nil {
			Fail(fmt.Sprintf("Failed to load YAML MapSlice from '%s': %s", input, err.Error()))
		}

		if len(content.Documents) > 1 {
			Fail(fmt.Sprintf("Failed to load YAML MapSlice from '%s': Provided file contains more than one document", input))
		}

		switch content.Documents[0].(type) {
		case yaml.MapSlice:
			return content.Documents[0].(yaml.MapSlice)
		}

		Fail(fmt.Sprintf("Failed to load YAML MapSlice from '%s': Document #0 in YAML is not of type MapSlice, but is %s", input, reflect.TypeOf(content.Documents[0])))
	}

	// Load YAML by parsing the actual string as YAML if it was not a file location
	doc := singleDoc(input)
	switch mapslice := doc.(type) {
	case yaml.MapSlice:
		return mapslice
	}

	Fail(fmt.Sprintf("Failed to use YAML, parsed data is not a YAML MapSlice:\n%s\n", input))
	return nil
}

func list(input string) []interface{} {
	doc := singleDoc(input)

	switch tobj := doc.(type) {
	case []interface{}:
		return tobj

	case []yaml.MapSlice:
		return ytbx.SimplifyList(tobj)
	}

	Fail(fmt.Sprintf("Failed to use YAML, parsed data is not a slice of any kind:\n%s\nIt was parsed as: %#v", input, doc))
	return nil
}

func singleDoc(input string) interface{} {
	docs, err := ytbx.LoadYAMLDocuments([]byte(input))
	if err != nil {
		Fail(fmt.Sprintf("Failed to parse as YAML:\n%s\n\n%v", input, err))
	}

	if len(docs) > 1 {
		Fail(fmt.Sprintf("Failed to use YAML, because it contains multiple documents:\n%s\n", input))
	}

	return docs[0]
}

func grab(obj interface{}, path string) interface{} {
	value, err := ytbx.Grab(obj, path)
	if err != nil {
		out, _ := neat.ToYAMLString(obj)
		Fail(fmt.Sprintf("Failed to grab by path %s from %s", path, out))
	}

	return value
}

func grabError(obj interface{}, path string) string {
	value, err := ytbx.Grab(obj, path)
	Expect(value).To(BeNil())
	return err.Error()
}
