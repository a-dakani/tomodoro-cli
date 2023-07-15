package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

var teamList *TeamList

type TeamList struct {
	Teams      []Team `json:"tl"`
	ConfigPath string `json:"config_path"`
	FilePath   string `json:"file_path"`
}

type Team struct {
	Name  string `json:"name"`
	Slug  string `json:"slug"`
	Focus int64  `json:"focus"`
	Pause int64  `json:"pause"`
}

func NewTeamList(configPath, filePath string) *TeamList {
	if teamList != nil {
		return teamList
	}
	teamList = &TeamList{
		ConfigPath: configPath,
		FilePath:   filePath,
		Teams:      []Team{},
	}
	teamList.load()
	return teamList
}

func (tl *TeamList) load() {
	_, err := os.Stat(tl.FilePath)
	if os.IsNotExist(err) {
		tl.init()
		tl.load()
	}
	bytes, err := os.ReadFile(tl.FilePath)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(bytes, &tl.Teams)
	if err != nil {
		panic(err)
	}
}

func (tl *TeamList) init() {
	// create teamList file with empty array
	err := tl.Save()
	if err != nil {
		panic(err)
	}
}

func (tl *TeamList) Save() error {
	if _, err := os.Stat(tl.FilePath); errors.Is(err, os.ErrNotExist) {
		err = os.MkdirAll(tl.ConfigPath, os.ModePerm)
		if err != nil {
			return err
		}
		_, err = os.Create(tl.FilePath)
		if err != nil {
			return err
		}
	}

	bytes, err := json.MarshalIndent(tl.Teams, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(tl.FilePath, bytes, 0600)
	if err != nil {
		return err
	}

	return nil
}

func (tl *TeamList) TeamExists(team Team) bool {
	for _, t := range tl.Teams {
		if t.Slug == team.Slug {
			return true
		}
	}
	return false
}

func (tl *TeamList) AddTeam(team Team) error {
	// check if team already exists if it does replace it with the newModel one else append it
	for i, t := range tl.Teams {
		if t.Slug == team.Slug {
			tl.Teams[i] = team
			return tl.Save()
		}
	}
	tl.Teams = append(tl.Teams, team)
	return tl.Save()
}

func (tl *TeamList) RemoveTeam(team Team) error {
	for i, t := range tl.Teams {
		if t.Slug == team.Slug {
			tl.Teams = append(tl.Teams[:i], tl.Teams[i+1:]...)
			return tl.Save()
		}
	}
	return nil
}

// FilterValue returns the value to filter the teamList by
func (t Team) FilterValue() string {
	return t.Name
}

// Title returns the title of the teamList item
func (t Team) Title() string {
	return t.Slug
}

// Description returns the description of the teamList item
func (t Team) Description() string {
	return fmt.Sprintf("Focus: %d min\nPause: %d min", t.Focus/1000000000/60, t.Pause/1000000000/60)
}
