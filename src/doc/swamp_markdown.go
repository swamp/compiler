package doc

import (
	"bytes"
	"fmt"
	"log"
	"strings"

	"github.com/swamp/compiler/src/ast/codewriter"
	"github.com/swamp/compiler/src/parser"
	"github.com/swamp/compiler/src/runestream"
	"github.com/swamp/compiler/src/tokenize"
)

type ConvertState = uint8

const (
	Normal ConvertState = iota
	InAdmonition
	InCodeExample
)

func convertAdmonition(multiline string) (string, error) {
	var lines []string

	var admonitionLines []string

	var codeLines []string

	var admonitionHeader string

	convertState := Normal

	for _, line := range strings.Split(strings.TrimSuffix(multiline, "\n"), "\n") {
		switch {
		case convertState == InAdmonition:
			switch {
			case len(strings.TrimSpace(line)) == 0:
				continue

			case strings.HasPrefix(line, " "):
				rawLine := strings.TrimSpace(line)
				admonitionLines = append(admonitionLines, rawLine)
			default:
				parts := strings.Split(admonitionHeader, " ")
				lines = append(lines, fmt.Sprintf("<div class='admonition %v'>", parts[0]))
				lines = append(lines, "<p class='admonition-title'>warning</p>")
				lines = append(lines, "<p>")

				for _, admonitionLine := range admonitionLines {
					lines = append(lines, admonitionLine)
				}

				lines = append(lines, "</p></div>\n\n")

				convertState = Normal
				admonitionHeader = ""

				lines = append(lines, line)
			}
		case convertState == InCodeExample:
			switch {
			case len(strings.TrimSpace(line)) == 0:
				continue
			case strings.HasPrefix(strings.TrimSpace(line), "```"):
				swampCode := strings.Join(codeLines, "\n") + "\n"
				reader := strings.NewReader(swampCode)
				log.Printf("found:\n%v\n\n", swampCode)
				runeReader, runeErr := runestream.NewRuneReader(reader, "unknown filename")
				if runeErr != nil {
					return "", runeErr
				}
				tokenizer, tokenizeErr := tokenize.NewTokenizer(runeReader, true)
				if tokenizeErr != nil {
					return "", tokenizeErr
				}
				newParser := parser.NewParser(tokenizer, true)
				sourceFile, parseErr := newParser.Parse()
				if parseErr != nil {
					return "", parseErr
				}
				var buf bytes.Buffer
				colorer := &HtmlColorer{&buf}
				fmt.Fprintf(&buf, "\n\n<pre><code class=\"swamp\">\n")
				codewriter.WriteCodeUsingColorer(sourceFile, colorer, 0)
				fmt.Fprintf(&buf, "\n</pre></code>\n\n\n")
				log.Printf("code: \n'%v'\n", buf.String())
				lines = append(lines, buf.String())
				convertState = Normal
			default:
				log.Printf("adding:\n'%v'\n\n", line)
				codeLines = append(codeLines, line)
			}
		case strings.HasPrefix(line, "!!! "):
			admonitionHeader = line[4:]
			admonitionLines = nil
			convertState = InAdmonition
		case strings.HasPrefix(strings.TrimSpace(line), "```swamp"):
			convertState = InCodeExample
			codeLines = nil
		default:
			lines = append(lines, line)
		}
	}

	return strings.Join(lines, "\n"), nil
}

func ConvertSwampMarkdown(markdownString string) (string, error) {
	return convertAdmonition(markdownString)
}
