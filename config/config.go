package config

import (
	"encoding/json"
	"io/ioutil"
)

type UntappdConfig struct {
	ClientID     string
	ClientSecret string
}

type Config struct {
	Debug           bool
	Room            string
	User            string
	Nick            string
	Server          string
	CommandPrefix   string
	Untappd         *UntappdConfig
	QuotesDirectory string
	TellFile        string
}

// Default returns the default config.
func Default() Config {
	conf := Config{
		Debug:         false,
		Room:          "#test",
		Nick:          "glitz",
		User:          "glitz",
		Server:        "irc.rizon.net:6667",
		CommandPrefix: ".",
		Untappd: &UntappdConfig{
			ClientID:     "",
			ClientSecret: "",
		},
		QuotesDirectory: "_quotes",
		TellFile:        "_data/tell.json",
	}
	return conf
}

// Load loads the config from the specified json file.
func Load(filename string) (Config, error) {
	var config Config
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return config, err
	}
	err = json.Unmarshal(content, &config)
	return config, err
}
