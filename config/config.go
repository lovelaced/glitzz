package config

import (
	"encoding/json"
	"io/ioutil"
)

type UntappdConfig struct {
	ClientID     string
	ClientSecret string
}

type QuotesConfig struct {
	QuotesDirectory string
}

type TellConfig struct {
	TellFile string
}

type Config struct {
	Debug          bool
	Room           string
	User           string
	Nick           string
	Server         string
	CommandPrefix  string
	EnabledModules []string
	Untappd        UntappdConfig
	Quotes         QuotesConfig
	Tell           TellConfig
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
		EnabledModules: []string{
			"c3",
			"fourchan",
			"info",
			"pipes",
			"quotes",
			"reactions",
			"tell",
		},
		Untappd: UntappdConfig{
			ClientID:     "",
			ClientSecret: "",
		},
		Quotes: QuotesConfig{
			QuotesDirectory: "_quotes",
		},
		Tell: TellConfig{
			TellFile: "_data/tell.json",
		},
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
