package config

import (
	"encoding/json"
	"errors"
	"os"
	"path"
	"time"
)

const (
	baseURLV1      = "https://api.tomodoro.de/api/v1/"
	baseWSURLV1    = "wss://api.tomodoro.de/api/v1/"
	teamsFileName  = "teamList.json"
	configFileName = "config.json"
)

var (
	configPath     = getConfigFilePath()
	teamsFilePath  = path.Join(configPath, teamsFileName)
	configFilePath = path.Join(configPath, configFileName)
)

type Config struct {
	BaseURLV1         string        `json:"base_url_v1"`
	BaseWSURLV1       string        `json:"base_ws_url_v1"`
	HTTPClientTimeout time.Duration `json:"http_client_timeout"`
}

var cfg *Config

func LoadConfig() *Config {
	if cfg != nil {
		return cfg
	}
	cfg = &Config{}

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
