package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/runvelocity/cli/internal/api"
	"github.com/runvelocity/cli/internal/models"
)

type FunctionDeleteResponseMsg struct {
	Function models.Function
}

type DeleteModel struct {
	Spinner   spinner.Model
	ApiClient api.ApiClient
	Error     error
	Function  *models.Function
	Name      string
}

func (m DeleteModel) DeleteFunction(functionName string) tea.Msg {
	function, err := m.ApiClient.DeleteFunction(functionName)
	if err != nil {
		return ErrorMsg{
			Error: err,
		}
	}
	return FunctionDeleteResponseMsg{
		Function: *function,
	}
}

func (m DeleteModel) IsDeleted() bool {
	return m.Function != nil
}

func (m DeleteModel) Init() tea.Cmd {
	return tea.Batch(m.Spinner.Tick,
		func() tea.Msg {
			return m.DeleteFunction(m.Name)
		})
}

func (m DeleteModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

	case FunctionDeleteResponseMsg:
		m.Function = &msg.Function
	case ErrorMsg:
		m.Error = msg.Error
	}

	return m, cmd
}

func (m DeleteModel) View() string {
	if m.Error != nil {
		return printError(m.Error)
	}
	if !m.IsDeleted() {
		return fmt.Sprintf("%s %s", m.Spinner.View(), "Deleting function..."+boldString("Press CTRL+C or q to quit"))
	}

	return printSuccess("Deleted function " + m.Function.Name)
}
