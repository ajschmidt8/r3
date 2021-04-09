package shared

import (
	"log"
	"os"
	"path"

	"github.com/go-git/go-git/v5"
)

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

	gitTree.Commit(commitMsg, &git.CommitOptions{})
}
