package cmd

// Copyright (c) 2018 Bhojpur Consulting Private Limited, India. All rights reserved.

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

import (
	. "github.com/bhojpur/deploy/pkg/entities"
	"github.com/spf13/cobra"
)

var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "applies an entity",
	Args:  cobra.MinimumNArgs(1),
	Long:  `Applies a entity yaml file to your system`,
	RunE: func(cmd *cobra.Command, args []string) error {
		p := &Parser{}

		safe, _ := cmd.Flags().GetBool("safe")

		entity, err := p.ReadEntity(args[0])
		if err != nil {
			return err
		}

		return entity.Apply(entityFile, safe)
	},
}

func init() {
	rootCmd.AddCommand(applyCmd)

	var flags = applyCmd.Flags()
	flags.Bool("safe", false,
		"Avoid to override existing entity if it has difference or if the id is used in a different way.")
}
