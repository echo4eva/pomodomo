/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"time"

	"github.com/echo4eva/pomogomo/ui"
	"github.com/spf13/cobra"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: runStart,
}

func runStart(cmd *cobra.Command, args []string) {
	minutes, _ := cmd.Flags().GetUint("minutes")
	task, _ := cmd.Flags().GetString("task")

	start := time.Now()
	end := start.Add(time.Minute * time.Duration(minutes))

	ui.Exec(start, end, task)
}

func init() {
	rootCmd.AddCommand(startCmd)
	startCmd.Flags().Uint("minutes", 30, "minutes for the timer")
	startCmd.Flags().String("task", "", "task doing")
}
