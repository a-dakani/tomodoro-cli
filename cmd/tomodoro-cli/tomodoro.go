package main

import (
	"context"
	"github.com/a-dakani/tomodoro-cli/pkg/tclient"
)

func getOrCreateTeam(teamName string) (Team, error) {
	ctx := context.Background()
	t, err := tc.GetTeam(ctx, teamName)
	if err != nil {
		if err.(*tclient.RequestError).NotFound() {
			_, err = tc.CreateTeam(ctx, teamName)
			if err != nil {
				return Team{}, err
			}
			t, err = tc.GetTeam(ctx, teamName)
			if err != nil {
				return Team{}, err
			}
		} else {
			return Team{}, err
		}
	}

	return Team{
		Name:  t.Name,
		Slug:  t.Slug,
		Focus: t.Settings.Focus,
		Pause: t.Settings.Pause,
	}, nil
}

func startFocus(team Team) error {
	ctx := context.Background()

	_, err := tc.StartTimer(ctx, team.Slug, team.Focus, "Focus")
	if err != nil {
		return err
	}

	return nil
}
func startPause(team Team) error {
	ctx := context.Background()

	_, err := tc.StartTimer(ctx, team.Slug, team.Pause, "Pause")
	if err != nil {
		return err
	}

	return nil
}
func stopTimer(team Team) error {
	ctx := context.Background()

	_, err := tc.StopTimer(ctx, team.Slug)
	if err != nil {
		return err
	}

	return nil
}
