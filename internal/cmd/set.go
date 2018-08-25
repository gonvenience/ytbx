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
	"github.com/HeavyWombat/dyff/pkg/v1/dyff"
	"github.com/HeavyWombat/ytbx/internal/pycgo"
	"github.com/HeavyWombat/ytbx/pkg/v1/ytbx"
	"github.com/spf13/cobra"
)

// setCmd represents the set command
var setCmd = &cobra.Command{
	Use:   "set [flags] <file> <path> <value>",
	Args:  cobra.ExactArgs(3),
	Short: "Set the value at a given path",
	Long:  "Set the value at a given path in the file.\n" + getPathHelp(),
	Run: func(cmd *cobra.Command, args []string) {
		location := args[0]
		pathString := args[1]
		newValue := args[2]
		if err := set(location, pathString, newValue); err != nil {
			exitWithError("Failed to set path in file", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(setCmd)
}

// set updates the provided YAML file at the given path. Like Get, the two
// styles of paths are supported. If the path points to a location in the YAML
// file, which does not exist, the required structure is created on the fly.
func set(location string, pathString string, newValue string) error {
	// translate dot style paths into go-patch style path
	if ytbx.IsDotStylePath(pathString) {
		inputfile, err := dyff.LoadFile(location)
		if err != nil {
			return err
		}

		path, err := ytbx.ParseDotStylePathString(pathString, inputfile.Documents[0])
		if err != nil {
			return err
		}

		pathString = path.ToGoPatchStyle()
	}

	return pycgo.UpdateYAML(location, pathString, newValue)
}
