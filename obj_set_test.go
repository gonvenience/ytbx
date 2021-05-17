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

package ytbx_test

import (
	"time"

	. "github.com/gonvenience/ytbx"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"gopkg.in/yaml.v3"
)

var _ = Describe("writing entries", func() {
	Context("writing using File struct", func() {
		var file *File

		BeforeEach(func() {
			file = &File{
				Location: "/foo/bar",
				Documents: []*yaml.Node{
					yml(`foo: bar`),
					yml(`bar: foo`),
				},
			}
		})

		AfterEach(func() {
			file = nil
		})

		It("should create an entry at document root level", func() {
			var path = parseGoPatch("1:/new")

			Expect(file.HasPath(path)).To(BeFalse())
			Expect(file.Set(path, "value")).To(BeNil())
			Expect(file.Get(path)).To(Equal("value"))
		})

		It("should create an entry in an intermediate map", func() {
			var path = parseGoPatch("/some/nested/key")

			Expect(file.HasPath(path)).To(BeFalse())
			Expect(file.Set(path, "value")).To(BeNil())
			Expect(file.Get(path)).To(Equal("value"))
		})

		It("should create an entry in an intermediate named entry list", func() {
			var path = parseGoPatch("/list/name=one/key")

			Expect(file.HasPath(path)).To(BeFalse())
			Expect(file.Set(path, "value")).To(BeNil())
			Expect(file.Get(path)).To(Equal("value"))
		})

		It("should create an entry in an intermediate simple list", func() {
			var path = parseGoPatch("/list/-1/key")

			Expect(file.HasPath(path)).To(BeFalse())
			Expect(file.Set(path, "value")).To(BeNil())
			Expect(file.Get(path)).To(Equal("value"))
		})

		It("should create an entry in a new simple list", func() {
			var path = parseGoPatch("/list/-1")

			Expect(file.HasPath(path)).To(BeFalse())
			Expect(file.Set(path, "value")).To(BeNil())
			Expect(file.Get(path)).To(Equal("value"))
		})

		It("should create an entry that is a map", func() {
			var path = parseGoPatch("/new")

			mapping := map[interface{}]interface{}{
				"key": "value",
				"foo": "bar",
			}

			Expect(file.HasPath(path)).To(BeFalse())
			Expect(file.Set(path, mapping)).To(BeNil())
			Expect(file.Get(path)).To(Equal(mapping))
		})

		It("should create an entry that is simple list", func() {
			var path = parseGoPatch("/new")

			listing := []interface{}{
				"one",
				"two",
				"three",
			}

			Expect(file.HasPath(path)).To(BeFalse())
			Expect(file.Set(path, listing)).To(BeNil())
			Expect(file.Get(path)).To(Equal(listing))
		})

		It("should create an entry that is named entry list", func() {
			var path = parseGoPatch("/new")

			listing := []interface{}{
				map[interface{}]interface{}{"name": "one"},
				map[interface{}]interface{}{"name": "two"},
				map[interface{}]interface{}{"name": "three"},
			}

			Expect(file.HasPath(path)).To(BeFalse())
			Expect(file.Set(path, listing)).To(BeNil())
			Expect(file.Get(path)).To(Equal(listing))
		})

		It("should create all YAML spec types", func() {
			var path = parseGoPatch("/new")

			mapping := map[interface{}]interface{}{
				"int":    1337,
				"float":  13.37,
				"string": "foobar",
				"bool":   true,
				"time":   time.Date(1955, time.November, 5, 9, 00, 00, 0, time.UTC),
				"null":   nil,
			}

			Expect(file.HasPath(path)).To(BeFalse())
			Expect(file.Set(path, mapping)).To(BeNil())
			Expect(file.Get(path)).To(BeEquivalentTo(mapping))
		})
	})
})
