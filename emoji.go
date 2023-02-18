// Package emoji terminal output.
package emoji

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"regexp"
	"unicode"

	"golang.org/x/text/transform"
)

//go:generate generateEmojiCodeMap -pkg emoji -o emoji_codemap.go

// Replace Padding character for emoji.
var (
	ReplacePadding = " "
)


var defaultTransfer = emojiTransfer{atEOF: true}

type emojiTransfer struct {
	atEOF    bool
	emojiBuf *bytes.Buffer
	buf      []byte
	transform.NopResetter
}

// NewEmojiTransfer return transform.Transfer implementd.
func NewEmojiTransfer() transform.Transformer {
	return &emojiTransfer{}
}

func (e *emojiTransfer) Transform(dst, src []byte, atEOF bool) (int, int, error) {
	e.atEOF = atEOF
	if e.buf == nil {
		e.buf = e.encode(bytes.NewBuffer(src)).Bytes()
	}

	n := copy(dst, e.buf)
	if len(e.buf) <= n {
		e.buf = nil
		return n, len(src), nil
	}
	e.buf = e.buf[n:]
	return n, 0, transform.ErrShortDst
}

// CodeMap gets the underlying map of emoji.
func CodeMap() map[string]string {
	return emojiCode()
}

// RevCodeMap gets the underlying map of emoji.
func RevCodeMap() map[string][]string {
	return emojiRevCode()
}

func AliasList(shortCode string) []string {
	return emojiRevCode()[emojiCode()[shortCode]]
}

// HasAlias flags if the given `shortCode` has multiple aliases with other
// codes.
func HasAlias(shortCode string) bool {
	return len(AliasList(shortCode)) > 1
}

// NormalizeShortCode normalizes a given `shortCode` to a deterministic alias.
func NormalizeShortCode(shortCode string) string {
	shortLists := AliasList(shortCode)
	if len(shortLists) == 0 {
		return shortCode
	}
	return shortLists[0]
}

// regular expression that matches :flag-[countrycode]:
var flagRegexp = regexp.MustCompile(":flag-([a-z]{2}):")

func emojize(x string) string {
	str, ok := emojiCode()[x]
	if ok {
		return str + ReplacePadding
	}
	if match := flagRegexp.FindStringSubmatch(x); len(match) == 2 {
		return regionalIndicator(match[1][0]) + regionalIndicator(match[1][1])
	}
	return x
}

// regionalIndicator maps a lowercase letter to a unicode regional indicator
func regionalIndicator(i byte) string {
	return string('\U0001F1E6' + rune(i) - 'a')
}

func compile(x string) string {
	if x == "" {
		return ""
	}

	return defaultTransfer.encode(bytes.NewBufferString(x)).String()
}

func (e *emojiTransfer) replaceEmoji(input *bytes.Buffer) string {
	emoji := bytes.NewBufferString(":")
	for {
		i, _, err := input.ReadRune()
		if err != nil {
			return e.replaceNotFindEnd(emoji)
		}

		if i == ':' && emoji.Len() == 1 {
			return emoji.String() + e.replaceEmoji(input)
		}

		emoji.WriteRune(i)
		switch {
		case unicode.IsSpace(i):
			return emoji.String()
		case i == ':':
			return emojize(emoji.String())
		}
	}
}

func (e *emojiTransfer) replaceNotFindEnd(emojiBuf *bytes.Buffer) string {
	if e.atEOF {
		return emojiBuf.String()
	}
	e.emojiBuf = emojiBuf
	return ""
}

func (e *emojiTransfer) mergeBuf(input *bytes.Buffer) *bytes.Buffer {
	if e.emojiBuf == nil {
		return input
	}
	return bytes.NewBuffer(append(e.emojiBuf.Bytes(), input.Bytes()...))
}

func (e *emojiTransfer) encode(input *bytes.Buffer) *bytes.Buffer {
	target := e.mergeBuf(input)
	e.emojiBuf = nil

	output := &bytes.Buffer{}
	output.Grow(input.Len())

	for {
		i, _, err := target.ReadRune()
		if err != nil {
			break
		}
		switch i {
		default:
			output.WriteRune(i)
		case ':':
			output.WriteString(e.replaceEmoji(target))
		}
	}
	return output
}

// Print is fmt.Print which supports emoji
func Print(a ...interface{}) (int, error) {
	return fmt.Print(compile(fmt.Sprint(a...)))
}

// Println is fmt.Println which supports emoji
func Println(a ...interface{}) (int, error) {
	return fmt.Println(compile(fmt.Sprint(a...)))
}

// Printf is fmt.Printf which supports emoji
func Printf(format string, a ...interface{}) (int, error) {
	return fmt.Print(compile(fmt.Sprintf(format, a...)))
}

// Fprint is fmt.Fprint which supports emoji
func Fprint(w io.Writer, a ...interface{}) (int, error) {
	return fmt.Fprint(w, compile(fmt.Sprint(a...)))
}

// Fprintln is fmt.Fprintln which supports emoji
func Fprintln(w io.Writer, a ...interface{}) (int, error) {
	return fmt.Fprintln(w, compile(fmt.Sprint(a...)))
}

// Fprintf is fmt.Fprintf which supports emoji
func Fprintf(w io.Writer, format string, a ...interface{}) (int, error) {
	return fmt.Fprint(w, compile(fmt.Sprintf(format, a...)))
}

// Sprint is fmt.Sprint which supports emoji
func Sprint(a ...interface{}) string {
	return compile(fmt.Sprint(a...))
}

// Sprintf is fmt.Sprintf which supports emoji
func Sprintf(format string, a ...interface{}) string {
	return compile(fmt.Sprintf(format, a...))
}

// Errorf is fmt.Errorf which supports emoji
func Errorf(format string, a ...interface{}) error {
	return errors.New(compile(Sprintf(format, a...)))
}
