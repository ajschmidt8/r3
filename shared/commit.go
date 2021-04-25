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

	// for key, val := range status {
	// fmt.Println("key:", key)
	// fmt.Println("val:", string(val.Worktree))
	// }

	// fmt.Printf("status: %v\n", status)
	// fmt.Printf("status length: %d\n", len(status))

	// FIXME: check ONLY for staged files.
	// This block won't run if the only change in a repo
	// is a new, untracked file. An empty commit will result
	if len(status) == 0 {
		return &NoChangesError{}
	}

	gitTree.Commit(commitMsg, &git.CommitOptions{})
	return nil
}
