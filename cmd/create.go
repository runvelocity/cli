/*
Copyright © 2023 Utibeabasi Umanah utibeabasiumanah6@gmail.com
*/
package cmd

import (
	"encoding/json"
	"io"
	"log"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/runvelocity/cli/internal/api"
	"github.com/runvelocity/cli/internal/tui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	functionName string
	filePath     string
	handler      string
)

func readBody(body io.ReadCloser) map[string]interface{} {
	var resp map[string]interface{}
	defer body.Close()
	b, err := io.ReadAll(body)
	if err != nil {
		log.Fatalln("A fatal error occured: " + err.Error())
	}
	err = json.Unmarshal(b, &resp)
	if err != nil {
		log.Fatalln("A fatal error occured: " + err.Error())
	}
	return resp
}

func initialCreateModel() tui.CreateModel {
	return tui.CreateModel{
		Function: nil,
		ApiClient: api.ApiClient{
			BaseUrl: viper.Get("managerurl").(string),
		},
		Error: nil,
		Spinner: spinner.New(
			spinner.WithSpinner(spinner.Dot),
			spinner.WithStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("205"))),
		),
		Name:     functionName,
		FilePath: filePath,
		Handler:  handler,
	}
}

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new velocity function",
	Run: func(cmd *cobra.Command, args []string) {
		// Send the UI for rendering
		createTuiModel := initialCreateModel()
		p := tea.NewProgram(createTuiModel)
		if _, err := p.Run(); err != nil {
			log.Fatalf("Alas, there's been an error: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
	createCmd.Flags().StringVar(&functionName, "name", "", "Function name")
	createCmd.Flags().StringVar(&filePath, "file-path", "", "Path to the zip file containing the function code")
	createCmd.Flags().StringVar(&handler, "handler", "", "Function handler")
	err := createCmd.MarkFlagRequired("name")
	if err != nil {
		log.Fatalln("A fatal error occured: " + err.Error())
	}
	err = createCmd.MarkFlagRequired("file-path")
	if err != nil {
		log.Fatalln("A fatal error occured: " + err.Error())
	}
	err = createCmd.MarkFlagRequired("handler")
	if err != nil {
		log.Fatalln("A fatal error occured: " + err.Error())
	}
}
