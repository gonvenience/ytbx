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

package ytbx_test

import (
	. "github.com/gonvenience/ytbx"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("deleting entries", func() {
	Context("deleting using File struct", func() {
		var file *File

		BeforeEach(func() {
			file = loadFile(assets("examples", "types.yml"))
		})

		AfterEach(func() {
			file = nil
		})

		It("should delete an entry in a map referenced by the path", func() {
			var path = parseGoPatch("/yaml/map/before")

			Expect(file.HasPath(path)).To(BeTrue())
			Expect(file.Del(path)).To(BeNil())
			Expect(file.HasPath(path)).To(BeFalse())
		})

		It("should delete an entry in a simple list referenced by the path", func() {
			var path = parseGoPatch("/yaml/simple-list/1")

			Expect(file.HasPath(path)).To(BeTrue())
			Expect(file.Del(path)).To(BeNil())

			list, err := GetPath(file.Documents[0], "/yaml/simple-list")
			Expect(err).ToNot(HaveOccurred())
			Expect(len(list.Content)).To(Equal(4))
		})

		It("should delete an entry in a named entry list referenced by the path", func() {
			var path = parseGoPatch("/yaml/named-entry-list-using-name/name=C")

			Expect(file.HasPath(path)).To(BeTrue())
			Expect(file.Del(path)).To(BeNil())

			list, err := GetPath(file.Documents[0], "/yaml/named-entry-list-using-name")
			Expect(err).ToNot(HaveOccurred())
			Expect(len(list.Content)).To(Equal(4))
		})
	})
})
