package emoji

import (
	"testing"
)

var beerKey  = ":beer:"
var beerText = " ビール!!!"

func TestPrintln(t *testing.T) {
	_, err := Println(createKeyStr(beerKey), beerText)
	if err != nil {
		t.Errorf("bad %s", err.Error())
	}
}

func TestSprint(t *testing.T) {
	convertBeer := Sprint(createKeyStr(beerKey), beerText)
	if convertBeer != (CodeMap[beerKey] + ReplacePadding + beerText) {
		t.Errorf("bad %s", convertBeer)
	}
}

func createKeyStr(key string) string {
	return  string(EscapeChar) + "{" + key + "}"
}
