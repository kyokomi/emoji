package main

import (
	"fmt"
	"github.com/kyokomi/emoji"
)

func main() {
	fmt.Println("Hello Wolrd Emoji!")

	emoji.Println("@{:beer:} Beer!!!")

	pizzaMessage := emoji.Sprint("I like @{:pizza:}!!")
	fmt.Println(pizzaMessage)

	dessert := emoji.Sprintf("%s @{:custard:}.", "This is")
	fmt.Println(dessert)
}

