package cmd

import (
	"github.com/ajschmidt8/r3/shared"
	"github.com/spf13/cobra"
)

// commitCmd represents the commit command
var commitCmd = &cobra.Command{
	Use:   "commit",
	Short: "Commit staged changes to active branch in each repository",
	Run: func(cmd *cobra.Command, args []string) {
		config := shared.ReadConfig()

		for _, repoName := range config.Repos {
			shared.Commit(repoName, config.CommitMsg)
		}
	},
}

func init() {
	rootCmd.AddCommand(commitCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// commitCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// commitCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
