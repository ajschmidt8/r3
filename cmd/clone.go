/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"github.com/ajschmidt8/rrr/shared"
	"github.com/spf13/cobra"
)

var createBranch bool

// cloneCmd represents the clone command
var cloneCmd = &cobra.Command{
	Use:   "clone",
	Short: "Clone repositories",
	Long: `Clone each repository listed in config.yaml.

By default, the "pr.base_branch" value will be checked out in each repository.
This command will create a fork of the desired repo if one does not exist already.`,
	Run: func(cmd *cobra.Command, args []string) {
		config := shared.ReadConfig()

		for _, repoName := range config.Repos {
			newBranchName := ""
			if createBranch {
				newBranchName = config.BranchName
			}
			shared.Clone(repoName, config.PR.BaseBranch, newBranchName)
		}
	},
}

func init() {
	rootCmd.AddCommand(cloneCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// cloneCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// cloneCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	cloneCmd.Flags().BoolVarP(&createBranch, "create-branch", "b", false, `Create and check out "branch_name" from config.yaml`)
}
