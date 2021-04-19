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
	"log"
	"os"
	"os/exec"
	"path"

	"github.com/ajschmidt8/rrr/shared"
	"github.com/spf13/cobra"
)

var (
	useInteractive bool
	addAll         bool
	doCommit       bool
	doPush         bool
	doPR           bool
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Clone and run change script in each repository",
	Long: `Clone and run change script in each repository,
then stage files interactively with git "add --patch".`,
	Run: func(cmd *cobra.Command, args []string) {
		reposDir := "repos"
		config := shared.ReadConfig()
		rootDir, _ := os.Getwd()
		scriptPath := path.Join(rootDir, "scr.sh")

		// Clone
		for _, repoName := range config.Repos {
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
			if useInteractive && addAll {
				log.Fatal(`Use "--all" or "--interactive", but not both.` + "\n")
			} else if useInteractive {
				addFlag = "-i"
			} else if addAll {
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
			if doCommit || doPush || doPR {
				shared.Commit(repoName, config.CommitMsg)
			}
			if doPush || doPR {
				shared.Push(repoName, config.BranchName, false)
			}
			if doPR {
				shared.PR(repoName, config.PR.RepoOwner, config.PR.Title, config.PR.Draft, config.PR.BaseBranch, config.BranchName, config.PR.Body, config.PR.MaintainersModify, config.PR.Labels)
			}
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
	runCmd.Flags().BoolVarP(&useInteractive, "interactive", "i", false, `Use "git add -i" instead of "git add -p".`)
	runCmd.Flags().BoolVarP(&addAll, "all", "A", false, `Use "git add -A" instead of "git add -p".`)
	runCmd.Flags().BoolVar(&doCommit, "commit", false, `Commits changes after they're made.`)
	runCmd.Flags().BoolVar(&doPush, "push", false, `Pushes changes after they're committed (implies --commit).`)
	runCmd.Flags().BoolVar(&doPR, "pr", false, `Opens a PR after changes are pushed (implies --push).`)
}
