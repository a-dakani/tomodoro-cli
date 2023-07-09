package main

import (
	"context"
	"github.com/a-dakani/tomodoro-cli/pkg/config"
)

func getTeam(teamName string) (config.Team, error) {
	ctx := context.Background()
	t, err := tc.GetTeam(ctx, teamName)

	if err != nil {
		return config.Team{}, err
	}

	return config.Team{
		Name:  t.Name,
		Slug:  t.Slug,
		Focus: t.Settings.Focus,
		Pause: t.Settings.Pause,
	}, nil
}

func startFocus(team config.Team) error {
	ctx := context.Background()

	_, err := tc.StartTimer(ctx, team.Slug, team.Focus, "Focus")
	if err != nil {
		return err
	}

	return nil
}
func startPause(team config.Team) error {
	ctx := context.Background()

	_, err := tc.StartTimer(ctx, team.Slug, team.Pause, "Pause")
	if err != nil {
		return err
	}

	return nil
}
func stopTimer(team config.Team) error {
	ctx := context.Background()

	_, err := tc.StopTimer(ctx, team.Slug)
	if err != nil {
		return err
	}

	return nil
}
