package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var version = "v0.1"

var rootCmd = &cobra.Command{
	Use:     "tftools",
	Short:   "Provides tools for working with terraform.",
	Version: version,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
}
