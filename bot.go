package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/thoj/go-ircevent"
	"io/ioutil"
	"net/http"
	"os"
)

func connect_irc() (irccon *irc.Connection, err error) {
	server := os.Getenv("IRC_SERVER")
	ircnick := os.Getenv("IRC_NICK")

	irccon = irc.IRC(ircnick, ircnick)
	irccon.Password = os.Getenv("IRC_PASSWORD")
	irccon.VerboseCallbackHandler = false
	irccon.Debug = false
	irccon.UseTLS = true
	irccon.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	irccon.AddCallback("PRIVMSG", func(e *irc.Event) {
		if e.Nick == "FleetBot" {
			err = send_to_slack(e.Nick, e.Message())
			if err != nil {
				fmt.Printf("Err %s", err)
			}
		} else {
			fmt.Printf("%s: %s\n", e.Nick, e.Message())
		}
	})

	irccon.AddCallback("PING", func(e *irc.Event) {
		err = send_to_slack(e.Nick, e.Message())
		if err != nil {
			fmt.Printf("Err %s", err)
		}
	})

	err = irccon.Connect(server)

	return irccon, err
}

func send_to_slack(from string, text string) (err error) {
	slack_payload := slack_message{
		Text:     text,
		Username: from,
	}
	payload, err := json.Marshal(slack_payload)
	if err != nil {
		return err
	}

	resp, err := http.Post(os.Getenv("WEBHOOK_URL"), "application/json", bytes.NewBuffer(payload))
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		body, _ := ioutil.ReadAll(resp.Body)
		err = errors.New(string(body))
	}

	return err
}

type slack_message struct {
	Text     string `json:"text"`
	Channel  string `json:"channel"`
	Username string `json:"username"`
	Icon     string `json:"icon_emoji"`
}

func main() {
	irccon, err := connect_irc()
	if err != nil {
		fmt.Printf("Err %s", err)
		return
	}

	irccon.Loop()
}
