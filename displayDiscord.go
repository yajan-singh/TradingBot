package main

import (
	"fmt"
	"log"

	// "os"
	// "os/signal"

	"time"

	"github.com/bwmarrin/discordgo"
	"golang.org/x/exp/slices"
)

var Dupes []string

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
	for i := range N {
		if !slices.Contains(Dupes, N[i].URL) {
			embed := &discordgo.MessageEmbed{
				URL:         N[i].URL,
				Title:       N[i].Title,
				Description: N[i].Teaser,
				Color:       0x00ff00,
				Author: &discordgo.MessageEmbedAuthor{
					Name: N[i].Ticker,
				},
			}

			Wiscord.ChannelMessageSendEmbed(cfg.Discord.MembershipChannelID, embed)
			if len(Dupes) >= 100 {
				Dupes = Dupes[1:]
				Dupes = append(Dupes, N[i].URL)
			} else {
				Dupes = append(Dupes, N[i].URL)
			}
		}
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
				fmt.Println("TokenWatch: ", token)
				if token != "ERROR" {
					cfg.Discord.Barer = token
				} else {
					fmt.Println("Cannot fetch token")
					return
				}
				N = Req_news(cfg.Discord.Barer)
			}
			for i := range N {
				if !slices.Contains(Dupes, N[i].URL) {
					embed := &discordgo.MessageEmbed{
						URL:         N[i].URL,
						Title:       N[i].Title,
						Description: N[i].Teaser,
						Color:       0x00ff00,
						Author: &discordgo.MessageEmbedAuthor{
							Name: N[i].Ticker,
						},
						Thumbnail: &discordgo.MessageEmbedThumbnail{
							URL: "https://cdn.iconscout.com/icon/premium/png-256-thumb/url-2879059-2393887.png",
						},
					}

					Wiscord.ChannelMessageSendEmbed(cfg.Discord.MembershipChannelID, embed)
					if len(Dupes) >= 100 {
						Dupes = Dupes[1:]
						Dupes = append(Dupes, N[i].URL)
					} else {
						Dupes = append(Dupes, N[i].URL)
					}
				}
			}
		}
	}
}
