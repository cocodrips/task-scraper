package scraper

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"regexp"
	"strings"
)

var (
	max_depth int = 20
)

type Message struct {
	Ok       bool `json:"ok"`
	Messages []struct {
		Text      string `json:"text"`
		User      string `json:"user"`
		TimeStamp string `json:"ts"`
	} `json:"messages"`
	HasMore bool `json:"has_more"`
}

var (
	endpoint = "https://slack.com/api/"
)

func channelHistoryUrl(timestamp string) (string, error) {
	u, err := url.Parse(endpoint)
	if err != nil {
		return "", err
	}

	u.Path = path.Join(u.Path, "channels.history")

	q := u.Query()
	q.Set("token", os.Getenv("OAUTH_TOKEN"))
	q.Set("channel", os.Getenv("CHANNEL"))
	q.Set("count", "1000")
	if timestamp != "" {
		q.Set("latest", timestamp)
	}

	u.RawQuery = q.Encode()
	return u.String(), nil
}

func EnglishCondition(message string) bool {
	return strings.HasPrefix(message, "今日の英語")
}

func EnglishOutput(message string) {
	for _, line := range strings.Split(message, "\n") {
		if !strings.Contains(line, "今日の英語") {
			fmt.Println(strings.Replace(line, "- ", "", -1))
		}
	}
}

func TaskCondition(message string) bool {
	re := regexp.MustCompile(`(\d{4}/\d{2}/\d{2}).?tasks.+`)
	return re.MatchString(message)
}

func TaskOutput(message string) {
	for _, line := range strings.Split(message, "\n") {
		fmt.Println(line)
	}
	fmt.Println("-----------")
}



func GetMessage(condition func(string) bool, output func(string)) {

	var responseJson Message
	timestamp := ""

	for i := 0; i < max_depth; i++ {
		// Create request
		url, _ := channelHistoryUrl(timestamp)
		req, _ := http.NewRequest("GET", url, nil)
		res, _ := http.DefaultClient.Do(req)

		defer res.Body.Close()
		body, _ := ioutil.ReadAll(res.Body)

		// Parse texta
		err := json.Unmarshal(body, &responseJson)
		if err != nil {
			fmt.Printf("Failed to parse\n")
		}

		for _, message := range responseJson.Messages {
			if message.User != os.Getenv("USER_ID") {
				continue
			}
			if !condition(message.Text) {
				continue
			}

			output(message.Text)
			timestamp = message.TimeStamp
		}
		if !responseJson.HasMore {
			break
		}
	}

}
