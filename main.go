package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
	"github.com/gin-gonic/gin"
)

var N []News
var cfg Config

type Config struct {
	Discord struct {
		Token               string `json:"token"`
		ServerID            string `json:"server_id"`
		MembershipChannelID string `json:"membership_channel_id"`
		Barer               string `json:"barer"`
	} `json:"discord"`
}

type PostRequest struct {
	Message  string `json:"message"`
	Discord  string `json:"discord"`
	Telegram string `json:"telegram"`
}

var Wiscord *discordgo.Session

func main() {
	file, _ := os.ReadFile("config.json")
	_ = json.Unmarshal([]byte(file), &cfg)
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
	var BotToken = cfg.Discord.Token
	Wiscord, err = discordgo.New("Bot " + BotToken)
	checkNilErr(err)
	Wiscord.Open()
	defer Wiscord.Close()

	// RESPFullAPI
	router := gin.Default()
	router.POST("/announcement", func(c *gin.Context) {
		var req PostRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		fmt.Println("Message: ")
		fmt.Println(req.Message)

		if req.Discord == "true" {
			Wiscord.ChannelMessageSend("1250111990395965550", req.Message)
		}
		if req.Telegram == "true" {
			fmt.Println("TELEGRAM")
		}
		c.JSON(http.StatusOK, gin.H{"status": "SENT"})
	})
	go router.Run(":1809")

	Run()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}
