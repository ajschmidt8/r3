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
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"

	"github.com/ajschmidt8/rrr/shared"
	"github.com/spf13/cobra"
)

var UseInteractive bool
var AddAll bool

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("run called")
		reposDir := "repos"
		config := shared.ReadConfig()
		rootDir, _ := os.Getwd()
		scriptPath := path.Join(rootDir, "scr.sh")
		fmt.Printf("scr path %s\n", scriptPath)

		// Clone
		for _, repoName := range config.Repos {
			fmt.Printf("Cloning %s\n", repoName)
			shared.Clone(repoName, config.PR.BaseBranch, config.BranchName)
		}

		// Make changes
		for _, repoName := range config.Repos {
			repoDir := path.Join(rootDir, reposDir, repoName)

			os.Chdir(repoDir)

			scrCmd := exec.Command(scriptPath)
			scrCmd.Stdout = os.Stdout
			scrCmd.Stdin = os.Stdin
			scrCmd.Stderr = os.Stderr
			err := scrCmd.Run()
			if err != nil {
				log.Fatal(err)
			}

			addFlag := "-p"
			if UseInteractive {
				addFlag = "-i"
			} else if AddAll {
				addFlag = "-A"
			}
			gitAddCmd := exec.Command("git", "add", addFlag)
			gitAddCmd.Stdout = os.Stdout
			gitAddCmd.Stdin = os.Stdin
			gitAddCmd.Stderr = os.Stderr
			gitAddCmd.Run()
		}
		os.Chdir(rootDir) // cd back to rootDir after script

		// Commit
		for _, repoName := range config.Repos {
			shared.Commit(repoName, config.CommitMsg)
			shared.Push(repoName, config.BranchName)
			shared.PR(config.PR.Title, config.PR.RepoOwner, repoName, config.PR.Draft, config.PR.BaseBranch, config.BranchName, config.PR.Body, config.PR.MaintainersModify, config.PR.Labels)
		}
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// runCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// runCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	runCmd.Flags().BoolVarP(&UseInteractive, "interactive", "i", false, `Use "git add -i" instead of "git add -p". Needed when you are adding new, untracked files to repos.`)
	runCmd.Flags().BoolVarP(&AddAll, "all", "A", false, `Use "git add -A" instead of "git add -p".`)
}
