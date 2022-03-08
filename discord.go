package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type DiscordWebhookObject struct {
	Type int `json:"type"`

	Id    string `json:"id"`
	Token string `json:"token"`

	Name string `json:"name"`
}

type DiscordSendObject struct {
	Username  string `json:"username"`
	AvatarURL string `json:"avatar_url"`
	Content   string `json:"content"`
}

type DiscordBot struct {
	Url string
	Obj *DiscordWebhookObject
}

func (bot *DiscordBot) Get() (*DiscordWebhookObject, error) {
	resp, err := http.Get(bot.Url)
	if err != nil { // TODO: Check Status Code
		return nil, err
	}

	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	resp_object := DiscordWebhookObject{}
	err = json.Unmarshal(buf, &resp_object)
	if err != nil {
		return nil, err
	}

	err = resp.Body.Close()
	if err != nil {
		return nil, err
	}

	bot.Obj = &resp_object

	return &resp_object, nil
}

func (bot *DiscordBot) Send(message string) error {
	if bot.Obj == nil {
		_, err := bot.Get()
		if err != nil {
			return err
		}
	}

	send_obj := DiscordSendObject{bot.Obj.Name, "", message}
	json_obj, err := json.Marshal(send_obj)
	if err != nil {
		return err
	}

	resp, err := http.Post(bot.Url, "application/json", bytes.NewReader(json_obj))
	if err != nil {
		return err
	}

	/*
		buf, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
	*/

	err = resp.Body.Close()
	if err != nil {
		return err
	}

	return nil
}

func DiscordBotInitialize(webhook_url string) (*DiscordBot, error) {
	res := &DiscordBot{webhook_url, nil}
	_, err := res.Get()
	if err != nil {
		return nil, err
	}

	return res, nil
}
