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

func initialListModel() tui.ListModel {
	return tui.ListModel{
		Functions: nil,
		ApiClient: api.ApiClient{
			BaseUrl: viper.Get("managerurl").(string),
		},
		Error: nil,
		Spinner: spinner.New(
			spinner.WithSpinner(spinner.Dot),
			spinner.WithStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("205"))),
		),
	}
}

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Fetch all functions",
	Run: func(cmd *cobra.Command, args []string) {

		// Send the UI for rendering
		listTuiModel := initialListModel()
		p := tea.NewProgram(listTuiModel)
		if _, err := p.Run(); err != nil {
			log.Fatalf("Alas, there's been an error: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
