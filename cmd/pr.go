package cmd

import (
	"github.com/Minnek-Digital-Studio/cominnek/pkg/github"
	"github.com/spf13/cobra"
)

var prCmd = &cobra.Command{
	Use:   "pr",
	Short: "Create a new pull request",
	Run: func(cmd *cobra.Command, args []string) {
		github.CreatePullRequest(ticket, baseBranch)
	},
}

func init() {
	prCmd.PersistentFlags().StringVarP(&ticket, "ticket", "t", "", "Ticket number")
	prCmd.PersistentFlags().StringVarP(&baseBranch, "base", "b", "develop", "Base branch")
	rootCmd.AddCommand(prCmd)
}