package cmd

import (
	"os"
	"path"

	"github.com/ajschmidt8/rrr/shared"
	"github.com/cli/oauth"
	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var VERSION = "0.0.0"

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "rrr",
	Long:    "A CLI tool for programmatically making changes across several RAPIDS repos.",
	Version: VERSION,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	// Find home directory.
	home, err := homedir.Dir()
	cobra.CheckErr(err)

	viper.AddConfigPath(home)
	viper.SetConfigName(".rrr")

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		file, err := os.Create(path.Join(home, ".rrr.yaml"))
		cobra.CheckErr(err)
		err = file.Chmod(0644)
		cobra.CheckErr(err)

		flow := &oauth.Flow{
			Hostname: "github.com",
			ClientID: os.Getenv("OAUTH_CLIENT_ID"),
			Scopes:   []string{"repo", "user"},
		}

		githubToken, err := flow.DetectFlow()
		if err != nil {
			panic(err)
		}
		viper.Set("gh_token", githubToken.Token)
		viper.WriteConfig()

		client, ctx := shared.GetGitHubClient()
		user, _, err := client.Users.Get(ctx, "")
		if err != nil {
			panic(err)
		}

		viper.Set("gh_username", user.GetLogin())
		viper.WriteConfig()
	}

}
