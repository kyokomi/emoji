package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

const emojoV5DBJsonURL = "https://raw.githubusercontent.com/CodeFreezr/emojo/master/db/v5/emoji-v5.json"

// Emojo json parse struct
type Emojo struct {
	No          int    `json:"No"`
	Emoji       string `json:"Emoji"`
	Category    string `json:"Category"`
	SubCategory string `json:"SubCategory"`
	Unicode     string `json:"Unicode"`
	Name        string `json:"Name"`
	Tags        string `json:"Tags"`
	Shortcode   string `json:"Shortcode"`
}

func createEmojoCodeMap() (map[string]string, error) {
	res, err := http.Get(emojoV5DBJsonURL)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	emojiFile, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var gs []Emojo
	if err := json.Unmarshal(emojiFile, &gs); err != nil {
		return nil, err
	}

	emojiCodeMap := make(map[string]string)
	for _, gemoji := range gs {
		shortCode := strings.Replace(gemoji.Shortcode, ":", "", 2)
		if len(shortCode) == 0 || len(gemoji.Emoji) == 0 {
			continue
		}
		code := gemoji.Emoji
		emojiCodeMap[shortCode] = fmt.Sprintf("%+q", strings.ToLower(code))
	}

	return emojiCodeMap, nil
}
