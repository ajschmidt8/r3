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
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"

	"github.com/ajschmidt8/rrr/shared"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/spf13/cobra"
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

		for _, repo := range config.Repos {
			fmt.Printf("Cloning %s\n", repo)
			repoDir := path.Join(cwd, reposDir, repo)

			if !dirExists(repoDir) {
				fmt.Printf("repos dir doesn't exist, cloning %v", config.PR.BaseBranch)
				_, err := git.PlainClone(repoDir, false, &git.CloneOptions{
					ReferenceName: plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", config.PR.BaseBranch)),
					Progress:      os.Stdout,
					URL:           fmt.Sprintf("git@github.com:rapidsai/%v.git", repo),
				})
				if err != nil {
					log.Fatalf("could not clone repo: %v", err)
				}
			} else {
				fmt.Printf("repos dir does exist!")
				gitRepo, err := git.PlainOpen(repoDir)
				if err != nil {
					log.Fatalf("could not open repo: %v", err)
				}
				gitTree, err := gitRepo.Worktree()
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
					SingleBranch:  true,
				})
				if !(err == nil || err == git.NoErrAlreadyUpToDate) {
					log.Fatalf("could not pull repo: %v", err)
				}
			}
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
			}
			gitAddCmd := exec.Command("git", "add", addFlag)
			gitAddCmd.Stdout = os.Stdout
			gitAddCmd.Stdin = os.Stdin
			gitAddCmd.Stderr = os.Stderr
			gitAddCmd.Run()
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
