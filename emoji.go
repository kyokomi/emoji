// emoji.go
//
//      emoji.Println("@{:beer:}Example Text")
//
package emoji

import (
    "bytes"
    "errors"
    "fmt"
    "io"
    "log"
)

const (
    EscapeChar = '@'       // Escape character for emoji syntax
	ReplacePadding = " "
)

// Mapping from character to concrete escape code.
var CodeMap = map[string]string{
    ":beer:": "\xF0\x9f\x8d\xba",
	":pizza:": "\xF0\x9F\x8D\x95",
	":custard:": "\xF0\x9F\x8D\xAE",
}

func Emojize(x string) string {
    result := x

	str, ok := CodeMap[string(x)]
	switch {
	case !ok:
		log.Printf("Wrong emoji syntax: %c", x)
	default:
		result = str + ReplacePadding
	}
    return result
}

// Handle state after meeting one '@'
func compileEmojiSyntax(input, output *bytes.Buffer) {
    i, _, err := input.ReadRune()
    if err != nil {
        // EOF got
        log.Print("Parse failed on emoji syntax")
        return
    }

    switch i {
    default:
        output.WriteString(string(i))
    case '{':
        emoji := bytes.NewBufferString("")
        for {
            i, _, err := input.ReadRune()
            if err != nil {
                log.Print("Parse failed on emoji syntax")
                break
            }
            if i == '}' {
                break
            }
			emoji.WriteRune(i)
        }
        output.WriteString(Emojize(emoji.String()))
    }
}

func compile(x string) string {
    if x == "" {
        return ""
    }

    input := bytes.NewBufferString(x)
    output := bytes.NewBufferString("")

    for {
        i, _, err := input.ReadRune()
        if err != nil {
            break
        }
        switch i {
        default:
            output.WriteRune(i)
        case EscapeChar:
            compileEmojiSyntax(input, output)
        }
    }
    return output.String()
}

func compileValues(a *[]interface{}) {
    for i, x := range *a {
        if str, ok := x.(string); ok {
            (*a)[i] = compile(str)
        }
    }
}

func Print(a ...interface{}) (int, error) {
    compileValues(&a)
    return fmt.Print(a...)
}

func Println(a ...interface{}) (int, error) {
    compileValues(&a)
    return fmt.Println(a...)
}

func Printf(format string, a ...interface{}) (int, error) {
    format = compile(format)
    return fmt.Printf(format, a...)
}

func Fprint(w io.Writer, a ...interface{}) (int, error) {
    compileValues(&a)
    return fmt.Fprint(w, a...)
}

func Fprintln(w io.Writer, a ...interface{}) (int, error) {
    compileValues(&a)
    return fmt.Fprintln(w, a...)
}

func Fprintf(w io.Writer, format string, a ...interface{}) (int, error) {
    format = compile(format)
    return fmt.Fprintf(w, format, a...)
}

func Sprint(a ...interface{}) string {
    compileValues(&a)
    return fmt.Sprint(a...)
}

func Sprintf(format string, a ...interface{}) string {
    format = compile(format)
    return fmt.Sprintf(format, a...)
}

func Errorf(format string, a ...interface{}) error {
    return errors.New(Sprintf(format, a...))
}

