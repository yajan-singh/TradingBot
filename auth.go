package main

import (
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"strings"
)

func Get_token() string {
	url := "https://charts.trendspider.com/authentication/1/member/login/"
	method := "POST"

	payload := strings.NewReader("email=takshpatel8844%40gmail.com&password=Eliteoptions.com&rememberme=forever&redirect=%2F")
	jar, _ := cookiejar.New(nil)
	client := &http.Client{Jar: jar}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return "ERROR"
	}
	req.Header.Add("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Add("accept-language", "en-US,en;q=0.9,fr;q=0.8")
	req.Header.Add("cache-control", "max-age=0")
	req.Header.Add("content-type", "application/x-www-form-urlencoded")
	req.Header.Add("cookie", "_vwo_uuid_v2=DD95199657BD234FA7F7B7AFE4F7F847C|7cd27df1eb26396719da0853da04e144; _vwo_uuid=DD95199657BD234FA7F7B7AFE4F7F847C; _vwo_ds=3%241717533515%3A4.08910658%3A%3A; _vwo_sn=0%3A1%3A%3A%3A1; _vis_opt_s=1%7C; _vis_opt_test_cookie=1; _vis_opt_exp_68_combi=3; _gcl_aw=GCL.1717533518.EAIaIQobChMIx_Xu8ubChgMVo07_AR3rfwcoEAAYASAAEgLFZfD_BwE; _gcl_gs=2.1.k1$i1717533512; _gcl_au=1.1.25559245.1717533518; _uetsid=6bd20d8022b211efa8831b686e8f0c35; _uetvid=6bd22ca022b211efa68beb53ee9a3943; _ga=GA1.1.1301801772.1717533519; _fprom_code=_r_viktoriia80; _fprom_track=b9f3396a-81de-47a7-b1c2-2d76d9b4d36f; _ga_S5PJ8PHJV4=GS1.1.1717533519.1.0.1717533523.56.0.0; _clck=gaqvo%7C2%7Cfmc%7C0%7C1616")
	req.Header.Add("dnt", "1")
	req.Header.Add("origin", "https://charts.trendspider.com")
	req.Header.Add("priority", "u=0, i")
	req.Header.Add("referer", "https://charts.trendspider.com/login/?reason=not_signed_in")
	req.Header.Add("sec-ch-ua", "\"Google Chrome\";v=\"125\", \"Chromium\";v=\"125\", \"Not.A/Brand\";v=\"24\"")
	req.Header.Add("sec-ch-ua-mobile", "?0")
	req.Header.Add("sec-ch-ua-platform", "\"macOS\"")
	req.Header.Add("sec-fetch-dest", "document")
	req.Header.Add("sec-fetch-mode", "navigate")
	req.Header.Add("sec-fetch-site", "same-origin")
	req.Header.Add("sec-fetch-user", "?1")
	req.Header.Add("upgrade-insecure-requests", "1")
	req.Header.Add("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0.0.0 Safari/537.36")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return "ERROR"
	}
	defer res.Body.Close()
	token := ""
	fmt.Println(res.StatusCode)
	fmt.Println(res)
	for i := range jar.Cookies(req.URL) {
		if jar.Cookies(req.URL)[i].Name == "auth_token" {
			token = jar.Cookies(req.URL)[i].Value
		}
	}
	if token == "" {
		fmt.Println("Token not found")
		return "ERROR"
	}
	return token
}
