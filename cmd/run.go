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
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"

	"github.com/ajschmidt8/rrr/shared"
	"github.com/go-git/go-git/v5"
	gitConfig "github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/google/go-github/v34/github"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
	"gopkg.in/yaml.v2"
)

var UseInteractive bool

func dirExists(path string) (exists bool) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		exists = false
	} else {
		exists = true
	}
	return
}

// runCmd represents the clone command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		var config shared.ConfigInterface
		fmt.Println("clone called")
		reposDir := "repos"
		cwd, _ := os.Getwd()
		ymlBytes, err := ioutil.ReadFile("config.yaml")
		if err != nil {
			log.Fatalf("cannot read config.yaml file: %v", err)
		}
		err = yaml.Unmarshal(ymlBytes, &config)
		if err != nil {
			log.Fatalf("cannot decode config.yaml: %v", err)
		}
		scriptPath := path.Join(cwd, "scr.sh")
		fmt.Printf("scr path %s\n", scriptPath)

		for _, repoName := range config.Repos {
			var gitTree *git.Worktree
			var gitRepo *git.Repository
			fmt.Printf("Cloning %s\n", repoName)
			repoDir := path.Join(cwd, reposDir, repoName)

			if !dirExists(repoDir) {
				fmt.Printf("repos dir doesn't exist, cloning %v", config.PR.BaseBranch)
				gitRepo, err = git.PlainClone(repoDir, false, &git.CloneOptions{
					ReferenceName: plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", config.PR.BaseBranch)),
					Progress:      os.Stdout,
					URL:           fmt.Sprintf("git@github.com:rapidsai/%v.git", repoName),
					RemoteName:    "upstream",
				})
				if err != nil {
					log.Fatalf("could not clone repo: %v", err)
				}
				_, err = gitRepo.CreateRemote(&gitConfig.RemoteConfig{
					Name: "origin",
					URLs: []string{fmt.Sprintf("git@github.com:ajschmidt8/%v.git", repoName)},
				})
				if err != nil {
					log.Fatalf("could not create remote: %v", err)
				}
				gitTree, err = gitRepo.Worktree()
				if err != nil {
					log.Fatalf("could not get worktree: %v", err)
				}
			} else {
				fmt.Printf("repos dir does exist!")
				gitRepo, err = git.PlainOpen(repoDir)
				if err != nil {
					log.Fatalf("could not open repo: %v", err)
				}
				gitTree, err = gitRepo.Worktree()
				if err != nil {
					log.Fatalf("could not get worktree: %v", err)
				}
				err = gitTree.Checkout(&git.CheckoutOptions{
					Branch: plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", config.PR.BaseBranch)),
					Force:  true,
				})
				if err != nil {
					log.Fatalf("could not checkout in repo: %v", err)
				}
				err = gitTree.Clean(&git.CleanOptions{
					Dir: true,
				})
				if err != nil {
					log.Fatalf("could not clean repo: %v", err)
				}
				err = gitTree.Pull(&git.PullOptions{
					ReferenceName: plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", config.PR.BaseBranch)),
				})
				if !(err == nil || err == git.NoErrAlreadyUpToDate) {
					log.Fatalf("could not pull repo: %v", err)
				}
			}
			os.Chdir(repoDir)

			// No errors are thrown here if the branch does not exist
			err = gitRepo.Storer.RemoveReference(plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", config.BranchName)))
			if err != nil {
				log.Fatalf("could not delete branch: %v", err)
			}
			err = gitTree.Checkout(&git.CheckoutOptions{
				Branch: plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", config.BranchName)),
				Create: true,
			})

			if err != nil {
				log.Fatalf("could not checkout branch: %v", err)
			}

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
			}
			gitAddCmd := exec.Command("git", "add", addFlag)
			gitAddCmd.Stdout = os.Stdout
			gitAddCmd.Stdin = os.Stdin
			gitAddCmd.Stderr = os.Stderr
			gitAddCmd.Run()
			gitTree.Commit(config.CommitMsg, &git.CommitOptions{})
			err = gitRepo.Push(&git.PushOptions{
				RemoteName: "origin",
				RefSpecs:   []gitConfig.RefSpec{gitConfig.RefSpec(fmt.Sprintf("refs/heads/%[1]s:refs/heads/%[1]s", config.BranchName))},
			})
			if err != nil {
				log.Fatalf("could not push branch: %v", err)
			}

			// Open PR
			ctx := context.Background()
			ts := oauth2.StaticTokenSource(
				&oauth2.Token{AccessToken: ""},
			)
			tc := oauth2.NewClient(ctx, ts)

			client := github.NewClient(tc)
			newPR := &github.NewPullRequest{
				Title:               github.String(config.PR.Title),
				Head:                github.String(config.BranchName),
				Base:                github.String(config.PR.BaseBranch),
				Body:                github.String(config.PR.Body),
				MaintainerCanModify: github.Bool(config.PR.MaintainersModify),
			}

			pr, _, err := client.PullRequests.Create(ctx, config.PR.RepoOwner, repoName, newPR)
			if err != nil {
				log.Fatalf("could not create PR: %v", err)
			}

			_, _, err = client.Issues.AddLabelsToIssue(ctx, config.PR.RepoOwner, repoName, pr.GetNumber(), config.PR.Labels)
			if err != nil {
				log.Fatalf("could not add labels: %v", err)
			}

			fmt.Printf("\nPR created: %s\n", pr.GetHTMLURL())
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
}
