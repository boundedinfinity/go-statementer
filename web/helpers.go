package web

import (
	"fmt"

	"github.com/boundedinfinity/go-commoner/idiomatic/pather"
	"github.com/boundedinfinity/go-commoner/idiomatic/stringer"
)

func attrId(vs ...string) string {
	return stringer.Join("-", vs...)
}

func attrPath(vs ...string) string {
	return pather.Paths.Join(vs...)
}

func print(format string, a ...any) string {
	return fmt.Sprintf(format, a...)
}
