package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"tftools/tf"
)

var rootPlanSummaryOps struct {
	PlanFilePath      string
	Basic             bool
	ResourceTypes     bool
	ResourceAddresses bool
}

var rootPlanSummaryCmd = &cobra.Command{
	Use:   "summary",
	Short: "Generates a summary of the plan.",
	Run: func(cmd *cobra.Command, args []string) {

		// Read the plan file from disk.
		planFile, err := os.ReadFile(rootPlanSummaryOps.PlanFilePath)
		if err != nil {
			fmt.Printf("failed to read plan file: %s", err)
			os.Exit(1)
		}

		// Parse the plan file.
		tfPlan := &tf.Plan{}
		err = json.Unmarshal(planFile, tfPlan)
		if err != nil {
			fmt.Printf("failed to parse plan file: %s", err)
			os.Exit(1)
		}

		importByAddress := make([]string, 0, 0)
		destroyByAddress := make([]string, 0, 0)
		noopByAddress := make([]string, 0, 0)
		addByAddress := make([]string, 0, 0)
		changeByAddress := make([]string, 0, 0)
		unknownByAddress := make([]string, 0, 0)

		importByResourceType := make(map[string]int)
		destroyByResourceType := make(map[string]int)
		noopByResourceType := make(map[string]int)
		addByResourceType := make(map[string]int)
		changeByResourceType := make(map[string]int)
		unknownByResourceType := make(map[string]int)

		// Processing plan.
		for _, change := range tfPlan.ResourceChanges {
			for _, action := range change.Change.Actions {
				switch action {
				case "delete":
					destroyByAddress = append(destroyByAddress, change.Address)
					destroyByResourceType[change.Type]++
				case "no-op":
					if change.Change.Importing.ID != "" {
						importByAddress = append(importByAddress, change.Address)
						importByResourceType[change.Type]++
						continue
					}
					noopByAddress = append(noopByAddress, change.Address)
					noopByResourceType[change.Type]++
				case "create":
					addByAddress = append(addByAddress, change.Address)
					addByResourceType[change.Type]++
				case "update":
					changeByAddress = append(changeByAddress, change.Address)
					changeByResourceType[change.Type]++
				default:
					unknownByAddress = append(unknownByAddress, action)
					unknownByResourceType[change.Type]++
				}
			}
		}

		// Output based on options selected.
		switchFound := false

		if rootPlanSummaryOps.ResourceAddresses {
			switchFound = true
			printResourceAddress(
				importByAddress,
				destroyByAddress,
				noopByAddress,
				addByAddress,
				changeByAddress)
		}

		if rootPlanSummaryOps.ResourceTypes {
			printResourceTypes(
				importByResourceType,
				destroyByResourceType,
				noopByResourceType,
				addByResourceType,
				changeByResourceType,
				unknownByResourceType)
		}

		if rootPlanSummaryOps.Basic {
			switchFound = true
			printBasic(
				importByAddress,
				destroyByAddress,
				noopByAddress,
				addByAddress,
				changeByAddress)
		}

		if !switchFound {
			printBasic(
				importByAddress,
				destroyByAddress,
				noopByAddress,
				addByAddress,
				changeByAddress)
		}

	},
}

func init() {
	rootPlanSummaryCmd.Flags().StringVar(&rootPlanSummaryOps.PlanFilePath, "plan-file", "", "The path to the plan file in json.")
	rootPlanSummaryCmd.Flags().BoolVar(&rootPlanSummaryOps.Basic, "basic", false, "Display the basic section.")
	rootPlanSummaryCmd.Flags().BoolVar(&rootPlanSummaryOps.ResourceTypes, "resource-types", false, "Display summary based on actions for each resource type.")
	rootPlanSummaryCmd.Flags().BoolVar(&rootPlanSummaryOps.ResourceAddresses, "resource-addresses", false, "Display summary based on actions for each resource address.")
	rootPlanSummaryCmd.MarkFlagRequired("plan-file")
	rootPlanCmd.AddCommand(rootPlanSummaryCmd)

}

// Prints the basic information to stdout.
func printBasic(
	importByAddress []string,
	destroyByAddress []string,
	noopByAddress []string,
	addByAddress []string,
	changeByAddress []string,
) {

	fmt.Printf(
		"Plan: %d to import, %d to add, %d to change, %d to destroy\n",
		len(importByAddress), len(addByAddress), len(changeByAddress), len(destroyByAddress))
}

// Prints the modified resources by address to stdout.
func printResourceAddress(
	importByAddress []string,
	destroyByAddress []string,
	noopByAddress []string,
	addByAddress []string,
	changeByAddress []string,
) {

	for _, address := range importByAddress {
		fmt.Printf("> %s\n", address)
	}
	for _, address := range destroyByAddress {
		fmt.Printf("- %s\n", address)
	}
	for _, address := range addByAddress {
		fmt.Printf("+ %s\n", address)
	}
	for _, address := range changeByAddress {
		fmt.Printf("~ %s\n", address)
	}

}

// Prints a summary of changes based on
func printResourceTypes(
	importByResourceType map[string]int,
	destroyByResourceType map[string]int,
	noopByResourceType map[string]int,
	addByResourceType map[string]int,
	changeByResourceType map[string]int,
	unknownByResourceType map[string]int) {

	for resource, count := range importByResourceType {
		fmt.Printf("Importing %4d %s\n", count, resource)
	}

	for resource, count := range destroyByResourceType {
		fmt.Printf("Destroying %4d %s\n", count, resource)
	}

	for resource, count := range addByResourceType {
		fmt.Printf("Adding %4d %s\n", count, resource)
	}

	for resource, count := range changeByResourceType {
		fmt.Printf("Changing  %4d %s\n", count, resource)
	}

	for resource, count := range unknownByResourceType {
		fmt.Printf("Unknown Change %4d %s\n", count, resource)
	}

}
