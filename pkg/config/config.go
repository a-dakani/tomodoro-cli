package config

import (
	"encoding/json"
	"errors"
	"os"
	"path"
	"time"
)

const (
	baseURLV1                 = "https://api.tomodoro.de/api/v1/"
	baseWSURLV1               = "wss://api.tomodoro.de/api/v1/"
	teamsFileName             = "teamList.json"
	configFileName            = "config.json"
	notificationTitle         = "Tomodoro"
	notificationSoundFileName = "default.mp3"
	notificationImageFileName = "logo.png"
)

var (
	configPath     = getConfigFilePath()
	teamsFilePath  = path.Join(configPath, teamsFileName)
	configFilePath = path.Join(configPath, configFileName)
)

type Config struct {
	BaseURLV1             string        `json:"base_url_v1"`
	BaseWSURLV1           string        `json:"base_ws_url_v1"`
	HTTPClientTimeout     time.Duration `json:"http_client_timeout"`
	ConfigPath            string        `json:"config_path"`
	ConfigFilePath        string        `json:"config_file_path"`
	TeamsFilePath         string        `json:"teams_file_path"`
	NotificationTitle     string        `json:"notification_title"`
	NotificationSoundPath string        `json:"notification_sound"`
	NotificationImagePath string        `json:"notification_image"`
}

var cfg *Config

func LoadConfig() *Config {
	if cfg != nil {
		return cfg
	}
	cfg = &Config{
		ConfigPath:     configPath,
		ConfigFilePath: configFilePath,
		TeamsFilePath:  teamsFilePath,
	}

	cfg.load()
	return cfg
}
func (c *Config) load() {
	_, err := os.Stat(configFilePath)
	if os.IsNotExist(err) {
		c.init()
		c.load()
	}
	// load config file
	bytes, err := os.ReadFile(configFilePath)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(bytes, c)
	if err != nil {
		panic(err)
	}
}

func (c *Config) init() {
	c.BaseURLV1 = baseURLV1
	c.BaseWSURLV1 = baseWSURLV1
	c.HTTPClientTimeout = time.Minute
	c.NotificationTitle = notificationTitle
	c.NotificationSoundPath = path.Join(configPath, notificationSoundFileName)
	c.NotificationImagePath = path.Join(configPath, notificationImageFileName)

	err := c.save()
	if err != nil {
		panic(err)
	}
}

func (c *Config) save() error {
	if _, err := os.Stat(configFilePath); errors.Is(err, os.ErrNotExist) {
		err = os.MkdirAll(configPath, os.ModePerm)
		if err != nil {
			return err
		}
		_, err = os.Create(configFilePath)
		if err != nil {
			return err
		}
	}

	bytes, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	err = os.WriteFile(configFilePath, bytes, 0600)
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
