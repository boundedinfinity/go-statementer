package runtime

import (
	"io/fs"
	"path/filepath"

	"github.com/boundedinfinity/go-commoner/idiomatic/stringer"
)

func (this *Runtime) walkSource(fn func(path string, info fs.FileInfo, err error) error) error {
	return filepath.Walk(this.Config.SourceDir, func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		if stringer.StartsWith(path, this.Config.RepositoryDir) {
			return nil
		}

		var validExt bool

		for _, ext := range this.Config.AllowedExts {
			if stringer.EndsWith(path, ext) {
				validExt = true
				break
			}
		}

		if !validExt {
			return nil
		}

		if err := fn(path, info, err); err != nil {
			return err
		}

		return nil
	})
}
