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

package ytbx

import (
	"bytes"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/gonvenience/bunt"
	yamlv3 "gopkg.in/yaml.v3"
)

var goPatchRegEx = regexp.MustCompile(`^((\d+):)?(/.*)$`)
var dotRegEx = regexp.MustCompile(`^((\d+):)?(.*)$`)

// PathStyle is a custom type for supported path styles
type PathStyle int

// Supported styles are the Dot-Style (used by Spruce for example) and GoPatch
// Style which is used by BOSH
const (
	DotStyle PathStyle = iota
	GoPatchStyle
)

type pathSectionType int

const (
	mappingEntry pathSectionType = iota
	namedListEntry
	indexedListEntry
	undefinedEntry
)

// Finder describes the general interface for Path structures
type Finder interface {
	DocumentIdx() int
	RootDescription() string
	GoPatchStyle() string
	DotStyle() string
}

type section interface {
	sectionType() pathSectionType
	goPatchStyle() string
	dotStyle() string
}

// Path points to a section in a data structure by using names to identify the
// location.
// Example:
//   ---
//   sizing:
//     api:
//       count: 2
// For example, `sizing.api.count` points to the key `sizing` of the root
// element and in there to the key `api` and so on and so forth.
type Path struct {
	root     *File
	docIdx   int
	sections []section
}

// TODO: remove me
var _ Finder = &Path{}

func newPath(path Path, section section) Path {
	return Path{
		root:     path.root,
		docIdx:   path.docIdx,
		sections: append(path.sections, section),
	}
}

// NewPathWithNamedEntryListSection creates a new path based on the provided
// path by adding a new named entry list section
func NewPathWithNamedEntryListSection(path Path, key string, name string) Path {
	return newPath(path, listNamedSection{id: key, name: name})
}

// NewPathWithNamedEntrySection creates a new path based on the provided path
// by adding a new named entry section
func NewPathWithNamedEntrySection(path Path, name string) Path {
	return newPath(path, mappingNameSection{name: name})
}

// NewPathWithIndexedEntrySection creates a new path based on the provided path
// by adding a new indexed entry section
func NewPathWithIndexedEntrySection(path Path, idx int) Path {
	return newPath(path, listIdxSection{idx: idx})
}

func (p Path) String() string {
	return p.GoPatchStyle()
}

// GoPatchStyle returns the path as a GoPatch style string.
func (p *Path) GoPatchStyle() string {
	if len(p.sections) == 0 {
		return "/"
	}

	var buf bytes.Buffer
	for _, section := range p.sections {
		buf.WriteString("/")
		buf.WriteString(section.goPatchStyle())
	}

	buf.WriteString(p.optionalDocIdxString())

	return buf.String()
}

// DotStyle returns the path as a Dot-Style string.
func (p *Path) DotStyle() string {
	if len(p.sections) == 0 {
		return "(root)"
	}

	var result = make([]string, len(p.sections))
	for i, section := range p.sections {
		result[i] = section.dotStyle()
	}

	return strings.Join(result, ".") + p.optionalDocIdxString()
}

// DocumentIdx returns the document index (document in the file)
func (p *Path) DocumentIdx() int {
	return p.docIdx
}

// RootDescription returns a description of the root level of this path, which
// could be the number of the respective document inside a YAML or if available
// the name of the document
func (p *Path) RootDescription() string {
	if p.root != nil && p.docIdx < len(p.root.Names) {
		return p.root.Names[p.docIdx]
	}

	// Note: human style counting that starts with 1
	return fmt.Sprintf("document #%d", p.docIdx+1)
}

// Parent returns the parent of the provided Path
func (p *Path) Parent() (Path, error) {
	if len(p.sections) == 0 {
		return Path{}, fmt.Errorf("path %s does not have a parent", p)
	}

	return Path{
		docIdx:   p.docIdx,
		sections: p.sections[:len(p.sections)-1],
	}, nil
}

func (p *Path) optionalDocIdxString() string {
	if p.root != nil && len(p.root.Documents) > 1 {
		return fmt.Sprintf("  (document #%d)", p.docIdx+1)
	}

	return ""
}

type listIdxSection struct{ idx int }

func (s listIdxSection) sectionType() pathSectionType { return indexedListEntry }
func (s listIdxSection) goPatchStyle() string         { return fmt.Sprintf("%d", s.idx) }
func (s listIdxSection) dotStyle() string             { return fmt.Sprintf("%d", s.idx) }

type listNamedSection struct{ id, name string }

func (s listNamedSection) sectionType() pathSectionType { return namedListEntry }
func (s listNamedSection) goPatchStyle() string         { return bunt.Sprintf("_%s_=%s", s.id, s.name) }
func (s listNamedSection) dotStyle() string             { return bunt.Sprintf("_%s_", s.name) }

type mappingNameSection struct{ name string }

func (s mappingNameSection) sectionType() pathSectionType { return mappingEntry }
func (s mappingNameSection) goPatchStyle() string         { return s.name }
func (s mappingNameSection) dotStyle() string             { return s.name }

type undefSection struct{ raw string }

func (s undefSection) sectionType() pathSectionType { return undefinedEntry }
func (s undefSection) goPatchStyle() string         { return s.raw }
func (s undefSection) dotStyle() string             { return s.raw }

// ComparePathsByValue returns all Path structure that have the same path value
func ComparePathsByValue(fromLocation string, toLocation string, duplicatePaths []Path) ([]Path, error) {
	from, err := LoadFile(fromLocation)
	if err != nil {
		return nil, err
	}

	to, err := LoadFile(toLocation)
	if err != nil {
		return nil, err
	}

	if len(from.Documents) > 1 || len(to.Documents) > 1 {
		return nil, fmt.Errorf(
			"input files have more than one document, which is not supported yet",
		)
	}

	var duplicatePathsWithTheSameValue []Path
	for _, path := range duplicatePaths {
		fromValue, err := GetPath(from.Documents[0], path.GoPatchStyle())
		if err != nil {
			return nil, err
		}

		toValue, err := GetPath(to.Documents[0], path.GoPatchStyle())
		if err != nil {
			return nil, err
		}

		if reflect.DeepEqual(fromValue, toValue) {
			duplicatePathsWithTheSameValue = append(duplicatePathsWithTheSameValue, path)
		}
	}
	return duplicatePathsWithTheSameValue, nil
}

// ComparePaths returns all duplicate Path structures between two documents.
func ComparePaths(fromLocation string, toLocation string, compareByValue bool) ([]Path, error) {
	var duplicatePaths []Path

	pathsFromLocation, err := ListPaths(fromLocation)
	if err != nil {
		return nil, err
	}
	pathsToLocation, err := ListPaths(toLocation)
	if err != nil {
		return nil, err
	}

	lookup := map[string]struct{}{}
	for _, pathsFrom := range pathsFromLocation {
		lookup[pathsFrom.GoPatchStyle()] = struct{}{}
	}

	for _, pathsTo := range pathsToLocation {
		if _, ok := lookup[pathsTo.GoPatchStyle()]; ok {
			duplicatePaths = append(duplicatePaths, pathsTo)
		}
	}

	if !compareByValue {
		return duplicatePaths, nil
	}

	return ComparePathsByValue(fromLocation, toLocation, duplicatePaths)
}

// ListPaths returns all paths in the documents using the provided choice of
// path style.
func ListPaths(location string) ([]Path, error) {
	inputfile, err := LoadFile(location)
	if err != nil {
		return nil, err
	}

	var paths []Path
	for idx, document := range inputfile.Documents {
		root := Path{docIdx: idx}

		traverseTree(root, nil, document, func(path Path, _ *yamlv3.Node, _ *yamlv3.Node) {
			paths = append(paths, path)
		})
	}

	return paths, nil
}

// IsPathInTree returns whether the provided path is in the given YAML structure
func IsPathInTree(tree *yamlv3.Node, pathString string) (bool, error) {
	searchPath, err := ParsePathString(pathString)
	if err != nil {
		return false, err
	}

	resultChan := make(chan bool)

	go func() {
		for _, node := range tree.Content {
			traverseTree(Path{}, nil, node, func(path Path, _ *yamlv3.Node, _ *yamlv3.Node) {
				if path.GoPatchStyle() == searchPath.GoPatchStyle() {
					resultChan <- true
				}
			})

			resultChan <- false
		}
	}()

	return <-resultChan, nil
}

// ParsePathString returns a path by parsing a string representation
// of a path, which can be one of the supported types.
func ParsePathString(pathString string) (*Path, error) {
	if goPatchRegEx.MatchString(pathString) {
		return ParseGoPatchStylePathString(pathString)
	}

	return ParseDotStylePathString(pathString)
}

// ParseGoPatchStylePathString returns a path by parsing a string representation
// which is assumed to be a GoPatch style path.
func ParseGoPatchStylePathString(path string) (*Path, error) {
	matches := goPatchRegEx.FindStringSubmatch(path)
	if matches == nil {
		return nil, NewInvalidPathError(GoPatchStyle, path,
			"failed to parse path string, because path does not match expected format",
		)
	}

	var documentIdx int
	if len(matches[2]) > 0 {
		var err error
		documentIdx, err = strconv.Atoi(matches[2])
		if err != nil {
			return nil, NewInvalidPathError(GoPatchStyle, path,
				"failed to parse path string, because path does not match expected format",
			)
		}
	}

	// Reset path variable to only contain the raw path string
	path = matches[3]

	// Special case for root path
	if path == "/" {
		return &Path{docIdx: documentIdx}, nil
	}

	// Hacky solution to deal with escaped slashes, replace them with a "safe"
	// replacement string that is later resolved into a simple slash
	path = strings.Replace(path, `\/`, `%2F`, -1)

	var elements []section
	for i, section := range strings.Split(path, "/") {
		if i == 0 {
			continue
		}

		keyNameSplit := strings.SplitN(section, "=", 2)
		switch len(keyNameSplit) {
		case 1:
			if idx, err := strconv.Atoi(keyNameSplit[0]); err == nil {
				elements = append(elements, listIdxSection{
					idx: idx,
				})

			} else {
				elements = append(elements, mappingNameSection{
					name: strings.Replace(keyNameSplit[0], `%2F`, "/", -1),
				})
			}

		case 2:
			elements = append(elements, listNamedSection{
				id:   strings.Replace(keyNameSplit[0], `%2F`, "/", -1),
				name: strings.Replace(keyNameSplit[1], `%2F`, "/", -1),
			})
		}
	}

	return &Path{docIdx: documentIdx, sections: elements}, nil
}

// ParseDotStylePathString returns a path by parsing a string representation
// which is assumed to be a Dot-Style path.
func ParseDotStylePathString(path string) (*Path, error) {
	matches := dotRegEx.FindStringSubmatch(path)
	if matches == nil {
		return nil, NewInvalidPathError(GoPatchStyle, path,
			"failed to parse path string, because path does not match expected format",
		)
	}

	var documentIdx int
	if len(matches[2]) > 0 {
		var err error
		documentIdx, err = strconv.Atoi(matches[2])
		if err != nil {
			return nil, NewInvalidPathError(GoPatchStyle, path,
				"failed to parse path string, cannot parse document index: %s", matches[2],
			)
		}
	}

	// Reset path variable to only contain the raw path string
	path = matches[3]

	var elements []section
	for _, section := range strings.Split(path, ".") {
		if idx, err := strconv.Atoi(section); err == nil {
			elements = append(elements, listIdxSection{idx})

		} else {
			elements = append(elements, undefSection{section})
		}
	}

	return &Path{docIdx: documentIdx, sections: elements}, nil
}
