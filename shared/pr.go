package shared

import (
	"fmt"
	"log"

	"github.com/google/go-github/v34/github"
	"github.com/spf13/viper"
)

func PR(repoName string, repoOwner string, title string, draft bool, baseBranch string, headBranch string, body string, maintainerModify bool, labels []string) {

	client, ctx := GetGitHubClient()

	prHead := headBranch
	prAuthor := viper.GetString("gh_username")

	// Prefix head branch with "<author>:" if opening PR from a fork
	if repoOwner != prAuthor {
		prHead = prAuthor + ":" + headBranch
	}
	newPR := &github.NewPullRequest{
		Title:               github.String(title),
		Base:                github.String(baseBranch),
		Head:                github.String(prHead),
		Body:                github.String(body),
		MaintainerCanModify: github.Bool(maintainerModify),
		Draft:               github.Bool(draft),
	}

	pr, _, err := client.PullRequests.Create(ctx, repoOwner, repoName, newPR)
	if err != nil {
		log.Fatalf("could not create PR: %v", err)
	}

	_, _, err = client.Issues.AddLabelsToIssue(ctx, repoOwner, repoName, pr.GetNumber(), labels)
	if err != nil {
		log.Fatalf("could not add labels: %v", err)
	}

	fmt.Printf("\nPR created: %s\n", pr.GetHTMLURL())
}
