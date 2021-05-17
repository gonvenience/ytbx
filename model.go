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

package ytbx

import (
	"fmt"

	yamlv3 "gopkg.in/yaml.v3"
)

// File represents the actual file (local, or fetched remotely) that needs to
// be processed. It can contain multiple documents, where a document is a map
// or a list of things.
type File struct {
	Location  string
	Note      string
	Documents []*yamlv3.Node
	Names     []string
}

// ReadWriter is able to get and set as well as delete entries in a File
type ReadWriter interface {
	Get(path Path) (interface{}, error)
	Set(path Path, value interface{}) error
	Del(path Path) error

	HasPath(path Path) (bool, error)
}

func (f File) validateDocumentIdx(idx int) error {
	if idx >= 0 && idx < len(f.Documents) {
		return nil
	}

	return fmt.Errorf("document index %d is out of bounds", idx)
}

// Get retrieves the value from the File located at the given Path
func (f File) Get(path Path) (interface{}, error) {
	if err := f.validateDocumentIdx(path.DocumentIdx()); err != nil {
		return nil, err
	}

	result, err := Get(f.Documents[path.DocumentIdx()], path)
	if err != nil {
		return nil, err
	}

	return asType(result)
}

// Set creates or updates the value in File located at the given Path
func (f File) Set(path Path, value interface{}) error {
	if err := f.validateDocumentIdx(path.DocumentIdx()); err != nil {
		return err
	}

	return Set(f.Documents[path.DocumentIdx()], path, value)
}

// Del removes the given Path in the File
func (f File) Del(path Path) error {
	if err := f.validateDocumentIdx(path.DocumentIdx()); err != nil {
		return err
	}

	_, err := Delete(f.Documents[path.DocumentIdx()], path)
	return err
}

// HasPath checks whether the given Path is available in the File
func (f File) HasPath(path Path) (bool, error) {
	if err := f.validateDocumentIdx(path.DocumentIdx()); err != nil {
		return false, err
	}

	_, err := Get(f.Documents[path.DocumentIdx()], path)
	return err == nil, nil
}
