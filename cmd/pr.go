package cmd

import (
	"github.com/ajschmidt8/r3/shared"
	"github.com/spf13/cobra"
)

// prCmd represents the pr command
var prCmd = &cobra.Command{
	Use:   "pr",
	Short: `Open a PR for each repository according to the information in config.yaml`,
	Run: func(cmd *cobra.Command, args []string) {
		config := shared.ReadConfig()

		for _, repoName := range config.Repos {
			shared.PR(repoName, config.PR.RepoOwner, config.PR.Title, config.PR.Draft, config.PR.BaseBranch, config.BranchName, config.PR.Body, config.PR.MaintainersModify, config.PR.Labels)
		}
	},
}

func init() {
	rootCmd.AddCommand(prCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// prCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// prCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
