package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

func Pushover(cfg Config, rt Route, msg Message) (err error) {
	params := url.Values{}
	params.Add("token", cfg.Server.Pushover)
	params.Add("user", rt.User)
	params.Add("title", msg.Src + " → " + msg.Dst)
	params.Add("message", msg.Msg)

	resp, err := http.PostForm("https://api.pushover.net/1/messages.json", params)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	log.Printf("MT: sending message to pushover: %s", body)
	return
}
