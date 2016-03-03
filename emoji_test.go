package emoji

import (
	"bytes"
	"fmt"
	"testing"
)

const (
	beerKey  = ":beer:"
	beerText = " ビール!!!"
)

var testFText = "test " + emojize(beerKey) + beerText
var testText = emojize(beerKey) + beerText

func TestMultiColons(t *testing.T) {
	var buf bytes.Buffer
	_, err := Fprint(&buf, "A :smile: and another: :smile:")
	if err != nil {
		t.Error("Fprint ", err)
	}

	testCase := "A " + emojize(":smile:") + " and another: " + emojize(":smile:")
	if buf.String() != testCase {
		t.Error("Fprint ", buf.String(), "!=", testCase)
	}
}

func TestContinuityColons(t *testing.T) {
	var buf bytes.Buffer
	_, err := Fprint(&buf, "::smile:")
	if err != nil {
		t.Error("Fprint ", err)
	}

	testCase := ":" + emojize(":smile:")
	if buf.String() != testCase {
		t.Error("Fprint ", buf.String(), "!=", testCase)
	}
}

func TestCodeMap(t *testing.T) {
	m := CodeMap()
	if &emojiCodeMap == &m {
		t.Error("emojiCodeMap != EmojiCodeMap")
	}
}

func TestPrint(t *testing.T) {
	_, err := Print(beerKey, beerText)
	if err != nil {
		t.Error("Print ", err)
	}
}

func TestPrintln(t *testing.T) {
	_, err := Println(beerKey, beerText)
	if err != nil {
		t.Error("Println ", err)
	}
}

func TestPrintf(t *testing.T) {
	_, err := Printf("%s "+beerKey+beerText, "test")
	if err != nil {
		t.Error("Printf ", err)
	}
}

func TestFprint(t *testing.T) {
	var buf bytes.Buffer
	_, err := Fprint(&buf, beerKey+beerText)
	if err != nil {
		t.Error("Fprint ", err)
	}

	if buf.String() != testText {
		t.Error("Fprint ", buf.String(), testText)
	}
}

func TestFprintln(t *testing.T) {
	var buf bytes.Buffer
	_, err := Fprintln(&buf, beerKey+beerText)
	if err != nil {
		t.Error("Fprintln ", err)
	}

	if buf.String() != (testText + "\n") {
		t.Error("Fprintln ", buf.String(), (testText + "\n"))
	}
}

func TestFprintf(t *testing.T) {
	var buf bytes.Buffer
	_, err := Fprintf(&buf, "%s "+beerKey+beerText, "test")
	if err != nil {
		t.Error("Fprintf ", err)
	}

	if buf.String() != testFText {
		t.Error("Fprintf ", buf.String(), testFText)
	}
}

func TestSprint(t *testing.T) {
	convertBeer := Sprint(beerKey, beerText)
	if convertBeer != testText {
		t.Error("Sprint ", convertBeer, testText)
	}
}

func TestSprintf(t *testing.T) {
	convertBeer := Sprintf("%s "+beerKey+beerText, "test")
	if convertBeer != testFText {
		t.Error("Sprintf ", convertBeer, testFText)
	}
}

func TestErrorf(t *testing.T) {
	error := Errorf("%s "+beerKey+beerText, "test")
	if error.Error() != testFText {
		t.Error("Errorf ", error, testFText)
	}
}

func TestSprintMulti(t *testing.T) {
	convertBeer := Sprint(beerKey, beerText, beerKey, beerText)
	if convertBeer != (testText + testText) {
		t.Error("Sprint ", convertBeer, testText)
	}
	fmt.Println(convertBeer)
}
