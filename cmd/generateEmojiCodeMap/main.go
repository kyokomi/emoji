package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/format"
	"log"
	"os"
	"sort"
	"text/template"
)

var pkgName string
var fileName string

func init() {
	log.SetFlags(log.Llongfile)

	flag.StringVar(&pkgName, "pkg", "emoji", "output package")
	flag.StringVar(&fileName, "o", "../../emoji_codemap.go", "output file")
	flag.Parse()
}

// TemplateData emoji_codemap.go template
type TemplateData struct {
	PkgName    string
	CodeMap    map[string]string
	RevCodeMap map[string][]string
}

const templateMapCode = `
package {{.PkgName}}

import (
	"sync"
)

// NOTE: THIS FILE WAS PRODUCED BY THE
// EMOJICODEMAP CODE GENERATION TOOL (github.com/kyokomi/emoji/cmd/generateEmojiCodeMap)
// DO NOT EDIT

var emojiCodeMap map[string]string
var emojiCodeMapInitOnce = sync.Once{}

func emojiCode() map[string]string {
	emojiCodeMapInitOnce.Do(func() {
		emojiCodeMap = map[string]string{
			{{range $key, $val := .CodeMap}}":{{$key}}:": {{$val}},
		{{end}}}
	})
	return emojiCodeMap
}

var emojiRevCodeMap map[string][]string
var emojiRevCodeMapInitOnce = sync.Once{}

func emojiRevCode() map[string][]string {
	emojiRevCodeMapInitOnce.Do(func() {
		emojiRevCodeMap = map[string][]string{
			{{range $key, $val := .RevCodeMap}} {{$key}}: { {{range $val}} ":{{.}}:", {{end}} },
		{{end}}}
	})
	return emojiRevCodeMap
}
`

func createCodeMap() (map[string]string, map[string][]string, error) {
	log.Printf("creating gemoji code map")
	emojiCodeMap, err := createGemojiCodeMap()
	if err != nil {
		return nil, nil, err
	}

	log.Printf("creating emojo code map")
	emojoCodeMap, err := createEmojoCodeMap()
	if err != nil {
		return nil, nil, err
	}
	for k, v := range emojoCodeMap {
		emojiCodeMap[k] = v
	}

	log.Printf("creating unicode code map")
	unicodeorgCodeMap, err := createUnicodeorgMap()
	if err != nil {
		return nil, nil, err
	}
	for k, v := range unicodeorgCodeMap {
		emojiCodeMap[k] = v
	}

	log.Printf("creating emoji code map")
	emojiDataCodeMap, err := createEmojiDataCodeMap()
	if err != nil {
		return nil, nil, err
	}
	for k, v := range emojiDataCodeMap {
		emojiCodeMap[k] = v
	}

	log.Printf("creating reverse emoji code map")
	emojiRevCodeMap := make(map[string][]string)
	for shortName, unicode := range emojiCodeMap {
		emojiRevCodeMap[unicode] = append(emojiRevCodeMap[unicode], shortName)
	}

	// ensure deterministic ordering for aliases
	for _, value := range emojiRevCodeMap {
		sort.Slice(value, func(i, j int) bool {
			if len(value[i]) == len(value[j]) {
				return value[i] < value[j]
			}
			return len(value[i]) < len(value[j])
		})
	}

	return emojiCodeMap, emojiRevCodeMap, nil
}

func createCodeMapSource(pkgName string, emojiCodeMap map[string]string, emojiRevCodeMap map[string][]string) ([]byte, error) {
	// Template GenerateSource

	var buf bytes.Buffer
	t := template.Must(template.New("template").Parse(templateMapCode))
	if err := t.Execute(&buf, TemplateData{PkgName: pkgName, CodeMap: emojiCodeMap, RevCodeMap: emojiRevCodeMap}); err != nil {
		return nil, err
	}

	// gofmt

	bts, err := format.Source(buf.Bytes())
	if err != nil {
		fmt.Print(buf.String())
		return nil, fmt.Errorf("gofmt: %s", err)
	}

	return bts, nil
}

func main() {
	emojiCodeMap, emojiRevCodeMap, err := createCodeMap()
	if err != nil {
		log.Fatalln(err)
	}

	codeMapSource, err := createCodeMapSource(pkgName, emojiCodeMap, emojiRevCodeMap)
	if err != nil {
		log.Fatalln(err)
	}

	os.Remove(fileName)

	file, err := os.Create(fileName)
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	if _, err := file.Write(codeMapSource); err != nil {
		log.Fatalln(err)
	}
}
