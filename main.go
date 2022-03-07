package main

import (
	"fmt"
	"log"
	"time"

	"github.com/nicklaw5/helix/v2"
)

func main() {
	config, err := ParseConfig()
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(config)

	/*
	access_token, token_type, err := GetOAuth(config.Twitch.Id)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(access_token, token_type)
	*/

	access_token, _ := "xef09nmfoo4j4ilmnfnkjlke40vvys", "bearer"

	twitch_client, err := helix.NewClient(&helix.Options{ClientID: config.Twitch.Id, UserAccessToken: access_token})
	if err != nil {
		log.Fatalln(err)
	}
	
	discord_bot, err := DiscordBotInitialize(config.Discord.Webhook)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(discord_bot)

	/*
	err = discord_bot.Send("Hello World!")
	if err != nil {
		log.Fatalln(err)
	}
	*/

	for {
		resp, err := twitch_client.GetStreams(&helix.StreamsParams{First: 1, UserLogins: []string{config.Twitch.User}})
		if err != nil {
			log.Fatalln(err, resp.Error, resp.ErrorMessage)
		}
		if resp.StatusCode != 200 {
			log.Fatalln("Unexpected Status Code: ", resp.StatusCode)
		}

		fmt.Println("Remaining:", resp.GetRateLimitRemaining())
		if resp.GetRateLimitRemaining() == 0 {
			time.Sleep(60 * time.Second)
			continue
		}

		fmt.Println(resp.Data)

		time.Sleep(5 * time.Second)
	}
}