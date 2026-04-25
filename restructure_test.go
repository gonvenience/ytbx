// Copyright © 2019 The Homeport Team
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
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	. "github.com/gonvenience/ytbx"

	"go.yaml.in/yaml/v3"
)

var _ = Describe("Restructure order of map keys", func() {
	var example *yaml.Node

	var keys []string
	var err error

	Context("with an input example", func() {
		JustBeforeEach(func() {
			keys, err = ListStringKeys(example)
		})

		Context("that is a Concourse like schema", func() {
			Context("at document root level", func() {
				BeforeEach(func() {
					example = yml("{ groups: [], jobs: [], resources: [], resource_types: [] }")
					RestructureObject(example)
				})

				It("should not error", func() {
					Expect(err).ToNot(HaveOccurred())
				})

				It("should have the expected order of fields", func() {
					Expect(keys).To(BeEquivalentTo([]string{"jobs", "resources", "resource_types", "groups"}))
				})
			})

			Context("at task level", func() {
				BeforeEach(func() {
					example = yml("{ source: {}, name: {}, type: {}, privileged: {} }")
					RestructureObject(example)
				})

				It("should not error", func() {
					Expect(err).ToNot(HaveOccurred())
				})

				It("should have the expected order of fields", func() {
					Expect(keys).To(BeEquivalentTo([]string{"name", "type", "source", "privileged"}))
				})
			})

			Context("at document root level inside the resources list", func() {
				BeforeEach(func() {
					example = yml("{ resources: [ { privileged: false, source: { branch: foo, paths: [] }, name: myname, type: mytype } ] }")
					RestructureObject(example)
				})

				JustBeforeEach(func() {
					keys, err = ListStringKeys(example.Content[1].Content[0])
				})

				It("should not error", func() {
					Expect(err).ToNot(HaveOccurred())
				})

				It("should have the expected order of fields", func() {
					Expect(keys).To(BeEquivalentTo([]string{"name", "type", "source", "privileged"}))
				})
			})
		})

		Context("that has no particular schema", func() {
			BeforeEach(func() {
				example = yml(`{"list":["one","two","three"], "some":{"deep":{"structure":{"where":{"you":{"loose":{"focus":{"one":1,"two":2}}}}}}}, "name":"here", "release":"this"}`)
				RestructureObject(example)
			})

			It("should not error", func() {
				Expect(err).ToNot(HaveOccurred())
			})

			It("should have the expected order of fields", func() {
				Expect(keys).To(BeEquivalentTo([]string{"name", "release", "list", "some"}))
			})
		})

		Context("that is Kubernetes like schema", func() {
			Context("e.g. kustomization.yaml file", func() {
				BeforeEach(func() {
					example = yml(assets("kustomize/kustomization.yaml"))
					RestructureObject(example)
				})

				JustBeforeEach(func() {
					keys, err = ListStringKeys(example.Content[0])
				})

				It("should not error", func() {
					Expect(err).ToNot(HaveOccurred())
				})

				It("should have the expected order of fields", func() {
					Expect(keys).To(Equal([]string{"apiVersion", "kind", "resources", "configMapGenerator"}))
				})
			})
		})
	})
})
