package cmd

import (
	"github.com/ajschmidt8/r3/shared"
	"github.com/spf13/cobra"
)

var deleteBranch bool

// pushCmd represents the push command
var pushCmd = &cobra.Command{
	Use:   "push",
	Short: `Push "branch_name" branch from config.yaml to forked repositories`,
	Run: func(cmd *cobra.Command, args []string) {
		config := shared.ReadConfig()

		for _, repoName := range config.Repos {
			shared.Push(repoName, config.BranchName, deleteBranch)
		}
	},
}

func init() {
	rootCmd.AddCommand(pushCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// pushCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// pushCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	pushCmd.Flags().BoolVarP(&deleteBranch, "delete", "d", false, `Instead of pushing branch, deletes branch from remote if it exists already.`)
}
