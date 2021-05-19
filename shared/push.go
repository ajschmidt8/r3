package shared

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/go-git/go-git/v5"
	gitConfig "github.com/go-git/go-git/v5/config"
)

func Push(repoName string, headBranchName string, deleteBranch bool) {
	reposDir := "repos"
	cwd, _ := os.Getwd()
	repoDir := path.Join(cwd, reposDir, repoName)

	gitRepo, err := git.PlainOpen(repoDir)
	if err != nil {
		log.Fatalf("could not open repo: %v", err)
	}

	if deleteBranch {
		gitRemote, err := gitRepo.Remote("origin")
		if err != nil {
			log.Fatalf("could not get remote: %v", err)
		}
		err = gitRemote.Push(&git.PushOptions{
			RefSpecs: []gitConfig.RefSpec{gitConfig.RefSpec(":refs/heads/" + headBranchName)},
		})
		if err == git.NoErrAlreadyUpToDate {
			fmt.Printf("No \"%s\" branch to delete in repo \"%s\". Skipping...\n", headBranchName, repoName)
			return
		}
		if err != nil {
			log.Fatalf("could not delete remote branch: %v", err)
		}
		fmt.Printf("Deleted branch: %s\n", headBranchName)
		return
	}

	err = gitRepo.Push(&git.PushOptions{
		RemoteName: "origin",
		RefSpecs:   []gitConfig.RefSpec{gitConfig.RefSpec(fmt.Sprintf("refs/heads/%[1]s:refs/heads/%[1]s", headBranchName))},
	})
	if err != nil {
		log.Fatalf("could not push branch: %v", err)
	}
}
