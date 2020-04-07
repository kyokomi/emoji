package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const unicodeorgURL = "http://www.unicode.org/emoji/charts/emoji-list.html"

func createUnicodeorgMap() (map[string]string, error) {
	res, err := http.Get(unicodeorgURL)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
	}

	return generateUnicodeorgCodeMap(res.Body)
}

// UnicodeorgEmoji unicode.org emoji
type UnicodeorgEmoji struct {
	No            int
	Code          string
	ShortName     string
	OtherKeywords []string
}

var shortNameReplaces = []string{
	":", "",
	",", "",
	"⊛", "", // \U+229B
	"“", "", // \U+201C
	"”", "", // \U+201D
}

func generateUnicodeorgCodeMap(body io.ReadCloser) (map[string]string, error) {
	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		return nil, err
	}

	var emojis []*UnicodeorgEmoji
	doc.Find("table").First().Find("tr").Each(func(i int, selection *goquery.Selection) {
		var cols []string
		selection.Find("td").Each(func(j int, s *goquery.Selection) {
			cols = append(cols, s.Text())
		})

		if len(cols) != 5 {
			return
		}

		unicodeEmoji := UnicodeorgEmoji{}
		unicodeEmoji.No, err = strconv.Atoi(cols[0])
		if err != nil {
			log.Println("ERROR: no", err)
			return
		}
		codes := strings.Fields(cols[1])
		var sb strings.Builder
		for _, code := range codes {
			code = strings.ReplaceAll(code, "U+", "")
			s, err := strconv.ParseInt(code, 16, 32)
			if err != nil {
				log.Println("ERROR: code", err)
				return
			}
			sb.WriteRune(rune(s))
		}
		unicodeEmoji.Code = sb.String()
		shortName := strings.NewReplacer(shortNameReplaces...).Replace(cols[3])
		unicodeEmoji.ShortName = strings.Replace(strings.TrimSpace(shortName), " ", "_", -1)
		unicodeEmoji.OtherKeywords = strings.Fields(cols[4])
		emojis = append(emojis, &unicodeEmoji)
	})

	emojiCodeMap := make(map[string]string)
	for _, emoji := range emojis {
		emojiCodeMap[emoji.ShortName] = fmt.Sprintf("%+q", emoji.Code)
	}

	return emojiCodeMap, nil
}
