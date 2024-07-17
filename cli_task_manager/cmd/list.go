package cmd

import (
	"fmt"
	"github.com/xanish/gophercises/cli_task_manager/task_manager"
	"log"

	"github.com/spf13/cobra"
	"github.com/xanish/gophercises/cli_task_manager/task"
)

var pending bool
var completed bool

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Displays a list of tasks present in your TODOs.",
	Long: `Displays a list of tasks present in your TODOs. Can be filtered to show only Pending or Completed tasks.

Usage:
  task list --pending
  task list --completed
  task list -p
  task list -c
`,
	Run: func(cmd *cobra.Command, args []string) {
		var status string
		if pending {
			status = task.StatusPending
		} else if completed {
			status = task.StatusCompleted
		}

		var tasks []task.Task

		taskManager, err := task_manager.NewTaskManager(task_manager.DefaultDatabaseFilePath)
		defer taskManager.Close()

		if err != nil {
			log.Fatalf("failed to initialize task manager %v", err)
		}

		tasks, err = taskManager.List(status)
		if err != nil {
			log.Fatal(err)
		}

		for _, task := range tasks {
			fmt.Println(task)
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	listCmd.Flags().BoolVarP(&pending, "pending", "p", false, "display pending tasks")
	listCmd.Flags().BoolVarP(&completed, "completed", "c", false, "display completed tasks")

	listCmd.MarkFlagsMutuallyExclusive(task.StatusPending, task.StatusCompleted)
}
