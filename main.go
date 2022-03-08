package main

import (
	"fmt"
	"log"
	"time"

	"github.com/nicklaw5/helix/v2"
)

const (
	REQUIRE_DELTA = 2 // How many requests are nessessary to change online -> offline or vice versa.
	// This is in-case we hit two different servers, one with a more up-to-date cache
	// than the other.
)

func main() {
	config, err := ParseConfig()
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(config)

	access_token, _, err := GetOAuth(config.Twitch.Id)
	if err != nil {
		log.Fatalln(err)
	}

	twitch_client, err := helix.NewClient(&helix.Options{ClientID: config.Twitch.Id, UserAccessToken: access_token})
	if err != nil {
		log.Fatalln(err)
	}

	discord_bot, err := DiscordBotInitialize(config.Discord.Webhook)
	if err != nil {
		log.Fatalln(err)
	}

	time_since_offline, time_since_online := 0, 0
	for {
		resp, err := twitch_client.GetStreams(&helix.StreamsParams{First: 1, UserLogins: []string{config.Twitch.User}})
		if err != nil {
			log.Fatalln(err, resp.Error, resp.ErrorMessage)
		}

		if resp.GetRateLimitRemaining() == 0 {
			time.Sleep(60 * time.Second) // Wait a minute for requests to reset, this isn't nessessary
			// because we get about 900 requests/minute, but just in case
			continue
		}

		if len(resp.Data.Streams) > 0 { // Online
			if time_since_online >= REQUIRE_DELTA {
				log.Println("Stream Online")
			}

			time_since_online = 0
			time_since_offline += 1
		} else { // Offline
			if time_since_offline >= REQUIRE_DELTA {
				log.Println("Stream Offline")

				discord_bot.Send(config.Discord.Message)
			}

			time_since_offline = 0
			time_since_online += 1
		}

		time.Sleep(5 * time.Second) // Picked 5 seconds semi-arbitrarily to avoid DoS-ing Twitch
	}
}
