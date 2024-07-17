package cmd

import (
	"log"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/xanish/gophercises/cli_task_manager/task_manager"
)

// doCmd represents the do command
var doCmd = &cobra.Command{
	Use:   "do",
	Short: "Mark a task on your TODO list as complete",
	Long: `Mark a task on your TODO list as complete.

Usage:
  task do [id]
`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			log.Fatalf("expected id as argument")
		}

		id, err := strconv.Atoi(args[0])
		if err != nil {
			log.Fatalf("invalid task id %v", err)
		}

		taskManager, err := task_manager.NewTaskManager(task_manager.DefaultDatabaseFilePath)
		defer taskManager.Close()

		if err != nil {
			log.Fatalf("failed to initialize task manager %v", err)
		}

		err = taskManager.Complete(id)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(doCmd)
}
