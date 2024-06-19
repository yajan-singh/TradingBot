package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

var N []News
var cfg Config

type Config struct {
	Discord struct {
		Token                  string `json:"token"`
		ServerID               string `json:"server_id"`
		MembershipChannelID    string `json:"membership_channel_id"`
		Barer                  string `json:"barer"`
		anouncement_channel_id string `json:"anouncement_channel_id"`
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
	Token    string `json:"token"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

var Wiscord *discordgo.Session

func validateToken(tokenString string) bool {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte("secret"), nil
	})
	if err != nil {
		fmt.Println(err)
		return false
	}

	if _, ok := token.Claims.(jwt.MapClaims); ok {
		return true
	} else {
		return false
	}
}

func main() {
	file, _ := os.ReadFile("config.json")
	_ = json.Unmarshal([]byte(file), &cfg)
	N = Req_news(cfg.Discord.Barer)
	for {
		if N != nil {
			break
		}
		token := Get_token()
		fmt.Println("TokenMain: ", token)
		if token != "ERROR" {
			cfg.Discord.Barer = token
		} else {
			fmt.Println("Cannot fetch token")
			return
		}
		N = Req_news(cfg.Discord.Barer)
	}
	var BotToken = cfg.Discord.Token
	Wiscord, err = discordgo.New("Bot " + BotToken)
	checkNilErr(err)
	Wiscord.Open()
	defer Wiscord.Close()

	// RESPFullAPI
	router := gin.Default()
	config := cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"POST", "GET"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return origin == "https://github.com"
		},
		MaxAge: 12 * time.Hour,
	}

	router.Use(cors.New(config))
	router.POST("/login", func(c *gin.Context) {
		fmt.Println("Login")
		var req LoginRequest

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if req.Email == "admin@eliteoptions.com" && req.Password == "Eliteoptions.com" {
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"email": req.Email,
				"exp":   time.Now().Add(time.Hour * 24).Unix(),
			})
			tokenString, err := token.SignedString([]byte("secret"))
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"status": "INTERNAL_SERVER_ERROR"})
				return
			}
			c.JSON(http.StatusOK, gin.H{"status": "LOGGED_IN", "token": tokenString})
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"status": "UNAUTHORIZED"})
		}
	})
	router.POST("/announcement", func(c *gin.Context) {

		var req PostRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		fmt.Println(req.Token)
		if !validateToken(req.Token) {
			c.JSON(http.StatusUnauthorized, gin.H{"status": "UNAUTHORIZED"})
			return
		}
		if req.Discord == "true" {

			Wiscord.ChannelMessageSend(cfg.Discord.anouncement_channel_id, req.Message)
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
