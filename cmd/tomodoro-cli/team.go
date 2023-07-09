package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

var teamList *TeamList

type TeamList []Team

type Team struct {
	Name  string `json:"name"`
	Slug  string `json:"slug"`
	Focus int64  `json:"focus"`
	Pause int64  `json:"pause"`
}

func LoadTeams() *TeamList {
	if teamList != nil {
		return teamList
	}
	teamList = &TeamList{}
	teamList.load()
	return teamList
}

func (tl *TeamList) load() {
	_, err := os.Stat(cfg.TeamsFilePath)
	if os.IsNotExist(err) {
		tl.init()
		tl.load()
	}
	bytes, err := os.ReadFile(cfg.TeamsFilePath)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(bytes, tl)
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
	if _, err := os.Stat(cfg.TeamsFilePath); errors.Is(err, os.ErrNotExist) {
		err = os.MkdirAll(cfg.ConfigPath, os.ModePerm)
		if err != nil {
			return err
		}
		_, err = os.Create(cfg.TeamsFilePath)
		if err != nil {
			return err
		}
	}

	bytes, err := json.MarshalIndent(tl, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(cfg.TeamsFilePath, bytes, 0600)
	if err != nil {
		return err
	}

	return nil
}

func (tl *TeamList) TeamExists(team Team) bool {
	for _, t := range *tl {
		if t.Slug == team.Slug {
			return true
		}
	}
	return false
}

func (tl *TeamList) AddTeam(team Team) error {
	// check if team already exists if it does replace it with the newModel one else append it
	for i, t := range *tl {
		if t.Slug == team.Slug {
			(*tl)[i] = team
			return tl.Save()
		}
	}
	*tl = append(*tl, team)
	return tl.Save()
}

func (tl *TeamList) RemoveTeam(team Team) error {
	for i, t := range *tl {
		if t.Slug == team.Slug {
			*tl = append((*tl)[:i], (*tl)[i+1:]...)
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
