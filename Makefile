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

.PHONY: all clean test gobin build upx

version := $(shell git describe --tags --abbrev=0 2>/dev/null || (git rev-parse HEAD | cut -c-8))
sources := $(wildcard cmd/ytbx/*.go internal/cmd/*.go pkg/v1/ytbx/*.go)

all: test

clean:
	GO111MODULE=on go clean -i -r -cache
	rm -rf binaries

test:
	GO111MODULE=on ginkgo -r --randomizeAllSpecs --randomizeSuites --failOnPending --trace --race --nodes=2 --compilers=2

gobin:
	GO111MODULE=on go build -ldflags='-s -w -extldflags "-static"' -o ${GOPATH}/bin/ytbx cmd/ytbx/main.go

build: binaries/ytbx-linux-amd64 binaries/ytbx-darwin-amd64 binaries/ytbx-windows-amd64

binaries/ytbx-linux-amd64: $(sources)
	GO111MODULE=on CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -tags netgo -ldflags='-s -w -extldflags "-static" -X github.com/HeavyWombat/ytbx/internal/cmd.version=$(version)' -o binaries/ytbx-linux-amd64 cmd/ytbx/main.go

binaries/ytbx-darwin-amd64: $(sources)
	GO111MODULE=on CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -tags netgo -ldflags='-s -w -extldflags "-static" -X github.com/HeavyWombat/ytbx/internal/cmd.version=$(version)' -o binaries/ytbx-darwin-amd64 cmd/ytbx/main.go

binaries/ytbx-windows-amd64: $(sources)
	GO111MODULE=on CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -tags netgo -ldflags='-s -w -extldflags "-static" -X github.com/HeavyWombat/ytbx/internal/cmd.version=$(version)' -o binaries/ytbx-windows-amd64 cmd/ytbx/main.go
