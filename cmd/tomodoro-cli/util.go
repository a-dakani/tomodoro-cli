package main

import (
	"encoding/json"
	"errors"
	"os"
)

func readTeamsFile() ([]Team, error) {
	var teams []Team
	if _, err := os.Stat(teamsFile); errors.Is(err, os.ErrNotExist) {
		err = os.MkdirAll(configPath, os.ModePerm)
		if err != nil {
			return teams, err
		}
		_, err = os.Create(teamsFile)
		if err != nil {
			return nil, err
		}

		bytes, err := json.MarshalIndent(teams, "", "  ")
		if err != nil {
			return teams, err
		}

		err = os.WriteFile(teamsFile, bytes, os.ModePerm)

		return teams, err
	}

	bytes, err := os.ReadFile(teamsFile)
	if err != nil {
		return teams, err
	}

	err = json.Unmarshal(bytes, &teams)
	if err != nil {
		return teams, err
	}

	return teams, nil
}

func removeTeamFromFile(team Team) error {
	teams, err := readTeamsFile()
	if err != nil {
		return err
	}

	for i, t := range teams {
		if t.Slug == team.Slug {
			teams = append(teams[:i], teams[i+1:]...)
			break
		}
	}

	bytes, err := json.MarshalIndent(teams, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(teamsFile, bytes, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

func addTeamToFile(team Team) error {
	teamExists := false

	fileTeams, err := readTeamsFile()
	if err != nil {
		return err
	}

	// check if team already exists if it does replace it with the newModel one else append it
	for i, t := range fileTeams {
		if t.Slug == team.Slug {
			fileTeams[i] = team
			teamExists = true

			break
		}
	}

	if !teamExists {
		fileTeams = append(fileTeams, team)
	}

	bytes, err := json.MarshalIndent(fileTeams, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(teamsFile, bytes, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

func getConfigFilePath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	return homeDir + "/.config/tomodoro"
}
