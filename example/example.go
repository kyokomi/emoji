package main

import (
	"flag"

	"github.com/kyokomi/emoji"
)

func main() {
	emojiKeyword := flag.String("e", ":beer: Beer!!!", "emoji name")
	flag.Parse()

	emoji.Println(*emojiKeyword)
}
