/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/echo4eva/pomogomo/internal/database"
	"github.com/spf13/cobra"
)

// viewCmd represents the view command
var viewCmd = &cobra.Command{
	Use:   "view",
	Short: "View tasks types.",
	Long:  `View task types`,
	Args:  cobra.NoArgs,
	Run:   runTaskView,
}

func runTaskView(cmd *cobra.Command, args []string) {
	db, err := database.New()
	if err != nil {
		fmt.Println(err)
	}
	tasks, err := db.RetrieveTasks()
	if err != nil {
		fmt.Println(err)
	}

	if len(tasks) == 0 {
		fmt.Println("No task types")
	}

	fmt.Println("ID\tName")
	for id, task := range tasks {
		fmt.Printf("%d\t%s\n", id, task.Name)
	}
}

func init() {
	taskCmd.AddCommand(viewCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// viewCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// viewCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
