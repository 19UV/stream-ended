package main

import (
	"fmt"
	"log"

	helix "github.com/nicklaw5/helix/v2"
)

func main() {
	config, err := ParseConfig()
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(config)

	access_token, token_type, err := GetOAuth(config.Twitch.Id)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(access_token, token_type)

	/*
	discord_bot, err := DiscordBotInitialize(config.Discord.Webhook)
	if err != nil {
		log.Fatalln(err)
	}

	err = discord_bot.Send("Hello World!")
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(discord_bot)
	*/
}