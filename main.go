package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

var N []News
var cfg Config

type discordCodeRequest struct {
	Code string `json:"code"`
}
type UserInfoResponse struct {
	ID                   string      `json:"id"`
	Username             string      `json:"username"`
	Avatar               string      `json:"avatar"`
	Discriminator        string      `json:"discriminator"`
	PublicFlags          json.Number `json:"public_flags"`
	Flags                string      `json:"flags"`
	Banner               string      `json:"banner"`
	AccentColor          string      `json:"accent_color"`
	GlobalName           string      `json:"global_name"`
	AvatarDecorationData string      `json:"avatar_decoration_data"`
	BannerColor          string      `json:"banner_color"`
	Clan                 string      `json:"clan"`
	MFAEnabled           string      `json:"mfa_enabled"`
	Locale               string      `json:"locale"`
	PremiumType          string      `json:"premium_type"`
	Email                string      `json:"email"`
	Verified             string      `json:"verified"`
}
type discordCodeResponse struct {
	TokenType    string `json:"token_type"`
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
}

type Config struct {
	Discord struct {
		Token                 string `json:"token"`
		Secret                string `json:"secret"`
		Id                    string `json:"id"`
		ServerID              string `json:"server_id"`
		MembershipChannelID   string `json:"membership_channel_id"`
		Barer                 string `json:"barer"`
		Safe_channel_id       string `json:"safe_channel_id"`
		Aggressive_channel_id string `json:"aggressive_channel_id"`
	} `json:"discord"`
	Telegram struct {
		Token            string `json:"token"`
		SafeChatID       string `json:"safe_chat_id"`
		AggressiveChatID string `json:"aggressive_chat_id"`
	} `json:"telegram"`
}

type PostRequest struct {
	Message  string `json:"message"`
	Discord  string `json:"discord"`
	Telegram string `json:"telegram"`
	Token    string `json:"token"`
	Type     string `json:"type"`
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
			if req.Type == "aggressive" {
				fmt.Println(cfg.Discord.Aggressive_channel_id)
				Wiscord.ChannelMessageSend(cfg.Discord.Aggressive_channel_id, req.Message)
			}
			if req.Type == "safe" {
				Wiscord.ChannelMessageSend(cfg.Discord.Safe_channel_id, req.Message)

			}
			if req.Type == "both" {
				Wiscord.ChannelMessageSend(cfg.Discord.Safe_channel_id, req.Message)
				Wiscord.ChannelMessageSend(cfg.Discord.Aggressive_channel_id, req.Message)
			}
		}
		if req.Telegram == "true" {
			resp := &http.Response{}
			baseURL := "https://api.telegram.org/bot" + cfg.Telegram.Token + "/sendMessage"
			if req.Type == "safe" {
				params := url.Values{}
				params.Add("chat_id", cfg.Telegram.SafeChatID)
				params.Add("text", req.Message)

				resp, err = http.Get(baseURL + "?" + params.Encode())
			} else if req.Type == "aggressive" {
				params := url.Values{}
				params.Add("chat_id", cfg.Telegram.AggressiveChatID)
				params.Add("text", req.Message)

				resp, err = http.Get(baseURL + "?" + params.Encode())
			} else {
				params := url.Values{}
				params.Add("chat_id", cfg.Telegram.SafeChatID)
				params.Add("text", req.Message)

				_, _ = http.Get(baseURL + "?" + params.Encode())

				params = url.Values{}
				params.Add("chat_id", cfg.Telegram.AggressiveChatID)
				params.Add("text", req.Message)

				resp, err = http.Get(baseURL + "?" + params.Encode())
			}
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
	router.POST("/discordtoken", func(c *gin.Context) {

		var temp discordCodeRequest
		if err := c.ShouldBindJSON(&temp); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		var code = temp.Code
		url := "https://discord.com/api/v10/oauth2/token"
		method := "POST"

		payload := strings.NewReader("grant_type=authorization_code&code=" + code + "&redirect_uri=https%3A%2F%2Fapi.rollintrades.com%3A3000%2Fbuy")

		client := &http.Client{}
		req, err := http.NewRequest(method, url, payload)

		if err != nil {
			fmt.Println(err)
			return
		}
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Add("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(cfg.Discord.Id+":"+cfg.Discord.Secret))+"==")
		req.Header.Add("Cookie", "__cfruid=ca60a239a74cf77d972dd3f37de310957ad863c3-1719434344; __dcfduid=2086805833fc11ef8ad7ee895a1966da; __sdcfduid=2086805833fc11ef8ad7ee895a1966dabb64ffed8fe077cff16d81dd824bd75afa6b7ad6948c60825098f4f60d57cc62; _cfuvid=dTeEv5tUO18W2t4hZmhxCNcCj7KjGJ2RlCzC7bwbVtY-1719434344243-0.0.1.1-604800000")

		res, err := client.Do(req)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer res.Body.Close()

		var result discordCodeResponse
		body, _ := ioutil.ReadAll(res.Body)
		err = json.Unmarshal(body, &result)
		if err != nil {
			fmt.Println(err)
			return
		}

		url = "https://discord.com/api/users/@me"

		req, err = http.NewRequest("GET", url, nil)
		if err != nil {
			fmt.Println("error creating request: ", err)
		}

		req.Header.Set("Authorization", "Bearer "+result.AccessToken)

		client = &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("error making request: ", err)
		}
		defer resp.Body.Close()

		var userInfo UserInfoResponse
		body2, _ := ioutil.ReadAll(resp.Body)
		err = json.Unmarshal(body2, &userInfo)
		if err != nil {
			fmt.Println("error unmarshalling response: ", err)
		}

		c.JSON(http.StatusOK, gin.H{"data": userInfo})
	})
	go router.RunTLS(":1809", "certificate.crt", "private.key")
	// go router.Run("localhost:1809")
	Run()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}
