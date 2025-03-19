package singleInput

import (
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type (
	errMsg error
)

type InputField struct {
	Placeholder        string
	ValidationFunction textinput.ValidateFunc
	Value              string
	Exit               bool
}

type model struct {
	textInput textinput.Model
	err       error
	in        *InputField
}

func InitialModel(in *InputField) model {
	ti := textinput.New()
	ti.Prompt = fmt.Sprintf("")
	
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 3

	return model{
		textInput: ti,
		err:       nil,
		in:        in,
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			return m, tea.Quit
		case tea.KeyCtrlC, tea.KeyEsc:
			m.in.Exit = true
			return m, tea.Quit
		}

	// We handle errors just like any other message
	case errMsg:
		m.err = msg
		return m, nil
	}

	m.textInput, cmd = m.textInput.Update(msg)
	m.in.Value = m.textInput.Value()

	return m, cmd
}

func (m model) View() string {
	return fmt.Sprintf(
		"%-30s %s",
		m.in.Placeholder+":",
		m.textInput.View(),
	) + "\n"
}
