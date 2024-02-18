package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"strings"
	"tftools/tf"
)

var cliOpts = struct {
	StateFilePath     string
	ResourcePrefix    string
	NewResourcePrefix string
}{}

var rootStateImportGenerate = &cobra.Command{
	Use:   "generate",
	Short: "Generates Terraform import statements for moving resources between states.",
	Run: func(cmd *cobra.Command, args []string) {

		// Read the state file from disk.
		stateFile, err := os.ReadFile(cliOpts.StateFilePath)
		if err != nil {
			fmt.Printf("failed to read state file: %s", err)
			os.Exit(1)
		}

		// Parse the state file.
		state := &tf.State{}
		err = json.Unmarshal(stateFile, state)
		if err != nil {
			fmt.Printf("failed to parse state file: %s", err)
			os.Exit(1)
		}

		// Retrieve all the resources based on the prefix provided.
		matchedResources := make(map[string]string)
		for _, r := range state.Values.RootModule.Resources {
			// Process all local resources of the module.
			if isPrefix(r.Address, cliOpts.ResourcePrefix) {
				matchedResources[r.Address] = r.Values.ID
			}

			// Process all sub modules.
			for _, c := range state.Values.RootModule.ChildModules {
				processChildModules(c, cliOpts.ResourcePrefix, matchedResources)
			}
		}

		// Output all the import statements.
		if cliOpts.NewResourcePrefix == "" {
			cliOpts.NewResourcePrefix = cliOpts.ResourcePrefix
		}

		for address, id := range matchedResources {
			newAddress := strings.Replace(address, cliOpts.ResourcePrefix, cliOpts.NewResourcePrefix, 1)
			// Escape quotes
			//newAddress = strings.Replace(newAddress, "\"", "\\\"", 100)
			fmt.Printf("#%s\nimport {\n    to = %s\n    id = \"%s\"\n}\n", address, newAddress, id)
		}

	},
}

func init() {
	rootStateImportGenerate.Flags().StringVar(&cliOpts.StateFilePath, "state-file", "", "The path to the state file in json format.")
	rootStateImportGenerate.Flags().StringVar(&cliOpts.ResourcePrefix, "resource-prefix", "", "The prefix to limit resources that will have imports generated.")
	rootStateImportGenerate.Flags().StringVar(&cliOpts.NewResourcePrefix, "new-resource-prefix", "", "An optional prefix that will replace the prefix provided for scoping.")
	rootStateImportGenerate.MarkFlagRequired("state-file")
	rootStateImportCmd.AddCommand(rootStateImportGenerate)
}

func isPrefix(address string, prefix string) bool {
	if len(address) <= len(prefix) {
		return false
	}

	if address[0:len(prefix)] != prefix {
		return false
	}

	return true
}

func processChildModules(cm tf.ChildModule, prefix string, resourceMap map[string]string) {
	for _, r := range cm.Resources {
		if isPrefix(r.Address, prefix) {
			resourceMap[r.Address] = r.Values.ID
		}
	}

	for _, r := range cm.ChildModules {
		processChildModules(r, prefix, resourceMap)
	}
}
