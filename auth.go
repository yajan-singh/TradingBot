package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"strings"
)

func logout() {

	url := "https://charts.trendspider.com/logout/"
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Add("accept-language", "en-US,en;q=0.9,fr;q=0.8")
	req.Header.Add("cookie", "_vwo_uuid_v2=DD95199657BD234FA7F7B7AFE4F7F847C|7cd27df1eb26396719da0853da04e144; _vwo_uuid=DD95199657BD234FA7F7B7AFE4F7F847C; _vwo_ds=3%241717533515%3A4.08910658%3A%3A; _vis_opt_s=1%7C; _vis_opt_test_cookie=1; _vis_opt_exp_68_combi=3; _gcl_aw=GCL.1717533518.EAIaIQobChMIx_Xu8ubChgMVo07_AR3rfwcoEAAYASAAEgLFZfD_BwE; _gcl_gs=2.1.k1$i1717533512; _gcl_au=1.1.25559245.1717533518; _uetvid=6bd22ca022b211efa68beb53ee9a3943; _fprom_code=_r_viktoriia80; _fprom_track=b9f3396a-81de-47a7-b1c2-2d76d9b4d36f; _ga_S5PJ8PHJV4=GS1.1.1717533519.1.0.1717533523.56.0.0; _clck=gaqvo%7C2%7Cfmc%7C0%7C1616; _ga=GA1.2.1301801772.1717533519; _gid=GA1.2.233047549.1718659918; _dca32=http://10.0.4.47:80; _3efe6=http://10.0.4.119:80; _36d96=http://10.0.4.172:80; __stripe_mid=9623245f-2cc4-4cf2-9bc6-4f33774740e6e44e5a; __stripe_sid=f79797b3-e908-46db-9458-477e5863f741941d3c; _gat=1; auth_token="+cfg.Discord.Barer+"; wpSGCacheBypass=1; trendspider_user=7b226e616d65223a2254616b736820506174656c222c22706c616e4964223a226f32345f7374616e646172645f6d6f222c22737461747573223a2270616964227d; mp_b107c9df7e4eceb1bfd4253740aa58b7_mixpanel=%7B%22distinct_id%22%3A%20%2283967%22%2C%22%24device_id%22%3A%20%2218fe4f9dbae129b-0d05c853b92818-1a525637-16a7f0-18fe4f9dbaf2981%22%2C%22%24initial_referrer%22%3A%20%22https%3A%2F%2Fcharts.trendspider.com%2Flogin%2F%3Freason%3Dnot_signed_in%22%2C%22%24initial_referring_domain%22%3A%20%22charts.trendspider.com%22%2C%22%24user_id%22%3A%20%2283967%22%7D; _ga_6KZXQ1M7SM=GS1.2.1718659918.5.1.1718660012.0.0.0; trendspider_user=; wpSGCacheBypass=")
	req.Header.Add("dnt", "1")
	req.Header.Add("priority", "u=0, i")
	req.Header.Add("referer", "https://charts.trendspider.com/")
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
		return
	}
	fmt.Println("logout:" + cfg.Discord.Barer)
	defer res.Body.Close()
}

func Get_token() string {
	logout()
	url := "https://charts.trendspider.com/authentication/1/member/login/"
	method := "POST"

	payload := strings.NewReader("email=takshpatel8844%40gmail.com&password=Printmoney@365&rememberme=forever&redirect=%2F")
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
	for i := range jar.Cookies(req.URL) {
		if jar.Cookies(req.URL)[i].Name == "auth_token" {
			token = jar.Cookies(req.URL)[i].Value
		}
	}
	if token == "" {
		fmt.Println("Token not found")
		return "ERROR"
	}
	if token != "ERROR" {
		cfg.Discord.Barer = token
		// update config file
		r, _ := json.Marshal(cfg)
		err = ioutil.WriteFile("config.json", r, 0644)
		if err != nil {
			fmt.Println(err)
		}

	}
	return token
}
