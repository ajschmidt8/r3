package shared

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/go-git/go-git/v5"
)

// If there are staged files, commits them to the
// currently active branch. Does nothing if there
// are no staged files
func Commit(repoName string, commitMsg string) {
	reposDir := "repos"
	cwd, _ := os.Getwd()
	repoDir := path.Join(cwd, reposDir, repoName)

	gitRepo, err := git.PlainOpen(repoDir)
	if err != nil {
		log.Fatalf("could not open repo: %v", err)
	}

	gitTree, err := gitRepo.Worktree()
	if err != nil {
		log.Fatalf("could not get worktree: %v", err)
	}
	status, err := gitTree.Status()
	if err != nil {
		log.Fatalf("could not get worktree status: %v", err)
	}

	hasStagedChanges := false
	for _, val := range status {
		if val.Staging != ' ' && val.Staging != '?' {
			hasStagedChanges = true
			break
		}
	}

	if !hasStagedChanges {
		fmt.Printf("No staged changes in %s repo. Skipping...\n", repoName)
		return
	}

	gitTree.Commit(commitMsg, &git.CommitOptions{})
}
