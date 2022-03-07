package main

import (
	"os"

	yaml "gopkg.in/yaml.v2"
)

type BotConfig struct {
	Discord struct {
		Webhook string `yaml:"webhook"`
	}

	Twitch struct {
		Id string `yaml:"id"`
		
		User string `yaml:"user"`
	}
}

type ConfigError struct {
	reason string
}

func (e ConfigError) Error() string {
	return e.reason
}

func ParseConfig() (*BotConfig, error) {
	data, err := os.ReadFile("./config.yml")
	if err != nil {
		return nil, err
	}

	config := BotConfig{}
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	if config.Discord.Webhook == "" {
		return nil, ConfigError{ "discord.webhook required, but not included" }
	}

	if config.Twitch.Id == "" {
		return nil, ConfigError{ "twitch.id required, but not included" }
	}

	if config.Twitch.User == "" {
		return nil, ConfigError{ "twitch.user required, but not included" }
	}

	return &config, nil
}