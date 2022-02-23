package crawler

import (
	"SecCrawler/register"
	"SecCrawler/utils"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"time"
)

type DongJian struct{}

func (crawler DongJian) Config() register.CrawlerConfig {
	return register.CrawlerConfig{
		Name:        "DongJian",
		Description: "洞见微信聚合",
	}
}

// Get 获取洞见微信聚合前24小时内文章。
func (crawler DongJian) Get() ([][]string, error) {
	client := &http.Client{
		Timeout: time.Duration(4) * time.Second,
	}
	req, err := http.NewRequest("GET", "http://wechat.doonsec.com/rss.xml", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.55 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Set("Sec-Fetch-Site", "none")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	bodyString := string(body)

	re := regexp.MustCompile(`<item><title>([\w\W]*?)</title><link>([\w\W]*?)</link><description>[\w\W]*?</description><author>[\w\W]*?</author><category>[\w\W]*?</category><pubDate>([\w\W]*?)</pubDate></item>`)
	result := re.FindAllStringSubmatch(strings.TrimSpace(bodyString), -1)

	var resultSlice [][]string
	fmt.Printf("[*] [DongJian] crawler result:\n%s\n\n", utils.CurrentTime())
	for _, match := range result {
		time_zone := time.FixedZone("CST", 8*3600)
		t, err := time.ParseInLocation(time.RFC1123Z, match[1:][2], time_zone)
		if err != nil {
			return nil, err
		}

		if !utils.IsIn24Hours(t.In(time_zone)) {
			// 默认时间顺序是从近到远
			break
		}

		// 去除title中的换行符
		re, _ = regexp.Compile(`\s{1,}`)
		match[1:][0] = re.ReplaceAllString(match[1:][0], "")

		fmt.Println(t.In(time_zone).Format("2006/01/02 15:04:05"))
		fmt.Println(match[1:][0])
		fmt.Printf("%s\n\n", match[1:][1])

		resultSlice = append(resultSlice, match[1:][0:2])
	}
	// slice中title和url调换位置，以符合统一的format
	for _, item := range resultSlice {
		item[0], item[1] = item[1], item[0]
	}
	if len(resultSlice) == 0 {
		return nil, errors.New("no records in the last 24 hours")
	}
	return resultSlice, nil

}
