package main

import (
	"fmt"
	"github.com/a-dakani/tomodoro-cli/pkg/tclient"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// errorMsg is a wrapper for error to implement the tea.Msg interface
type errorMsg error

type sessionState int

const (
	noTeams sessionState = iota
	showList
	showTimer
	showInput
)

type model struct {
	title          string
	state          sessionState
	sub            chan tclient.Message
	ws             *tclient.WebSocketClient
	input          textinput.Model
	timerName      string
	timerState     tclient.MessageType
	timerRemaining int64
	teamList       list.Model
	help           help.Model
	height         int
	width          int
	err            error
	windowTooSmall bool
}

func newModel() *model {
	ti := textinput.New()
	ti.CharLimit = 30
	ti.Placeholder = "Team Slug"
	ti.Focus()

	delegate := list.NewDefaultDelegate()
	delegate.SetHeight(3)
	delegate.ShortHelpFunc = func() []key.Binding { return []key.Binding{Keymap.Add, Keymap.Remove} }
	delegate.FullHelpFunc = func() [][]key.Binding { return [][]key.Binding{{Keymap.Add, Keymap.Remove}} }
	tl := list.New([]list.Item{}, delegate, initialListHeight, initialListWidth)
	tl.Title = "Teams"

	return &model{
		title:          "Tomodoro",
		state:          noTeams,
		sub:            make(chan tclient.Message, 100),
		input:          ti,
		timerName:      "Inactive",
		timerState:     tclient.TimerStopped,
		timerRemaining: 0,
		teamList:       tl,
		help:           help.New(),
		height:         minimalWindowHeight,
		width:          minimalWindowWidth,
		err:            nil,
		windowTooSmall: false,
	}
}

func (m *model) Init() tea.Cmd {
	m.loadTeams()
	return nil
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	var cmd tea.Cmd

	// if msg type is tea.KeyMsg and is CtrlC, return tea.Quit
	if msg, ok := msg.(tea.KeyMsg); ok {
		if key.Matches(msg, Keymap.Quit) {
			return m, tea.Quit
		}
	}
	// if msg type is tea.WindowSizeMsg, set the width and height of the model
	if msg, ok := msg.(tea.WindowSizeMsg); ok {
		// if the window is too small, set the windowTooSmall flag to true
		if msg.Width < minimalWindowWidth || msg.Height < minimalWindowHeight {
			m.windowTooSmall = true
		} else {
			m.windowTooSmall = false
		}

		m.height = msg.Height - stylesHeight - 1
		m.width = msg.Width - stylesWidth
		m.help.Width = msg.Width - stylesWidth - 2
	}

	// if the window is too small, return the model and don't update
	if m.windowTooSmall {
		return m, nil
	}
	// else check the state of the model
	switch m.state {
	case noTeams:
		if msg, ok := msg.(tea.KeyMsg); ok {
			if key.Matches(msg, Keymap.Add) {
				m.input.Reset()
				m.state = showInput
			}
		}
	case showInput:
		switch msg := msg.(type) {
		case Team:
			if err := teams.AddTeam(msg); err != nil {
				m.err = err
				return m, nil
			}
			// reload the list of teams
			m.loadTeams()
			m.teamList.SetWidth(m.width)
			m.teamList.SetHeight(m.height)
			m.teamList.Help.Width = m.width - stylesWidth
			m.state = showList
			m.teamList, _ = m.teamList.Update(msg)
		case tea.KeyMsg:
			switch {
			case msg.Type == tea.KeyEnter:
				// reset the error in case it was set
				m.err = nil
				return m, m.addTeam()
			case key.Matches(msg, Keymap.Back):
				m.err = nil
				if len(m.teamList.Items()) == 0 {
					m.state = noTeams
				} else {
					m.teamList.SetWidth(m.width)
					m.teamList.SetHeight(m.height)
					m.teamList.Help.Width = m.width - stylesWidth
					m.state = showList
					m.teamList, _ = m.teamList.Update(msg)
				}
			}
		}

		m.input, cmd = m.input.Update(msg)
	case showList:
		if msg, ok := msg.(tea.KeyMsg); ok {
			switch {
			case msg.Type == tea.KeyEnter:
				m.state = showTimer
				return m, tea.Batch(m.joinTeam(), m.waitForActivity())
			case key.Matches(msg, Keymap.Add):
				m.input.Reset()
				m.state = showInput
			case key.Matches(msg, Keymap.Remove):
				if err := teams.RemoveTeam(m.teamList.Items()[m.teamList.Index()].(Team)); err != nil {
					m.err = err
				}

				m.teamList.RemoveItem(m.teamList.Index())

				if len(m.teamList.Items()) == 0 {
					m.state = noTeams
				}
			}
		}

		m.teamList.SetWidth(m.width)
		m.teamList.SetHeight(m.height)
		m.teamList.Help.Width = m.width - stylesWidth
		m.teamList, cmd = m.teamList.Update(msg)
	case showTimer:
		switch msg := msg.(type) {
		case tclient.Message:
			switch msg.Type {
			case tclient.Tick:
				if m.timerState != tclient.TimerStarted {
					m.timerState = tclient.TimerStarted
				}

				m.timerName = msg.Payload.Name
				m.timerRemaining = msg.Payload.RemainingTime
			case tclient.TimerStarted:
				m.timerState = tclient.TimerStarted
				m.timerName = msg.Payload.Name
				n.Notify(fmt.Sprintf("Timer %s started", msg.Payload.Name))
			case tclient.TimerStopped:
				m.timerRemaining = 0
				m.timerName = "Inactive"
				m.timerState = tclient.TimerStopped
				n.Notify(fmt.Sprintf("Timer %s stopped", msg.Payload.Name))
			case tclient.Connecting:
				m.timerState = tclient.Connecting
			case tclient.Connected:
				m.timerState = tclient.Connected
			case tclient.Terminating:
				m.timerRemaining = 0
				m.timerName = "Inactive"
				m.timerState = tclient.Terminating
			case tclient.Error:
				m.timerState = tclient.Error
				m.err = msg.Error
			}

			return m, m.waitForActivity()
		case tea.KeyMsg:
			switch {
			case msg.Type == tea.KeyEnter:
				m.state = showTimer
			case key.Matches(msg, Keymap.Back):
				m.state = showList
			case key.Matches(msg, Keymap.StartFocus):
				m.err = nil
				return m, m.startFocus()
			case key.Matches(msg, Keymap.StartPause):
				m.err = nil
				return m, m.startPause()
			case key.Matches(msg, Keymap.StopTimer):
				m.err = nil
				return m, m.stopTimer()
			case key.Matches(msg, Keymap.ShowFullHelp):
				m.help.ShowAll = !m.help.ShowAll
			}
		}
	}

	return m, cmd
}

func (m *model) View() string {
	var output string
	output += m.renderTitle()

	if m.windowTooSmall {
		t := fmt.Sprintf(
			"Window too small. Please resize.\n\nMinimum width: %d\nMinimum height: %d",
			minimalWindowWidth, minimalWindowHeight)
		output += addHelp(t, "q/ctrl quit", m.height)

		return appStyle.Width(m.width).Height(m.height + 1).Render(output)
	}

	if m.err != nil {
		output += m.err.Error() + "\n"
	}

	switch m.state {
	case showList:
		output += m.teamList.View()
	case showTimer:
		t := renderTimer(m.teamList.SelectedItem().(Team), m.timerRemaining, m.timerName, string(m.timerState))
		output += addHelp(t, m.help.View(Keymap), m.height)
	case showInput:
		output += m.input.View()
	case noTeams:
		output += "No teams found. to fetch a team from tomodoro press `+`"
	default:
		output += "Something went wrong. Please try again."
	}

	return appStyle.Width(m.width).Height(m.height).Render(output)
}

func (m *model) loadTeams() {
	items := make([]list.Item, len(*teams))
	for i, team := range *teams {
		items[i] = team
	}

	if len(items) != 0 {
		m.state = showList
	}

	m.teamList.SetItems(items)
}

func (m *model) addTeam() tea.Cmd {
	return func() tea.Msg {
		team, err := getTeam(m.input.Value())
		if err != nil {
			m.err = err
			return errorMsg(err)
		}
		return team
	}
}

func (m *model) startFocus() tea.Cmd {
	return func() tea.Msg {
		err := startFocus(m.teamList.Items()[m.teamList.Index()].(Team))
		if err != nil {
			return errorMsg(err)
		}

		return m.waitForActivity()
	}
}

func (m *model) startPause() tea.Cmd {
	return func() tea.Msg {
		err := startPause(m.teamList.Items()[m.teamList.Index()].(Team))
		if err != nil {
			return errorMsg(err)
		}

		return m.waitForActivity()
	}
}

func (m *model) stopTimer() tea.Cmd {
	return func() tea.Msg {
		err := stopTimer(m.teamList.Items()[m.teamList.Index()].(Team))
		if err != nil {
			return errorMsg(err)
		}

		return m.waitForActivity()
	}
}

func (m *model) joinTeam() tea.Cmd {
	return func() tea.Msg {
		slug := m.teamList.SelectedItem().(Team).Slug
		// if there is already a websocket connection, check if it is the same team
		if m.ws != nil {
			if m.ws.Slug == slug {
				return nil
			}

			m.ws.Stop()
			m.ws = tclient.NewWebSocketClient(cfg.BaseWSURLV1, m.teamList.SelectedItem().(Team).Slug)
			m.ws.Start()

			for {
				for elem := range m.ws.OutChan {
					m.sub <- elem
				}
			}
		}

		m.ws = tclient.NewWebSocketClient(cfg.BaseWSURLV1, m.teamList.SelectedItem().(Team).Slug)
		m.ws.Start()

		for {
			for {
				for elem := range m.ws.OutChan {
					m.sub <- elem
				}
			}
		}
	}
}

func (m *model) waitForActivity() tea.Cmd {
	return func() tea.Msg {
		return <-m.sub
	}
}

func (m *model) renderTitle() string {
	return titleStyle.Width(m.width-2).Render(m.title) + "\n"
}
