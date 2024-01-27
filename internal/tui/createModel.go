package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/runvelocity/cli/internal/api"
	"github.com/runvelocity/cli/internal/models"
)

type FunctionCreateResponseMsg struct {
	Function models.Function
}

type CreateModel struct {
	Spinner   spinner.Model
	ApiClient api.ApiClient
	Error     error
	Function  *models.Function
	Name      string
	FilePath  string
	Handler   string
}

func (m CreateModel) CreateFunction(functionName, filePath, handler string) tea.Msg {
	function, err := m.ApiClient.CreateFunction(functionName, filePath, handler)
	if err != nil {
		return ErrorMsg{
			Error: err,
		}
	}
	return FunctionCreateResponseMsg{
		Function: *function,
	}
}

func (m CreateModel) IsCreated() bool {
	return m.Function != nil
}

func (m CreateModel) Init() tea.Cmd {
	return tea.Batch(m.Spinner.Tick,
		func() tea.Msg {
			return m.CreateFunction(m.Name, m.FilePath, m.Handler)
		})
}

func (m CreateModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

	case FunctionCreateResponseMsg:
		m.Function = &msg.Function
	case ErrorMsg:
		m.Error = msg.Error
	}

	return m, cmd
}

func (m CreateModel) View() string {
	if m.Error != nil {
		return printError(m.Error)
	}
	if !m.IsCreated() {
		return fmt.Sprintf("%s %s", m.Spinner.View(), "Creating function..."+boldString("Press CTRL+C or q to quit"))
	}

	return printSuccess("Created function " + m.Function.Name)
}
