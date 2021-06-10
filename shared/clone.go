package shared

import (
	"fmt"
	"log"
	"os"
	"path"
	"strings"

	"github.com/go-git/go-git/v5"
	gitConfig "github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/google/go-github/v34/github"
	"github.com/spf13/viper"
)

func dirExists(path string) (exists bool) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		exists = false
	} else {
		exists = true
	}
	return
}

func createFork(repoName string) {
	client, ctx := GetGitHubClient()
	// the "create fork" point can be called regardless of whether a fork already exists or not
	client.Repositories.CreateFork(ctx, "rapidsai", repoName, &github.RepositoryCreateForkOptions{})
}

func Clone(repoName string, baseBranchName string, headBranchName string) {
	var gitTree *git.Worktree
	var gitRepo *git.Repository
	var err error
	reposDir := "repos"
	cwd, _ := os.Getwd()
	repoName = strings.ToLower(repoName)
	repoDir := path.Join(cwd, reposDir, repoName)
	fmt.Printf("Cloning %s\n", repoName)

	if !dirExists(repoDir) {
		createFork(repoName)
		gitRepo, err = git.PlainClone(repoDir, false, &git.CloneOptions{
			ReferenceName: plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", baseBranchName)),
			URL:           fmt.Sprintf("git@github.com:rapidsai/%v.git", repoName),
			RemoteName:    "upstream",
		})
		if err != nil {
			log.Fatalf("could not clone repo: %v", err)
		}

		gitTree, err = gitRepo.Worktree()
		if err != nil {
			log.Fatalf("could not get worktree: %v", err)
		}
		_, err = gitRepo.CreateRemote(&gitConfig.RemoteConfig{
			Name: "origin",
			URLs: []string{fmt.Sprintf("git@github.com:%s/%s.git", viper.GetString("gh_username"), repoName)},
		})
		if err != nil {
			log.Fatalf("could not create remote: %v", err)
		}
	} else {
		gitRepo, err = git.PlainOpen(repoDir)
		if err != nil {
			log.Fatalf("could not open repo: %v", err)
		}
		gitTree, err = gitRepo.Worktree()
		if err != nil {
			log.Fatalf("could not get worktree: %v", err)
		}

		// Check for local reference. Create it if it doesn't exist
		_, err = gitRepo.Reference(plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", baseBranchName)), false)
		if err != nil {
			// refs/heads/<localBranchName>
			localBranchReferenceName := plumbing.NewBranchReferenceName(baseBranchName)
			// refs/remotes/upstream/<remoteBranchName>
			remoteReferenceName := plumbing.NewRemoteReferenceName("upstream", baseBranchName)

			err = gitRepo.CreateBranch(&gitConfig.Branch{Name: baseBranchName, Remote: "upstream", Merge: localBranchReferenceName})
			if err != nil {
				log.Fatalf("could not create branch: %v", err)
			}
			newReference := plumbing.NewSymbolicReference(localBranchReferenceName, remoteReferenceName)
			err = gitRepo.Storer.SetReference(newReference)
			if err != nil {
				log.Fatalf("could not set reference: %v", err)
			}
		}

		err = gitTree.Checkout(&git.CheckoutOptions{
			Branch: plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", baseBranchName)),
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
			ReferenceName: plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", baseBranchName)),
			RemoteName:    "upstream",
		})
		if !(err == nil || err == git.NoErrAlreadyUpToDate) {
			log.Fatalf("could not pull repo: %v", err)
		}
	}

	if headBranchName != "" {
		// No errors are thrown here if the branch does not exist
		err = gitRepo.Storer.RemoveReference(plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", headBranchName)))
		if err != nil {
			log.Fatalf("could not delete branch: %v", err)
		}
		err = gitTree.Checkout(&git.CheckoutOptions{
			Branch: plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", headBranchName)),
			Force:  true,
			Create: true,
		})

		if err != nil {
			log.Fatalf("could not checkout branch: %v", err)
		}
	}

}
