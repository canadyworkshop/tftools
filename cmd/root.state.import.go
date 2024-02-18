package cmd

import "github.com/spf13/cobra"

var rootStateImportCmd = &cobra.Command{
	Use:   "import",
	Short: "Method to work with state imports.",
}

func init() {
	rootStateCmd.AddCommand(rootStateImportCmd)
}
