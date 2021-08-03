package doc

import (
	"fmt"
	"strings"
)

func convertAdmonition(multiline string) string {
	var lines []string

	var admonitionLines []string

	var admonitionHeader string

	isInAdmonition := false

	for _, line := range strings.Split(strings.TrimSuffix(multiline, "\n"), "\n") {
		switch {
		case isInAdmonition:
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
				isInAdmonition = false
				admonitionHeader = ""
				lines = append(lines, line)
			}
		case strings.HasPrefix(line, "!!! "):
			admonitionHeader = line[4:]
			admonitionLines = nil
			isInAdmonition = true
		default:
			lines = append(lines, line)
		}
	}

	return strings.Join(lines, "\n")
}

func ConvertSwampMarkdown(markdownString string) string {
	return convertAdmonition(markdownString)
}
