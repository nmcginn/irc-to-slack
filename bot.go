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
	"strings"
)

var prev_msg string

func connect_irc() (irccon *irc.Connection, err error) {
	server := os.Getenv("IRC_SERVER")
	ircnick := os.Getenv("IRC_NICK")

	irccon = irc.IRC(ircnick, ircnick)
	irccon.Password = os.Getenv("IRC_PASSWORD")
	irccon.VerboseCallbackHandler = false
	irccon.Debug = false
	irccon.UseTLS = true
	irccon.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	irccon.AddCallback("001", func(e *irc.Event) {})

	irccon.AddCallback("PRIVMSG", func(e *irc.Event) {
		if e.Nick == "sphere" {
			err = send_to_slack(e.Nick, e.Message())
			if err != nil {
				fmt.Printf("Err %s", err)
			}
		} else {
			fmt.Printf("%s: %s\n", e.Nick, e.Message())
		}
	})

	irccon.AddCallback("NOTICE", func(e *irc.Event) {
		fmt.Printf("NOTICE %s: %s\n", e.Nick, e.Message())
		if err != nil {
			fmt.Printf("Err %s", err)
		}
	})

	err = irccon.Connect(server)

	return irccon, err
}

func send_to_slack(from string, text string) (err error) {
	if prev_msg == text {
		return nil
	} else {
		prev_msg = text
	}
	slack_payload := slack_message{
		Text:     "@here " + text,
		Username: from,
	}
	payload, err := json.Marshal(slack_payload)
	if err != nil {
		return err
	}

	var resp *http.Response
	if strings.Contains(text, "supers") || strings.Contains(text, "supercaps") || strings.Contains(text, "titans") {
		resp, err = http.Post(os.Getenv("WEBHOOK_URL_S"), "application/json", bytes.NewBuffer(payload))
	} else {
		resp, err = http.Post(os.Getenv("WEBHOOK_URL"), "application/json", bytes.NewBuffer(payload))
	}
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
	prev_msg = ""
	irccon, err := connect_irc()
	if err != nil {
		fmt.Printf("Err %s", err)
		return
	}

	irccon.Loop()
}
