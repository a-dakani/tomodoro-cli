package main

import (
	"github.com/a-dakani/tomodoro-cli/pkg/config"
	"github.com/a-dakani/tomodoro-cli/pkg/notifier"
	"github.com/a-dakani/tomodoro-cli/pkg/tclient"
	tea "github.com/charmbracelet/bubbletea"
	"log"
)

var (
	cfg = config.LoadConfig()
	tl  = NewTeamList(cfg.ConfigPath, cfg.TeamsFilePath)
	n   = notifier.NewNotifier(cfg.NotificationTitle, cfg.NotificationImagePath)
	tc  = tclient.NewHttpClient(cfg.BaseURLV1)
)

func main() {
	m := newModel()
	program := tea.NewProgram(m, tea.WithAltScreen())

	_, err := program.Run()
	if err != nil {
		log.Fatal(err)
	}
}
