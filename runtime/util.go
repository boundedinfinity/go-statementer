package runtime

import (
	"fmt"
	"io"
	"os"

	"github.com/boundedinfinity/go-commoner/pather"
)

func prepOutDir(dir string) error {
	if pather.PathExists(dir) {
		ok, err := pather.IsDir(dir)

		if err != nil {
			return err
		}

		if ok {
			if err := os.RemoveAll(dir); err != nil {
				return err
			}
		} else {
			return fmt.Errorf("not a directory: %v", dir)
		}
	}

	if err := pather.DirEnsure(dir); err != nil {
		return err
	}

	return nil
}

func appendFile(dst string, src ...string) error {
	dstHandle, err := os.Create(dst)

	if err != nil {
		return err
	}

	defer dstHandle.Close()

	appendNext := func(s string) error {
		srcHandle, err := os.Open(s)

		if err != nil {
			return err
		}

		defer srcHandle.Close()

		_, err = io.Copy(dstHandle, srcHandle)

		if err != nil {
			return err
		}

		return nil
	}

	for _, s := range src {
		if err := appendNext(s); err != nil {
			return err
		}
	}

	return nil
}

func copyFile(dst, src string) error {
	dstHandle, err := os.Create(dst)

	if err != nil {
		return err
	}

	defer dstHandle.Close()

	srcHandle, err := os.Open(src)

	if err != nil {
		return err
	}

	defer srcHandle.Close()

	_, err = io.Copy(dstHandle, srcHandle)

	if err != nil {
		return err
	}

	return nil
}

// func (t *Runtime) orig2txt(path string) string {
// 	return fu.SwapExt(path, t.config.InputExt, "txt")
// }

// func (t *Runtime) orig2sum(path string) string {
// 	return fu.SwapExt(path, t.config.InputExt, t.config.SumExt)
// }

// func (t *Runtime) sum2orig(path string) string {
// 	return fu.SwapExt(path, t.config.SumExt, t.config.InputExt)
// }

// func (t Runtime) hasSumExt(path string) bool {
// 	return strings.HasSuffix(path, t.config.SumExt)
// }

// func (t Runtime) isIgnorePath(path string) bool {
// 	fn := func(s string) bool {
// 		return strings.Contains(path, s)
// 	}

// 	return su.ExistFn(t.config.IgnorePaths, fn)
// }

// func (t Runtime) hasInputExt(path string) bool {
// 	return strings.HasSuffix(path, t.config.InputExt)
// }
