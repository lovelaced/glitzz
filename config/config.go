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

type SeenConfig struct {
	SeenFile string
}

type RemindersConfig struct {
	RemindersFile string
}

type Config struct {
	Rooms          []string
	User           string
	Nick           string
	Server         string
	TLS            bool
	CommandPrefix  string
	EnabledModules []string
	Untappd        UntappdConfig
	Quotes         QuotesConfig
	Tell           TellConfig
	Seen           SeenConfig
	Reminders      RemindersConfig
}

// Default returns the default config.
func Default() Config {
	conf := Config{
		Rooms:         []string{"#test"},
		Nick:          "glitz",
		User:          "glitz",
		Server:        "irc.rizon.net:6697",
		TLS:           true,
		CommandPrefix: ".",
		EnabledModules: []string{
			"c3",
			"fourchan",
			"info",
			"links",
			"pipes",
			"quotes",
			"reactions",
			"reminders",
			"seen",
			"tell",
			"vatsim",
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
		Seen: SeenConfig{
			SeenFile: "_data/seen.json",
		},
		Reminders: RemindersConfig{
			RemindersFile: "_data/reminders.json",
		},
	}
	return conf
}

// Load loads the config from the specified json file.
func Load(filename string) (Config, error) {
	var config Config = Default()
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return config, err
	}
	err = json.Unmarshal(content, &config)
	return config, err
}
