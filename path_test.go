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
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/gonvenience/ytbx"
)

var _ = Describe("path tests", func() {
	Context("parse dot-style path strings into a path", func() {
		It("should parse string with only map elements", func() {
			path, err := ParseDotStylePathString("yaml.structure.somekey")
			Expect(err).To(BeNil())
			Expect(path.DotStyle()).To(Equal("yaml.structure.somekey"))
			Expect(path.GoPatchStyle()).To(Equal("/yaml/structure/somekey"))
		})

		It("should parse string with simple list entry", func() {
			path, err := ParseDotStylePathString("simpleList.1")
			Expect(err).To(BeNil())
			Expect(path.DotStyle()).To(Equal("simpleList.1"))
			Expect(path.GoPatchStyle()).To(Equal("/simpleList/1"))
		})
	})

	Context("parse go-patch style path strings into paths", func() {
		It("should parse an input string using go-patch style into a path (only maps)", func() {
			path, err := ParseGoPatchStylePathString("/yaml/structure/somekey")
			Expect(err).To(BeNil())
			Expect(path.DotStyle()).To(Equal("yaml.structure.somekey"))
			Expect(path.GoPatchStyle()).To(Equal("/yaml/structure/somekey"))
		})

		It("should parse an input string using go-patch style into a path (maps and named-entry lists)", func() {
			path, err := ParseGoPatchStylePathString("/list/name=one/somekey")
			Expect(err).To(BeNil())
			Expect(path.DotStyle()).To(Equal("list.one.somekey"))
			Expect(path.GoPatchStyle()).To(Equal("/list/name=one/somekey"))
		})

		It("should parse an input string using go-patch style into a path (simple list)", func() {
			path, err := ParseGoPatchStylePathString("/simpleList/1")
			Expect(err).To(BeNil())
			Expect(path.DotStyle()).To(Equal("simpleList.1"))
			Expect(path.GoPatchStyle()).To(Equal("/simpleList/1"))
		})

		It("should parse an input string that points to the root of the tree structure", func() {
			path, err := ParseGoPatchStylePathString("/")
			Expect(err).To(BeNil())
			Expect(path.DotStyle()).To(Equal("(root)"))
			Expect(path.GoPatchStyle()).To(Equal("/"))
		})

		It("should parse real-life scenario paths with mixed types", func() {
			path, err := ParseGoPatchStylePathString("/resource_pools/name=concourse_resource_pool/cloud_properties/datacenters/0/clusters")
			Expect(err).ToNot(HaveOccurred())
			Expect(path.DotStyle()).To(Equal("resource_pools.concourse_resource_pool.cloud_properties.datacenters.0.clusters"))
			Expect(path.GoPatchStyle()).To(Equal("/resource_pools/name=concourse_resource_pool/cloud_properties/datacenters/0/clusters"))
		})

		It("should parse path strings with escaped slashes", func() {
			path, err := ParseGoPatchStylePathString("/foo/name=bar.com\\/id/string")
			Expect(err).ToNot(HaveOccurred())
			Expect(path.DotStyle()).To(Equal("foo.bar.com/id.string"))
			Expect(path.GoPatchStyle()).To(Equal("/foo/name=bar.com/id/string"))
		})

		It("should parse an input string using non-standard go-patch style with document index", func() {
			path, err := ParseGoPatchStylePathString("1:/yaml")
			Expect(err).To(BeNil())
			Expect(path.DotStyle()).To(Equal("yaml"))
			Expect(path.GoPatchStyle()).To(Equal("/yaml"))
		})
	})

	Context("compare paths between two files", func() {
		It("should find only duplicate paths", func() {
			list, err := ComparePaths(
				assets("testbed", "sample_a.yml"),
				assets("testbed", "sample_b.yml"),
				false,
			)

			Expect(err).ToNot(HaveOccurred())
			Expect(list).To(Equal([]Path{
				parseGoPatch("/yaml/structure/somekey"),
				parseGoPatch("/yaml/structure/dot"),
				parseGoPatch("/list/name=sametwo/somekey"),
			}))
		})

		It("should find only paths with the same value", func() {
			list, err := ComparePaths(
				assets("testbed", "sample_a.yml"),
				assets("testbed", "sample_b.yml"),
				true,
			)

			Expect(err).ToNot(HaveOccurred())
			Expect(list).To(Equal([]Path{
				parseGoPatch("/yaml/structure/dot"),
				parseGoPatch("/list/name=sametwo/somekey"),
			}))
		})
	})

	Context("checking for path in YAML", func() {
		It("should check whether the provided path is in the YAML", func() {
			example := yml(assets("examples", "types.yml"))

			Expect(IsPathInTree(example, "/yaml/map/before")).To(BeTrue())
			Expect(IsPathInTree(example, "/yaml/map/nope")).To(BeFalse())

			Expect(IsPathInTree(example, "/yaml/simple-list/0")).To(BeTrue())
			Expect(IsPathInTree(example, "/yaml/simple-list/5")).To(BeFalse())

			Expect(IsPathInTree(example, "/yaml/named-entry-list-using-name/name=A/foo")).To(BeTrue())
			Expect(IsPathInTree(example, "/yaml/named-entry-list-using-name/name=nope/foo")).To(BeFalse())
		})
	})
})
