//
// Copyright 2020 Chef Software, Inc.
// Author: Marc A. Paradise <marc.paradise@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package cmd

import (
	"errors"
	"fmt"
	"os"
	"os/exec"

	"github.com/chef/chef-workstation/components/main-chef-wrapper/dist"
	"github.com/spf13/cobra"
)

// Cobra usage notes

// 1. describe flag value parameter with backticks. This behavior is
// not documented clearly in the cobra or pflags libs, but backtick in
// the description allows us to the name the flag's parameter in help text.
//
// For example, this description:
//   "Read configuration from this path"
// Will give this help output:
//   -c, --config string   Read configuration from this path
//
// But using this one:
//   "Read configuration from `CONFIG_FILE_PATH`:
// Gives us:
//   -c, --config CONFIG_FILE_PATH  Read configuration from CONFIG_FILE_PATH
//
// The latter is more clear to the operator, so we prefer it.

type rootConfig struct {
	configFile        string
	licenseAcceptance string
	debug             bool
}

var (
	options rootConfig

	rootCmd = &cobra.Command{
		Use:   "chef",
		Short: "chef",
		// This flags means that no error info will be output by default when
		// a command fails.  This is good in the most common supported case -
		// dispatching to a command handled by another binary. In that case, we
		// rely on the binary to report its own errors, and don't want to show an
		// extra message if it exits non-zero.
		// For commands implemented internally, we will want to explicitly set
		// SilenceErrors: false on those commands, so that those errors are displayed.
		SilenceErrors: true,

		// Arg validation
		// Args: func(cmd *cobra.Command, args []string) {}
		// RunE: func(cmd *cobra.Command, args []string) { } error,
	}
)

func Execute() {
	var ee *exec.ExitError
	if err := rootCmd.Execute(); err != nil {
		if errors.As(err, &ee) {
			os.Exit(ee.ExitCode())
		}
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// TODO is there a a way to have these _only_ in child commands, without
	//      having them visible in root command? This would avoid us having to implement
	//      license handling prematurely in case someone wants to `chef --chef-license accept`

	// These flags are common to all child commands.  Some of them do not need config or debug,
	// so we can look at pushing this down; but it seems to make sense since it's present for more
	// commands than it isn't.
	rootCmd.PersistentFlags().StringVarP(&options.configFile, "config", "c", "", "Read configuration from `CONFIG_FILE_PATH`")
	rootCmd.PersistentFlags().StringVar(&options.licenseAcceptance, "chef-license", "",
		"Accept product license, where `ACCEPTANCE` is one of 'accept', 'accept-no-persist', or 'accept-silent'")
	rootCmd.PersistentFlags().BoolVarP(&options.debug, "debug", "d", false,
		"Enable debug output when available")
	rootCmd.PersistentFlags().BoolVarP(&options.debug, "version", "v", false,
		fmt.Sprintf("Show %s version information", dist.WorkstationProduct))
}

// TODO -
func passThroughCommand(targetPath string, cmdName string, args []string) error {

	var allArgs []string
	if cmdName != "" {
		allArgs = append([]string{cmdName}, args...)
	} else {
		allArgs = args
	}

	//
	cmd := exec.Command(targetPath, allArgs...)
	cmd.Env = os.Environ()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	return cmd.Run()

}
