package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
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
	Telegram struct {
		Token  string `json:"token"`
		ChatID string `json:"chat_id"`
	} `json:"telegram"`
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

		if req.Discord == "true" {
			Wiscord.ChannelMessageSend("1250111990395965550", req.Message)
		}
		if req.Telegram == "true" {
			baseURL := "https://api.telegram.org/bot" + cfg.Telegram.Token + "/sendMessage"
			params := url.Values{}
			params.Add("chat_id", cfg.Telegram.ChatID)
			params.Add("text", req.Message)

			resp, err := http.Get(baseURL + "?" + params.Encode())
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				body, _ := ioutil.ReadAll(resp.Body)
				fmt.Println("Error:", resp.StatusCode, string(body))
			} else {
				fmt.Println("Message sent successfully!")
			}
		}
		c.JSON(http.StatusOK, gin.H{"status": "SENT"})
	})
	go router.Run("localhost:1809")

	Run()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}
