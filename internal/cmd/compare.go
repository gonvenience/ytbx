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

	"github.com/HeavyWombat/ytbx/pkg/v1/ytbx"
	"github.com/spf13/cobra"
)

var comparePathsByValue bool

// compareCmd represents the compare command
var compareCmd = &cobra.Command{
	Use:   "compare <from_file> <to_file>",
	Args:  cobra.RangeArgs(1, 2),
	Short: "Compare YAML paths",
	Long:  `Compare YAML paths between two files.`,
	Run: func(cmd *cobra.Command, args []string) {
		list, err := ytbx.ComparePaths(args[0], args[1], ytbx.GoPatchStyle, comparePathsByValue)
		if err != nil {
			exitWithError("Failed to get paths from file", err)
		}
		for _, entry := range list {
			fmt.Println(entry)
		}
	},
}

func init() {
	rootCmd.AddCommand(compareCmd)
	compareCmd.PersistentFlags().BoolVar(&comparePathsByValue, "by-value", false, "compare only paths with same value")
}
