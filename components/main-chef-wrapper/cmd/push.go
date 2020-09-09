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
// LIMITATIONS UNDER THE LICENSE.
//

package cmd

import (
	"fmt"

	"github.com/chef/chef-workstation/components/main-chef-wrapper/dist"
	"github.com/spf13/cobra"
)

var pushCmd = &cobra.Command{
	Use:   "push POLICY_GROUP [ POLICY_FILE ]",
	Short: "Push a local policyfile lock to a policy group on the %s",
	Long: `Upload an existing Policyfile.lock.json to a %s, along
with all the cookbooks contained in the policy lock. The policy lock is applied
to a specific POLICY_GROUP, which is a set of nodes that share the same
run_list and cookbooks.

See the Policyfile documentation for more information:

https://docs.chef.io/policyfile/
`,

	Run: func(cmd *cobra.Command, args []string) {
		passThroughCommand("chef-cli", "push", args)
	},
}

func init() {
	pushCmd.Short = fmt.Sprintf(pushCmd.Short, dist.ServerProduct)
	pushCmd.Long = fmt.Sprintf(pushCmd.Long, dist.ServerProduct)
	rootCmd.AddCommand(pushCmd)
}
