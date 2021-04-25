package shared

import (
	"context"

	"github.com/google/go-github/v34/github"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
)

func GetGitHubClient() (client *github.Client, ctx context.Context) {
	ctx = context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: viper.GetString("gh_token")},
	)
	tc := oauth2.NewClient(ctx, ts)

	client = github.NewClient(tc)
	return
}
