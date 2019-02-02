package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/format"
	"log"
	"os"
	"text/template"
)

var pkgName string
var fileName string

func init() {
	log.SetFlags(log.Llongfile)

	flag.StringVar(&pkgName, "pkg", "main", "output package")
	flag.StringVar(&fileName, "o", "emoji_codemap.go", "output file")
	flag.Parse()
}

// TemplateData emoji_codemap.go template
type TemplateData struct {
	PkgName string
	CodeMap map[string]string
}

const templateMapCode = `
package {{.PkgName}}

// NOTE: THIS FILE WAS PRODUCED BY THE
// EMOJICODEMAP CODE GENERATION TOOL (github.com/kyokomi/emoji/cmd/generateEmojiCodeMap)
// DO NOT EDIT

// Mapping from character to concrete escape code.
var emojiCodeMap = map[string]string{
	{{range $key, $val := .CodeMap}}":{{$key}}:": {{$val}},
{{end}}
}
`

func createCodeMap() (map[string]string, error) {
	gemojiCodeMap, err := createGemojiCodeMap()
	if err != nil {
		return nil, err
	}

	emojoCodeMap, err := createEmojoCodeMap()
	if err != nil {
		return nil, err
	}
	for k, v := range emojoCodeMap {
		gemojiCodeMap[k] = v
	}

	unicodeorgCodeMap, err := createUnicodeorgMap()
	if err != nil {
		return nil, err
	}
	for k, v := range unicodeorgCodeMap {
		gemojiCodeMap[k] = v
	}

	return gemojiCodeMap, nil
}

func createCodeMapSource(pkgName string, emojiCodeMap map[string]string) ([]byte, error) {
	// Template GenerateSource

	var buf bytes.Buffer
	t := template.Must(template.New("template").Parse(templateMapCode))
	if err := t.Execute(&buf, TemplateData{PkgName: pkgName, CodeMap: emojiCodeMap}); err != nil {
		return nil, err
	}

	// gofmt

	bts, err := format.Source(buf.Bytes())
	if err != nil {
		fmt.Printf(string(buf.Bytes()))
		return nil, fmt.Errorf("gofmt: %s", err)
	}

	return bts, nil
}

func main() {
	codeMap, err := createCodeMap()
	if err != nil {
		log.Fatalln(err)
	}

	codeMapSource, err := createCodeMapSource(pkgName, codeMap)
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
