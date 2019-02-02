package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

const gemojiDBJsonURL = "https://raw.githubusercontent.com/github/gemoji/master/db/emoji.json"

// GemojiEmoji gemoji json parse struct
type GemojiEmoji struct {
	Aliases     []string `json:"aliases"`
	Description string   `json:"description"`
	Emoji       string   `json:"emoji"`
	Tags        []string `json:"tags"`
}

func createGemojiCodeMap() (map[string]string, error) {
	res, err := http.Get(gemojiDBJsonURL)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	emojiFile, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var gs []GemojiEmoji
	if err := json.Unmarshal(emojiFile, &gs); err != nil {
		return nil, err
	}

	emojiCodeMap := make(map[string]string)
	for _, gemoji := range gs {
		for _, a := range gemoji.Aliases {
			if len(a) == 0 || len(gemoji.Emoji) == 0 {
				continue
			}
			code := gemoji.Emoji
			emojiCodeMap[a] = fmt.Sprintf("%+q", strings.ToLower(code))
		}
	}

	return emojiCodeMap, nil
}
