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

package cmd

import (
	"fmt"

	"github.com/HeavyWombat/dyff/pkg/v1/dyff"
	"github.com/HeavyWombat/gonvenience/pkg/v1/neat"
	"github.com/spf13/cobra"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Args:  cobra.ExactArgs(2),
	Short: "Get the value at a given path",
	Long:  "Get the value at a given path in the file.\n" + getPathHelp(),
	Run: func(cmd *cobra.Command, args []string) {
		location := args[0]
		pathString := args[1]

		obj, err := get(location, pathString)
		if err != nil {
			exitWithError("Failed to get path from file", err)
		}

		boldKeys := true
		useIndentLines := true
		outoutProcessor := neat.NewOutputProcessor(useIndentLines, boldKeys, &neat.DefaultColorSchema)

		output, err := outoutProcessor.ToYAML(obj)
		if err != nil {
			exitWithError("Failed to render gathered data", err)
		}

		fmt.Print(output)
	},
}

func init() {
	rootCmd.AddCommand(getCmd)
}

// get returns the data that is stored in the provided YAML file at the given
// path. There are two types of path supported: A dot-style path
// (yaml.structure.one) and GoPatch-style paths (/yaml/structure/name=one).
func get(location string, pathString string) (interface{}, error) {
	inputfile, err := dyff.LoadFile(location)
	if err != nil {
		return nil, err
	}

	content, err := dyff.Grab(inputfile.Documents[0], pathString)
	if err != nil {
		return nil, err
	}

	return content, nil
}
