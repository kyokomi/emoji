package emoji

import (
	"bytes"
	"fmt"
	"math"
	"strings"
	"sync"
	"testing"
)

const (
	beerKey  = ":beer:"
	beerText = " ビール!!!"
	flag     = ":flag-us:"
	plusOne  = ":+1:"
)

var testFText = "test " + emojize(beerKey) + beerText
var testText = emojize(beerKey) + beerText

func TestFlag(t *testing.T) {
	f := emojize(flag)
	expected := "\U0001f1fA\U0001f1f8"
	if f != expected {
		t.Error("emojize ", f, "!=", expected)
	}
}

func TestPlusOne(t *testing.T) {
	f := emojize(plusOne)
	expected := "\U0001f44d "
	if f != expected {
		t.Error("emojize ", f, "!=", expected)
	}
}

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

func TestRevCodeMap(t *testing.T) {
	m := RevCodeMap()
	if &emojiRevCodeMap == &m {
		t.Error("emojiRevCodeMap != EmojiRevCodeMap")
	}
}

func TestHasAlias(t *testing.T) {
	hasAlias := HasAlias(":+1:")
	if !hasAlias {
		t.Error(":+1: doesn't have an alias")
	}
	hasAlias = HasAlias(":no-good:")
	if hasAlias {
		t.Error(":no-good: has an alias")
	}
}

func TestNoramlizeShortCode(t *testing.T) {
	test := ":thumbs_up:"
	expected := ":+1:"
	normalized := NormalizeShortCode(test)
	if normalized != expected {
		t.Errorf("Normalized %q != %q", test, expected)
	}
	test = ":no-good:"
	normalized = NormalizeShortCode(test)
	if normalized != test {
		t.Errorf("Normalized %q != %q", test, normalized)
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

// Copyright 2016 The Hugo Authors. All rights reserved.
// source: https://github.com/spf13/hugo/blob/master/helpers/emoji_test.go

func BenchmarkFprint(b *testing.B) {
	f := func(in []byte) []byte {
		buff := getBuffer()
		defer putBuffer(buff)
		Fprint(buff, string(in))

		bc := make([]byte, buff.Len())
		copy(bc, buff.Bytes())
		return bc
	}

	doBenchmarkEmoji(b, f)
}

func BenchmarkSprint(b *testing.B) {
	f := func(in []byte) []byte {
		return []byte(Sprint(string(in)))
	}

	doBenchmarkEmoji(b, f)
}

func doBenchmarkEmoji(b *testing.B, f func(in []byte) []byte) {
	type input struct {
		in     []byte
		expect []byte
	}

	data := []struct {
		input  string
		expect string
	}{
		{"A :smile: a day", Sprint("A :smile: a day")},
		{"A :smile: and a :beer: day keeps the doctor away", Sprint("A :smile: and a :beer: day keeps the doctor away")},
		{"A :smile: a day and 10 " + strings.Repeat(":beer: ", 10), Sprint("A :smile: a day and 10 " + strings.Repeat(":beer: ", 10))},
		{"No smiles today.", "No smiles today."},
		{"No smiles for you or " + strings.Repeat("you ", 1000), "No smiles for you or " + strings.Repeat("you ", 1000)},
	}

	var in = make([]input, b.N*len(data))
	var cnt = 0
	for i := 0; i < b.N; i++ {
		for _, this := range data {
			in[cnt] = input{[]byte(this.input), []byte(this.expect)}
			cnt++
		}
	}

	b.ResetTimer()
	cnt = 0
	for i := 0; i < b.N; i++ {
		for j := range data {
			currIn := in[cnt]
			cnt++
			result := f(currIn.in)
			// The Emoji implementations gives slightly different output.
			diffLen := len(result) - len(currIn.expect)
			diffLen = int(math.Abs(float64(diffLen)))
			if diffLen > 30 {
				b.Fatalf("[%d] emoji std, got \n%q but expected \n%q", j, result, currIn.expect)
			}
		}
	}
}

var bufferPool = &sync.Pool{
	New: func() interface{} {
		return &bytes.Buffer{}
	},
}

func getBuffer() (buf *bytes.Buffer) {
	return bufferPool.Get().(*bytes.Buffer)
}

func putBuffer(buf *bytes.Buffer) {
	buf.Reset()
	bufferPool.Put(buf)
}
