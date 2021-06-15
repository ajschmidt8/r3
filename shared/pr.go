package shared

import (
	"fmt"
	"time"

	"github.com/fatih/color"
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
		color.New(color.FgRed, color.Bold).Printf("Error opening PR for \"%s\": %v\n", repoName, err)
		return
	}
	fmt.Printf("%s\n", pr.GetHTMLURL())

	time.Sleep(1 * time.Second)

	_, _, err = client.Issues.AddLabelsToIssue(ctx, repoOwner, repoName, pr.GetNumber(), labels)
	if err != nil {
		color.New(color.FgRed, color.Bold).Printf("Error addings labels to \"%s\" PR: %v\n", repoName, err)
	}
}
