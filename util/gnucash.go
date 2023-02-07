package util

import (
	"strings"

	"github.com/boundedinfinity/go-commoner/slicer"
)

func FullAccountName2AccountName(s string) string {
	comps := strings.Split(s, "::")
	n := slicer.Last(comps)
	return n
}
