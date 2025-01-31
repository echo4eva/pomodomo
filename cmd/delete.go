/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/echo4eva/pomogomo/internal/database"
	"github.com/spf13/cobra"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Deletes a task type.",
	Long:  `Deletes a task type.`,
	Args:  cobra.ExactArgs(1),
	Run:   runTaskDelete,
}

func runTaskDelete(cmd *cobra.Command, args []string) {
	db, err := database.New()
	if err != nil {
		fmt.Println(err)
	}
	db.DeleteTask(args[0])

	fmt.Println("Successfully deleted task: ", args[0])
}

func init() {
	taskCmd.AddCommand(deleteCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// deleteCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// deleteCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
