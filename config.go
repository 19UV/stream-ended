package main

import (
	"os"
	"strings"
	"text/template"

	yaml "gopkg.in/yaml.v2"
)

const (
	DISCORD_WEBHOOK_DEFAULT = "<INSERT DISCORD WEBHOOK>"
	DISCORD_MESSAGE_DEFAULT = "@everyone {{.User}} just went Offline!"
	TWITCH_APP_ID_DEFAULT   = "<INSERT TWITCH APPLICATION ID>"
	TWITCH_CHANNEL_DEFAULT  = "19UV"
)

type BotConfig struct {
	Discord struct {
		Webhook string `yaml:"webhook"`
		Message string `yaml:"message"`
	}

	Twitch struct {
		Id   string `yaml:"id"`
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
	config.Discord.Webhook = DISCORD_WEBHOOK_DEFAULT
	config.Discord.Message = DISCORD_MESSAGE_DEFAULT
	config.Twitch.Id = TWITCH_APP_ID_DEFAULT
	config.Twitch.User = TWITCH_CHANNEL_DEFAULT

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	if config.Discord.Webhook == DISCORD_WEBHOOK_DEFAULT {
		return nil, ConfigError{"discord.webhook required, but not included"}
	}

	if config.Twitch.Id == TWITCH_APP_ID_DEFAULT {
		return nil, ConfigError{"twitch.id required, but not included"}
	}

	if config.Discord.Message == DISCORD_MESSAGE_DEFAULT {
		config.Discord.Message = "@everyone {{.User}} just went Offline!"
	}

	template, err := template.New("message").Parse(config.Discord.Message)
	if err != nil {
		return nil, err
	}

	buf := new(strings.Builder)
	err = template.Execute(buf, map[string]string{"User": config.Twitch.User})
	if err != nil {
		return nil, err
	}
	config.Discord.Message = buf.String()

	return &config, nil
}
