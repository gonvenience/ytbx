// Copyright © 2018 Matthias Diester
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
	"path/filepath"
	"runtime"
	"testing"

	. "github.com/HeavyWombat/dyff/pkg/v1/bunt"
	. "github.com/HeavyWombat/dyff/pkg/v1/dyff"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var assetsDirectory string

func TestYtbx(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "ytbx suite")
}

var _ = BeforeSuite(func() {
	ColorSetting = OFF
	FixedTerminalWidth = 80

	_, file, _, ok := runtime.Caller(0)
	Expect(ok).To(BeTrue())

	dir, err := filepath.Abs(filepath.Dir(file) + "/../../../assets")
	Expect(err).To(BeNil())

	assetsDirectory = dir
})
