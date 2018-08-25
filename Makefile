# Copyright © 2018 Matthias Diester
#
# Permission is hereby granted, free of charge, to any person obtaining a copy
# of this software and associated documentation files (the "Software"), to deal
# in the Software without restriction, including without limitation the rights
# to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
# copies of the Software, and to permit persons to whom the Software is
# furnished to do so, subject to the following conditions:
#
# The above copyright notice and this permission notice shall be included in
# all copies or substantial portions of the Software.
#
# THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
# IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
# FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
# AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
# LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
# OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
# THE SOFTWARE.

.PHONY: all clean test pythontest gotest build upx

os := $(shell uname | tr '[:upper:]' '[:lower:]')
arch := $(shell uname -m | sed 's/x86_64/amd64/')

all: test

clean:
	go clean -i -r -cache
	rm -rf internal/pycgo/updateYAML.c internal/pycgo/updateYAML.go internal/pycgo/__pycache__ binaries

internal/pycgo/updateYAML.c: internal/pycgo/updateYAML.py
	cython -3 --embed=updateYAML --output-file internal/pycgo/updateYAML.c internal/pycgo/updateYAML.py

internal/pycgo/updateYAML.go: internal/pycgo/updateYAML.go.template
	@scripts/createGoSourceFileFromTemplate.sh

test: pythontest gotest

pythontest:
	python3 internal/pycgo/updateYAML_test.py

gotest:
	ginkgo -r --nodes 1 --randomizeAllSpecs --randomizeSuites --race --trace

build: internal/pycgo/updateYAML.c internal/pycgo/updateYAML.go
	@mkdir -p binaries
	go build -a -ldflags='-s -w' -o binaries/ytbx-$(os)-$(arch) cmd/ytbx/main.go

upx: build
	upx -q binaries/ytbx-$(os)-$(arch)
