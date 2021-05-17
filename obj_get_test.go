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
	"strconv"

	. "github.com/gonvenience/ytbx"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/gonvenience/neat"
	yamlv3 "gopkg.in/yaml.v3"
)

func get(node *yamlv3.Node, path string) interface{} {
	v, err := GetPath(node, path)
	if err != nil {
		out, _ := neat.ToYAMLString(node)
		Fail(fmt.Sprintf("Failed to grab by path %s from %s", path, out))
	}

	switch v.Tag {
	case "!!str":
		return v.Value

	case "!!int":
		i, _ := strconv.Atoi(v.Value)
		return i
	}

	return v
}

func getToError(node *yamlv3.Node, path string) string {
	value, err := GetPath(node, path)
	Expect(value).To(BeNil())
	Expect(err).ToNot(BeNil())
	return err.Error()
}

var _ = Describe("reading entries", func() {
	Context("reading using File struct", func() {
		var file *File

		BeforeEach(func() {
			file = loadFile(assets("examples", "types.yml"))
		})

		AfterEach(func() {
			file = nil
		})

		It("should return the value referenced by the path", func() {
			Expect(file.Get(parseGoPatch("/yaml/map/before"))).To(BeEquivalentTo("after"))
			Expect(file.Get(parseGoPatch("/yaml/map/intA"))).To(BeEquivalentTo(42))
			Expect(file.Get(parseGoPatch("/yaml/simple-list/1"))).To(BeEquivalentTo("B"))
		})
	})

	Context("reading values by path", func() {
		It("should return the value referenced by the path", func() {
			example := yml(assets("examples", "types.yml"))
			Expect(get(example, "/yaml/map/before")).To(BeEquivalentTo("after"))
			Expect(get(example, "/yaml/map/intA")).To(BeEquivalentTo(42))
			Expect(get(example, "/yaml/map/mapA")).To(BeAsNode(yml(`{ key0: A, key1: A }`)))
			Expect(get(example, "/yaml/map/listA")).To(BeAsNode(list(`[ A, A, A ]`)))
			Expect(get(example, "/yaml/named-entry-list-using-name/name=B")).To(BeAsNode(yml(`{ name: B, foo: bar }`)))
			Expect(get(example, "/yaml/named-entry-list-using-key/key=B")).To(BeAsNode(yml(`{ key: B, foo: bar }`)))
			Expect(get(example, "/yaml/named-entry-list-using-id/id=B")).To(BeAsNode(yml(`{ id: B, foo: bar }`)))
			Expect(get(example, "/yaml/simple-list/1")).To(BeEquivalentTo("B"))
			Expect(get(example, "/yaml/named-entry-list-using-key/3")).To(BeAsNode(yml(`{ key: X, foo: bar }`)))

			example = yml(assets("bosh-yaml", "manifest.yml"))
			Expect(get(example, "/instance_groups/name=web/networks/name=concourse/static_ips/0")).To(BeEquivalentTo("XX.XX.XX.XX"))
			Expect(get(example, "/instance_groups/name=worker/jobs/name=baggageclaim/properties")).To(BeAsNode(yml(`{}`)))
		})

		It("should return the whole tree if root is referenced", func() {
			example := yml(assets("examples", "types.yml"))
			Expect(get(example, "/")).To(BeAsNode(example.Content[0]))
		})

		It("should return useful error messages", func() {
			example := yml(assets("examples", "types.yml"))
			Expect(getToError(example, "/yaml/simple-list/5")).To(BeEquivalentTo("failed to traverse tree, provided list index 5 is not in range: 0..4"))
			Expect(getToError(example, "/yaml/does-not-exist")).To(BeEquivalentTo("no key 'does-not-exist' found in map, available keys: map, simple-list, named-entry-list-using-name, named-entry-list-using-key, named-entry-list-using-id"))
			Expect(getToError(example, "/yaml/0")).To(BeEquivalentTo("failed to traverse tree, expected list but found type map at /yaml"))
			Expect(getToError(example, "/yaml/simple-list/foobar")).To(BeEquivalentTo("failed to traverse tree, expected map but found type list at /yaml/simple-list"))
			Expect(getToError(example, "/yaml/map/foobar=0")).To(BeEquivalentTo("failed to traverse tree, expected complex-list but found type map at /yaml/map"))
			Expect(getToError(example, "/yaml/named-entry-list-using-id/id=0")).To(BeEquivalentTo("there is no entry id=0 in the list"))
		})
	})

	Context("Trying to get values by path in an empty file", func() {
		It("should return a not found key error", func() {
			emptyFile := yml(assets("examples", "empty.yml"))
			Expect(getToError(emptyFile, "/does-not-exist")).To(
				BeEquivalentTo("failed to traverse tree, expected map but found type string at /"),
			)
		})
	})
})
