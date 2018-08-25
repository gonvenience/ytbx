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
	"os"

	"github.com/HeavyWombat/ytbx/pkg/v1/ytbx"

	"github.com/HeavyWombat/dyff/pkg/v1/neat"

	"github.com/HeavyWombat/dyff/pkg/v1/dyff"

	"github.com/HeavyWombat/dyff/pkg/v1/bunt"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "ytbx",
	Short: "YAML tool box",
	Long: `YAML tool box provides a set of commands to inspect, display, or modify
the content of a given YAML or JSON file.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// exitWithError exits program with given text and error message
func exitWithError(text string, err error) {
	if err != nil {
		fmt.Printf("%s: %s\n", text, bunt.Colorize(err.Error(), bunt.Red))

	} else {
		fmt.Print(text)
	}

	os.Exit(1)
}

func getPathHelp() string {
	exampleData, _ := dyff.LoadYAMLDocuments([]byte(`---
list:
- name: one
  somekey: foobar`))

	neatOutput, _ := neat.ToYAMLString(exampleData[0])

	examplePath := ytbx.Path{PathElements: []ytbx.PathElement{
		{Name: "list"},
		{Key: "name", Name: "one"},
		{Name: "somekey"},
	}}

	return fmt.Sprintf(`
There are two supported path syntaxes available: Dot-style (used by Spruce) or GoPatch-style (used by BOSH).

Example:
---
%s
Path in Dot-style would be: %s
The same path in GoPatch is: %s

`,
		neatOutput,
		examplePath.ToDotStyle(),
		examplePath.ToGoPatchStyle())
}
