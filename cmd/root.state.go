package cmd

import "github.com/spf13/cobra"

var rootStateCmd = &cobra.Command{
	Use:   "state",
	Short: "Method to work with states.",
}

func init() {
	rootCmd.AddCommand(rootStateCmd)
}
