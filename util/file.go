package util

import (
	"fmt"
	"io"
	"os"

	"github.com/boundedinfinity/go-commoner/pather"
)

func EnsureDelete(dir string) error {
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

	return nil
}

func AppendFile(dst string, src ...string) error {
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

func CopyFile(dst, src string) error {
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
