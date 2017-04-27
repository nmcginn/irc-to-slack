package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/thoj/go-ircevent"
)

func connect_irc() (irccon *irc.Connection, err error) {
	const server = "irc.freenode.net:7000"
	const ircnick = "dkanjus1"

	irccon = irc.IRC(ircnick, "IRCTestSSL")
	irccon.VerboseCallbackHandler = false
	irccon.Debug = false
	irccon.UseTLS = true
	irccon.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	irccon.AddCallback("PRIVMSG", func(e *irc.Event) {
		send_to_slack(e.Nick, e.Message())
	})
	err = irccon.Connect(server)

	return irccon, err
}

func send_to_slack(from string, text string) (err error) {
	slack_payload := slack_message{
		Text: text,
		Channel: "",
		Username: from,
		Icon: ":nyx:",
	}
	payload, err := json.Marshal(slack_payload)
	if err != nil {
		return err
	}

	fmt.Println(string(payload))
	return nil
}

type slack_message struct {
	Text string `json:"text"`
	Channel string `json:"channel"`
	Username string `json:"username"`
	Icon string `json:"icon_emoji"`
}

func main() {
	irccon, err := connect_irc()
	if err != nil {
		fmt.Printf("Err %s", err)
		return
	}

	irccon.Loop()
}
