package util

import (
	"fmt"
	"strings"
)

func PrintSep() string {
	sep := strings.Repeat("=", 120)
	return sep
}

func PrintLabeled(label string, v string) string {
	text := fmt.Sprintf("%30v", label)
	text = fmt.Sprintf("%v: %v", text, v)
	return text
}

func PrintLabeleds(label string, vs []string) {
	for i, v := range vs {
		ilabel := fmt.Sprintf("%v[%v]", label, i)
		PrintLabeled(ilabel, v)
	}
}
