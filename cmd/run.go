package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"sync"

	"github.com/ajschmidt8/r3/shared"
	"github.com/fatih/color"
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
		var wg sync.WaitGroup
		workerPoolSize := shared.ConcurrentClones
		dataCh := make(chan shared.CloneJob)

		for i := 0; i < workerPoolSize; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for job := range dataCh {
					shared.Clone(job.RepoName, config.PR.BaseBranch, config.BranchName)
				}
			}()
		}

		// Clone
		for i, repoName := range config.Repos {
			if i == 0 {
				color.New(color.Bold).Println("Cloning repos:")
			}
			dataCh <- shared.CloneJob{RepoName: repoName}
		}
		close(dataCh)
		wg.Wait()

		fmt.Println()
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

			if addFlag == "-A" {
				fmt.Print("Staged all changes for ")
			} else {
				fmt.Print("\nChanges for ")
			}
			color.New(color.FgGreen, color.Bold).Printf("%s\n", repoName)
			gitAddCmd := exec.Command("git", "add", addFlag)
			gitAddCmd.Stdout = os.Stdout
			gitAddCmd.Stdin = os.Stdin
			gitAddCmd.Stderr = os.Stderr
			gitAddCmd.Run()
		}
		os.Chdir(rootDir) // cd back to rootDir after script

		changedRepos := make([]string, 0, len(config.Repos))

		// Commit
		for i, repoName := range config.Repos {
			if doCommit || doPush || doPR {
				if i == 0 {
					fmt.Println()
					color.New(color.Bold).Println("Committing changes:")
				}
				err := shared.Commit(repoName, config.CommitMsg)
				if _, ok := err.(*shared.NoChangesError); ok {
					continue
				}
				changedRepos = append(changedRepos, repoName)
			}
		}

		// Push
		for i, repoName := range changedRepos {
			if doPush || doPR {
				if i == 0 {
					fmt.Println()
					color.New(color.Bold).Println("Pushing changes:")
				}
				shared.Push(repoName, config.BranchName, false)
			}
		}

		// PR
		for i, repoName := range changedRepos {
			if doPR {
				if i == 0 {
					fmt.Println()
					color.New(color.Bold).Println("Opening PRs:")
				}
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
