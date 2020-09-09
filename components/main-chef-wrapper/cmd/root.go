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
	// rootCmd represents the base command when called without any subcommands
	rootCmd = &cobra.Command{
		Use:   "chef",
		Short: "chef",
		// Uncomment the following line if your bare application
		// has an action associated with it:
		//	Run: func(cmd *cobra.Command, args []string) { },
	}
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// PersistentFlags are also available in child commands.
	rootCmd.PersistentFlags().StringVarP(&options.configFile, "config", "c", "", "Read configuration from `CONFIG_FILE_PATH`")
	// TODO - I think cobra has some support for lists of acceptable param values
	rootCmd.PersistentFlags().StringVar(&options.licenseAcceptance, "chef-license", "",
		"Accept product license, where `ACCEPTANCE` is one of 'accept', 'accept-no-persist', or 'accept-silent'")
	rootCmd.PersistentFlags().BoolVarP(&options.debug, "debug", "d", false,
		"Enable debug output when available")
	rootCmd.PersistentFlags().BoolVarP(&options.debug, "version", "v", false,
		fmt.Sprintf("Show %s version information", dist.WorkstationProduct))
}

func passThroughCommand(targetPath string, cmdName string, args []string) error {
	cmdArg := []string{cmdName}

	allArgs := append(cmdArg, args...)
	cmd := exec.Command(targetPath, allArgs...)
	cmd.Env = os.Environ()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil

	// TODO - verify that the cobra framework will pass along
	//        the error exit code from a called exec.
	//    A: (TODO mp 2020-09-09) Nope, it does not.
	// If we can cast this to an ExitError, then exit with the provided
	// exit code.
	// if exitError, ok := err.(*exec.ExitError); ok {
	//
	//   os.Exit(exitError.ExitCode())
	// }
	//
	// // Otherwise something like 'executable not found' and other
	// // non-exec errors
	// os.Exit(7)

}
