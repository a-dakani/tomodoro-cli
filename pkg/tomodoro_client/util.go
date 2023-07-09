package tomodoro_client

import "os"

func getTeamsFilePath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	return homeDir + "/.config/tomodoro-cli/teams.json"
}
