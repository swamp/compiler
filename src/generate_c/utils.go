package generate_c

import (
	"strings"
)

func indentationString(indentation int) string {
	return strings.Repeat("    ", indentation)
}
