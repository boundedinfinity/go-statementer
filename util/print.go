package util

import (
	"fmt"
	"strings"
)

func PrintSep() {
	sep := strings.Repeat("=", 120)
	fmt.Println(sep)
}

func PrintLabeled(label string, v string) {
	text := fmt.Sprintf("%30v", label)
	text = fmt.Sprintf("%v: %v", text, v)
	fmt.Println(text)
}

func PrintLabeleds(label string, vs []string) {
	for i, v := range vs {
		ilabel := fmt.Sprintf("%v[%v]", label, i)
		PrintLabeled(ilabel, v)
	}
}
