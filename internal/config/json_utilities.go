package cfg

import (
	"encoding/json"
	"os"
)

// const for ressource handling
const configFileName = "/.gatorconfig.json"

// helper function for ressource handling
func getConfigFilePath() (string, error) {
	homePath, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return homePath + configFileName, nil
}

// struct to represent JSON file structure including tags
type Config struct {
	DbURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

// function to read .gatorconfig.json 
func Read() (Config, error) {
	configFilePath, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}
	data, err := os.ReadFile(configFilePath)
	if err != nil {
		return Config{}, err
	}
	config := Config{}
	if err = json.Unmarshal(data, &config); err != nil {
		return Config{}, err
	}
	
	return config, nil
}

// private helper function for ressource handling
func write(c Config) error {
	configFilePath, err := getConfigFilePath()
	if err != nil {
		return err
	}
	data, err := json.Marshal(c)
	if err != nil {
		return err
	}
	err = os.WriteFile(configFilePath, data, 0600)
	if err != nil {
		return err
	}

	return nil
}

// method to write config struct to .gatorconfig.json
func (c *Config) SetUser(name string) error {
	c.CurrentUserName = name
	return write(*c)
}

