package main

import "path"

var (
	teamsFile     = path.Join(configPath, teamsFileName)
	configPath    = getConfigFilePath()
	teamsFileName = "teams.json"
)
