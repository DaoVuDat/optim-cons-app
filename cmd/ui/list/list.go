package list

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"golang-moaha-construction/internal/data"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

func InitializeList(
	output *Output,
	title string,
	items []list.Item,
	delegate list.ItemDelegate, width, height int,
) model {

	m := model{
		list:   list.New(items, delegate, width, height),
		output: output,
	}

	m.list.Title = title
	return m
}

type Output struct {
	Output *data.Item
	Exit   bool
}

type model struct {
	list   list.Model
	output *Output
	exit   bool
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			i, ok := m.list.SelectedItem().(data.Item)

			if ok {
				newItem := data.NewItem(i.Title(), i.Description())
				m.output.Output = &newItem
			}
			return m, tea.Quit
		case tea.KeyCtrlC, tea.KeyEsc:
			m.exit = true
			m.output.Exit = true
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return docStyle.Render(m.list.View())
}
