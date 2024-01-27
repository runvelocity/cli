/*
Copyright © 2023 Utibeabasi Umanah utibeabasiumanah6@gmail.com
*/
package cmd

import (
	"log"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/runvelocity/cli/internal/api"
	"github.com/runvelocity/cli/internal/tui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var functionNameToDelete string

func initialDeleteModel() tui.DeleteModel {
	return tui.DeleteModel{
		Function: nil,
		ApiClient: api.ApiClient{
			BaseUrl: viper.Get("managerurl").(string),
		},
		Error: nil,
		Spinner: spinner.New(
			spinner.WithSpinner(spinner.Dot),
			spinner.WithStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("205"))),
		),
		Name: functionNameToDelete,
	}
}

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a function",
	Run: func(cmd *cobra.Command, args []string) {
		// Send the UI for rendering
		deleteTuiModel := initialDeleteModel()
		p := tea.NewProgram(deleteTuiModel)
		if _, err := p.Run(); err != nil {
			log.Fatalf("Alas, there's been an error: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
	deleteCmd.Flags().StringVar(&functionNameToDelete, "name", "", "Function name to delete")
	err := deleteCmd.MarkFlagRequired("name")
	if err != nil {
		log.Fatalln("A fatal error occured: " + err.Error())
	}
}
