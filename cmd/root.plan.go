package cmd

import (
	"github.com/spf13/cobra"
)

var rootPlanCmd = &cobra.Command{
	Use:   "plan",
	Short: "Method to work with plans.",
}

func init() {
	rootCmd.AddCommand(rootPlanCmd)
}
