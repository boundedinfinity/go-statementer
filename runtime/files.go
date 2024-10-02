package runtime

import (
	"crypto/sha512"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"

	"github.com/boundedinfinity/go-commoner/idiomatic/extentioner"
	"github.com/boundedinfinity/statementer/model"
	"github.com/google/uuid"
)

func fileHash(path string) (string, error) {
	hasher := sha512.New()
	file, err := os.Open(path)

	if err != nil {
		return "", err
	}

	defer file.Close()

	if _, err = io.Copy(hasher, file); err != nil {
		return "", err
	}

	sum := hasher.Sum(nil)
	return base64.URLEncoding.EncodeToString(sum), nil
}

func (this *Runtime) HashSource(file *model.FileDescriptor) error {
	hash, err := fileHash(file.SourcePath)
	if err != nil {
		return err
	}

	file.Hash = hash
	return nil
}

func (this *Runtime) WalkSource() error {
	err := this.walkSource(func(path string, info fs.FileInfo, err error) error {
		found := this.state.Files.BySourcePath(path)
		var file *model.FileDescriptor

		switch len(found) {
		case 0:
			file = model.NewFileDescriptor()
			file.SourcePath = path
			file.Id = uuid.New()
			this.state.Files = append(this.state.Files, file)
		case 1:
			file = found[0]
		default:
			return errors.New("dup file descriptors: " + path)
		}

		if err := this.HashSource(file); err != nil {
			return err
		}

		file.Size = model.NewSize(info.Size())
		file.Extention = extentioner.Ext(path)

		return nil
	})

	if err != nil {
		return err
	}

	if err := this.SaveState(); err != nil {
		return err
	}

	return nil
}

func (this *Runtime) ShowDups() {
	for hash, files := range this.state.Files.Duplicates() {
		fmt.Println(hash + ":")
		for _, file := range files {
			fmt.Printf("\t%s: Size: %s\n", file.SourcePath, file.Size.Human())
		}
	}
}

func (this *Runtime) Files() model.FileDescriptors {
	return this.state.Files
}
