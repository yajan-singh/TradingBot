package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	// "os"
	// "os/signal"
	"strings"
	"time"
)

func checkNilErr(e error) {
	if e != nil {
		log.Fatal("Error message")
	}
}

var err error

func Run() {
	// var BotToken = cfg.Discord.Token
	// // create a session
	// Wiscord, err = discordgo.New("Bot " + BotToken)
	// checkNilErr(err)

	// add a event handler
	msg, err := Wiscord.ChannelMessages(cfg.Discord.MembershipChannelID, 100, "", "", "")
	if err == nil && len(msg) > 0 {
		fmt.Println("Message Found!")
	}
	flag := true
	for i := range N {
		for j := range msg {
			if strings.Contains(msg[j].Content, N[i].URL) {
				fmt.Println("DUPE FOUND!")
				flag = false
				break
			}
		}
		if flag {
			Wiscord.ChannelMessageSend(cfg.Discord.MembershipChannelID, N[i].URL)
			msg, err = Wiscord.ChannelMessages(cfg.Discord.MembershipChannelID, 100, "", "", "")
			if err == nil && len(msg) > 0 {
				fmt.Println("Message Found!")
			}
		}
		flag = true

	}

	// open session
	// Wiscord.Open()
	go watch()
	// defer Wiscord.Close() // close session, after function termination

	// keep bot running untill there is NO os interruption (ctrl + C)
	fmt.Println("Bot running....")
	// c := make(chan os.Signal, 1)
	// signal.Notify(c, os.Interrupt)
	// <-c

}

func watch() {
	for {
		if time.Now().Second()%5 == 0 {
			N = Req_news(cfg.Discord.Barer)
			for {
				if N != nil {
					break
				}
				token := Get_token()
				fmt.Println("Token: ", token)
				if token != "ERROR" {
					cfg.Discord.Barer = token
				} else {
					fmt.Println("Cannot fetch token")
					return
				}
				r, _ := json.Marshal(cfg)
				err = ioutil.WriteFile("output.json", r, 0644)
				N = Req_news(cfg.Discord.Barer)
			}
			msg, err := Wiscord.ChannelMessages(cfg.Discord.MembershipChannelID, 100, "", "", "")
			if err == nil && len(msg) > 0 {
				fmt.Println("Message: ", msg[0].Content)
			}
			flag := true
			for i := range N {
				for j := range msg {
					if strings.Contains(msg[j].Content, N[i].URL) {
						fmt.Println("DUPE FOUND!")
						flag = false
						break
					}
				}
				if flag {
					Wiscord.ChannelMessageSend(cfg.Discord.MembershipChannelID, N[i].URL)
					msg, err = Wiscord.ChannelMessages(cfg.Discord.MembershipChannelID, 100, "", "", "")
					if err == nil && len(msg) > 0 {
						fmt.Println("Message: ", msg[0].Content)
					}
				}
				flag = true
			}
		}
	}
}
