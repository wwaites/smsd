package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

func AndrewsArnoldMt(cfg Config, rt Route, msg Message) (err error) {
	params := url.Values{}
	if rt.User != "" {
		params.Add("username", rt.User)
	}
	if rt.Pass != "" {
		params.Add("password", rt.Pass)
	}
	params.Add("oa", msg.Src)
	params.Add("ud", msg.Msg)

	resp, err := http.PostForm("https://sms.aa.net.uk/sms.cgi", params)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	log.Printf("MT: sending message to andrews and arnold: %s", body)
	return
}

func AndrewsArnoldMo(cfg Config, rt Route, msg Message) (err error) {
	params := url.Values{}
	if rt.User != "" {
		params.Add("username", rt.User)
	}
	if rt.Pass != "" {
		params.Add("password", rt.Pass)
	}
	if rt.Src != "" {
		params.Add("oa", rt.Src)
	}
	params.Add("da", msg.Dst)
	params.Add("ud", msg.Msg)

	resp, err := http.PostForm("https://sms.aa.net.uk/sms.cgi", params)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	log.Printf("MO sending message to andrews and arnold: %s", body)
	return
}
