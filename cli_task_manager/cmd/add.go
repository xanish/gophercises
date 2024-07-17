package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/xanish/gophercises/cli_task_manager/task"
	"github.com/xanish/gophercises/cli_task_manager/task_manager"
)

var (
	title       string
	description []string
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new task to your TODO list",
	Long: `Add a new task to your TODO list.

Usage:
  task add --title="Task Title" --description="Line 1" --description="Line 2"
  task add --t="Task Title" --d="Line 1"
  task add --t="Task Title"
`,
	Run: func(cmd *cobra.Command, args []string) {
		taskManager, err := task_manager.NewTaskManager(task_manager.DefaultDatabaseFilePath)
		defer taskManager.Close()

		if err != nil {
			log.Fatalf("failed to initialize task manager %v", err)
		}

		err = taskManager.Add(task.NewTask(title, description))
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("successfully added task '%s'\n", title)
	},
}

func init() {
	rootCmd.AddCommand(addCmd)

	addCmd.Flags().StringVarP(&title, "title", "t", "", "title of the task to create")
	addCmd.Flags().StringArrayVarP(&description, "description", "d", nil, "optional task description")
}
