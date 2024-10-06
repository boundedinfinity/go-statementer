package web

import (
	"fmt"

	"github.com/boundedinfinity/go-commoner/idiomatic/pather"
	"github.com/boundedinfinity/go-commoner/idiomatic/stringer"
	"github.com/boundedinfinity/statementer/model"
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

func webSourcePath(config model.Config, path string) string {
	path = stringer.Replace(path, _PREFIX_PROCESSED_DIR, config.RepositoryDir)
	path = stringer.Replace(path, _PREFIX_SOURCE_DIR, config.SourceDir)
	return path
}
