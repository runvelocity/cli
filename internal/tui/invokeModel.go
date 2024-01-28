package tui

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/alecthomas/chroma/v2/quick"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/runvelocity/cli/internal/api"
	"github.com/runvelocity/cli/internal/models"
)

type FunctionInvokeResponseMsg struct {
	InvokeResponse *models.InvokeResponse
}

type InvokeModel struct {
	Spinner        spinner.Model
	ApiClient      api.ApiClient
	Error          error
	InvokeResponse *models.InvokeResponse
	FunctionName   string
	Payload        map[string]interface{}
}

func (m InvokeModel) InvokeFunction(functionName string, payload map[string]interface{}) tea.Msg {
	response, err := m.ApiClient.InvokeFunction(functionName, payload)
	if err != nil {
		return ErrorMsg{
			Error: err,
		}
	}
	return FunctionInvokeResponseMsg{
		InvokeResponse: response,
	}
}

func (m InvokeModel) IsInvoked() bool {
	return m.InvokeResponse != nil
}

func (m InvokeModel) Init() tea.Cmd {
	return tea.Batch(m.Spinner.Tick,
		func() tea.Msg {
			return m.InvokeFunction(m.FunctionName, m.Payload)
		})
}

func (m InvokeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {

	case spinner.TickMsg:
		m.Spinner, cmd = m.Spinner.Update(msg)
		return m, cmd
	case tea.KeyMsg:
		switch msg.String() {

		case "ctrl+c", "q":
			return m, tea.Quit

		}

	case FunctionInvokeResponseMsg:
		m.InvokeResponse = msg.InvokeResponse
	case ErrorMsg:
		m.Error = msg.Error
	}

	return m, cmd
}

func (m InvokeModel) View() string {
	if m.Error != nil {
		return printError(m.Error)
	}
	if !m.IsInvoked() {
		return fmt.Sprintf("%s %s", m.Spinner.View(), "Invoking function..."+boldString("Press CTRL+C or q to quit"))
	}

	jsonString, err := json.MarshalIndent(m.InvokeResponse.InvocationResponse, "", "   ")
	if err != nil {
		return printError(err)
	}

	var resultBuilder strings.Builder

	err = quick.Highlight(&resultBuilder, string(jsonString), "json", "terminal256", "pygments")
	if err != nil {
		fmt.Println("Error:", err)
		return printError(err)
	}

	return printSuccess("Status Code: " + fmt.Sprintf("%d", m.InvokeResponse.StatusCode) + "\n" + resultBuilder.String())
}
