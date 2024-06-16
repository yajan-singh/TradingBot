package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"time"

	"golang.org/x/exp/slices"
)

type News struct {
	Timestamp int    `json:"timestamp"`
	DaysAgo   string `json:"daysAgo"`
	Author    string `json:"author"`
	Title     string `json:"title"`
	Teaser    string `json:"teaser"`
	URL       string `json:"url"`
	ID        string `json:"id"`
	Ticker    string `json:"ticker"`
}

func Req_news(token string) []News {
	url := "https://charts.trendspider.com/authentication/1/api?key=&path=%2Fnon_market_data%2F1%2Ffeeds%2Fnews__flow%3Fquery%3Dbase64%3AeyJmb3JtYXQiOiJkaXNwbGF5X2Zsb3ciLCJzb3J0QnkiOnsiZmllbGQiOiJ0aW1lc3RhbXAiLCJkaXJlY3Rpb24iOiJkZXNjIn19"
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)

	}
	token = "Bearer " + token
	req.Header.Add("accept", "*/*")
	req.Header.Add("accept-language", "en-US,en;q=0.9,fr;q=0.8")
	req.Header.Add("authorization", token)
	req.Header.Add("dnt", "1")
	req.Header.Add("priority", "u=1, i")
	req.Header.Add("referer", "https://charts.trendspider.com/dashboard")
	req.Header.Add("sec-ch-ua", "\"Google Chrome\";v=\"125\", \"Chromium\";v=\"125\", \"Not.A/Brand\";v=\"24\"")
	req.Header.Add("sec-ch-ua-mobile", "?0")
	req.Header.Add("sec-ch-ua-platform", "\"macOS\"")
	req.Header.Add("sec-fetch-dest", "empty")
	req.Header.Add("sec-fetch-mode", "cors")
	req.Header.Add("sec-fetch-site", "same-origin")
	req.Header.Add("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0.0.0 Safari/537.36")
	req.Header.Add("x-client", "dashboard")
	req.Header.Add("x-requested-with", "XMLHttpRequest")
	req.Header.Add("x-requester", "data-flow")
	req.Header.Add("x-workspace-id", "40abd119448f56e7ae7b529ba8696d76")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
	}
	var news []News
	err = json.Unmarshal(body, &news)
	if err != nil {
		fmt.Println(err)
	}
	var tickers []string

	for i := range news {
		if (news[i].Timestamp < (int(time.Now().Unix()) - 660)) && slices.Contains(tickers, news[i].Ticker) {
			news = append(news[:i], news[i+1:]...)
		}
	}
	sort.Slice(news[:], func(i, j int) bool {
		return news[i].Timestamp < news[j].Timestamp
	})
	return news
}
