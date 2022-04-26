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
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/bhojpur/deploy/pkg/console"
	"github.com/bhojpur/deploy/pkg/executor"
	"github.com/bhojpur/deploy/pkg/logger"
	"github.com/bhojpur/deploy/pkg/schema"
	"github.com/bhojpur/deploy/pkg/version"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/twpayne/go-vfs"
)

func initLogger() logger.Interface {
	ll := log.New()
	switch strings.ToLower(os.Getenv("LOGLEVEL")) {
	case "error":
		ll.SetLevel(log.ErrorLevel)
	case "warning":
		ll.SetLevel(log.WarnLevel)
	case "debug":
		ll.SetLevel(log.DebugLevel)
	default:
		ll.SetLevel(log.InfoLevel)
	}
	return ll
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "depcfg",
	Short:   "A system configurator for Unix-like environments",
	Version: fmt.Sprintf("%s-g%s %s", version.Version, version.BuildCommit, version.BuildTime),
	Long: `Bhojpur Deploy loads cloud-init style yamls and applies them in the system.

For example:
	$> depcfg -s initramfs https://<deploy.yaml> /path/to/disk <definition.yaml> ...
	$> depcfg -s initramfs <deploy.yaml> <deploy2.yaml> ...
	$> depcfg def.yaml | depcfg -
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		stage, _ := cmd.Flags().GetString("stage")
		dot, _ := cmd.Flags().GetBool("dotnotation")

		ll := initLogger()
		runner := executor.NewExecutor(executor.WithLogger(ll))
		fromStdin := len(args) == 1 && args[0] == "-"

		ll.Infof("Bhojpur Deploy configure version %s", cmd.Version)
		if len(args) == 0 {
			ll.Fatal("depcfg needs at least one path or URL as argument")
		}
		stdConsole := console.NewStandardConsole(console.WithLogger(ll))

		if dot {
			runner.Modifier(schema.DotNotationModifier)
		}

		if fromStdin {
			std, err := ioutil.ReadAll(os.Stdin)
			if err != nil {
				return err
			}

			args = []string{string(std)}
		}

		return runner.Run(stage, vfs.OSFS, stdConsole, args...)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}

func init() {
	rootCmd.PersistentFlags().StringP("stage", "s", "default", "Stage to apply")
	rootCmd.PersistentFlags().BoolP("dotnotation", "d", false, "Parse input in dotnotation ( e.g. `stages.foo.name=..` ) ")
}
