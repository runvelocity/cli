package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/runvelocity/cli/internal/api"
	"github.com/runvelocity/cli/internal/models"
)

type FunctionListResponseMsg struct {
	Functions []models.Function
}

type ErrorMsg struct {
	Error error
}

type ListModel struct {
	Functions []models.Function
	Spinner   spinner.Model
	ApiClient api.ApiClient
	Error     error
}

func (m ListModel) ListFunctions() tea.Msg {
	functions, err := m.ApiClient.ListFunctions()
	if err != nil {
		return ErrorMsg{
			Error: err,
		}
	}
	return FunctionListResponseMsg{
		Functions: *functions,
	}
}

func (m ListModel) IsInitialized() bool {
	return m.Functions != nil
}

func (m ListModel) Init() tea.Cmd {
	return tea.Batch(m.Spinner.Tick,
		m.ListFunctions)
}

func (m ListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

	case FunctionListResponseMsg:
		m.Functions = msg.Functions
	case ErrorMsg:
		m.Error = msg.Error
	}

	return m, cmd
}

func (m ListModel) View() string {
	if m.Error != nil {
		return printError(m.Error)
	}
	if !m.IsInitialized() {
		return fmt.Sprintf("%s %s", m.Spinner.View(), "Fetching functions..."+boldString("Press CTRL+C or q to quit"))
	}

	columns := []table.Column{
		{Title: "UUID", Width: 20},
		{Title: "Name", Width: 10},
		{Title: "Handler", Width: 20},
	}

	rows := []table.Row{}

	for _, v := range m.Functions {
		rows = append(rows, table.Row{v.UUID, v.Name, v.Status, v.Handler})
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithHeight(5),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Cell.Foreground(color)
	t.SetStyles(s)
	t.Focus()
	return boldString("Press CTRL+C or q to quit") + "\n\n" + t.View()
}
