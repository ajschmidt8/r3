package cmd

import (
	"sync"

	"github.com/ajschmidt8/r3/shared"
	"github.com/spf13/cobra"
)

var createBranch bool

// cloneCmd represents the clone command
var cloneCmd = &cobra.Command{
	Use:   "clone",
	Short: "Clone repositories",
	Long: `Clone each repository listed in config.yaml.

By default, the "pr.base_branch" value will be checked out in each repository.
This command will create a fork of the desired repo if one does not exist already.`,
	Run: func(cmd *cobra.Command, args []string) {
		config := shared.ReadConfig()
		var wg sync.WaitGroup
		workerPoolSize := shared.ConcurrentClones
		dataCh := make(chan shared.CloneJob, workerPoolSize)

		for i := 0; i < workerPoolSize; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()

				for job := range dataCh {
					shared.Clone(job.RepoName, config.PR.BaseBranch, job.NewBranchName)
				}
			}()
		}

		for _, repoName := range config.Repos {
			newBranchName := ""
			if createBranch {
				newBranchName = config.BranchName
			}
			dataCh <- shared.CloneJob{RepoName: repoName, NewBranchName: newBranchName}
		}
		close(dataCh)
		wg.Wait()
	},
}

func init() {
	rootCmd.AddCommand(cloneCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// cloneCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// cloneCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	cloneCmd.Flags().BoolVarP(&createBranch, "create-branch", "b", false, `Create and check out "branch_name" from config.yaml`)
}
