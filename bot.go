package main

import (
	"crypto/tls"
	"fmt"
	"github.com/thoj/go-ircevent"
)

const server = "irc.freenode.net:7000"

func main() {
	ircnick := "dkanjus0"
	irccon := irc.IRC(ircnick, "IRCTestSSL")
	irccon.VerboseCallbackHandler = true
	irccon.Debug = true
	irccon.UseTLS = true
	irccon.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	irccon.AddCallback("001", func(e *irc.Event) { irccon.Join("#go-eventirc-test") })
	irccon.AddCallback("366", func(e *irc.Event) {})
	err := irccon.Connect(server)
	if err != nil {
		fmt.Printf("Err %s", err)
		return
	}
	irccon.Loop()
}
