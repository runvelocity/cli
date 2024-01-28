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

var functionNameToInvoke string

func initialInvokeModel(payload map[string]interface{}) tui.InvokeModel {
	return tui.InvokeModel{
		FunctionName: functionNameToInvoke,
		ApiClient: api.ApiClient{
			BaseUrl: viper.Get("managerurl").(string),
		},
		Error:   nil,
		Payload: payload,
		Spinner: spinner.New(
			spinner.WithSpinner(spinner.Dot),
			spinner.WithStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("205"))),
		),
	}
}

// invokeCmd represents the invoke command
var invokeCmd = &cobra.Command{
	Use:   "invoke",
	Short: "Invoke a function",
	Run: func(cmd *cobra.Command, args []string) {
		payload := map[string]interface{}{
			"key1": "value1",
			"key2": "value2",
		}
		// Send the UI for rendering
		invokeTuiModel := initialInvokeModel(payload)
		p := tea.NewProgram(invokeTuiModel)
		if _, err := p.Run(); err != nil {
			log.Fatalf("Alas, there's been an error: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(invokeCmd)
	invokeCmd.Flags().StringVar(&functionNameToInvoke, "name", "", "Function to invoke")
	err := invokeCmd.MarkFlagRequired("name")
	if err != nil {
		log.Fatalln("A fatal error occured: " + err.Error())
	}
}
