package cmd

import (
	"fmt"
	"io/ioutil"

	"github.com/ajschmidt8/rrr/shared"
	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Generate necessary config files",
	Long: `Generates the following config files:

  - scr.sh - script to be run in each repo
  - config.yaml - yaml file to define repos to change & PR information`,
	Run: func(cmd *cobra.Command, args []string) {
		base_branch := getLatestBranch()
		ioutil.WriteFile("scr.sh", []byte(shared.Script), 0755)
		ioutil.WriteFile("config.yaml", []byte(shared.Config(base_branch)), 0644)
		fmt.Println("Run `rrr -h` or visit https://github.com/ajschmidt8/rapids-repo-reviser for usage instructions.")
	},
}

func getLatestBranch() (branch string) {
	branch = "branch-0.xx  # branch to base your changes off of"
	client, ctx := shared.GetGitHubClient()
	repo, _, err := client.Repositories.Get(ctx, "rapidsai", "cudf")
	if err != nil {
		return
	}
	branch = *(repo.DefaultBranch)
	return
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
