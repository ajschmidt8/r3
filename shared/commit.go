package shared

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/go-git/go-git/v5"
)

type NoChangesError struct{}

func (e *NoChangesError) Error() string {
	return fmt.Sprintln("No staged files changes in worktree.")
}

// Commits the staged changes in a repo. Returns NoChangesError
// if there are no staged changes in the repo.
func Commit(repoName string, commitMsg string) error {
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
		fmt.Printf("  - No staged changes in \"%s\" repo. Skipping...\n", repoName)
		return &NoChangesError{}
	}

	gitTree.Commit(commitMsg, &git.CommitOptions{})
	fmt.Printf("  - %s\n", repoName)
	return nil
}
