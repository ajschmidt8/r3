package shared

import (
	"context"
	"fmt"
	"log"

	"github.com/google/go-github/v34/github"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
)

func PR(repoName string, repoOwner string, title string, draft bool, baseBranch string, headBranch string, body string, maintainerModify bool, labels []string) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: viper.GetString("gh_token")},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)
	newPR := &github.NewPullRequest{
		Title:               github.String(title),
		Base:                github.String(baseBranch),
		Head:                github.String(headBranch),
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
