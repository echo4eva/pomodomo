/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/echo4eva/pomogomo/internal/database"
	"github.com/spf13/cobra"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Creates a task type.",
	Long:  `Creates a task type.`,
	Args:  cobra.ExactArgs(1),
	Run:   runTaskCreate,
}

func runTaskCreate(cmd *cobra.Command, args []string) {
	db, err := database.New()
	if err != nil {
		fmt.Println(err)
	}
	db.CreateTask(database.Task{
		Name: args[0],
	})

	fmt.Println("Successfully created task: ", args[0])
}

func init() {
	taskCmd.AddCommand(createCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
