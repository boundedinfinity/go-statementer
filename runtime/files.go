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
	"github.com/boundedinfinity/go-commoner/idiomatic/slicer"
	"github.com/boundedinfinity/go-commoner/idiomatic/stringer"
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
	var hashes []string

	for _, path := range file.SourcePaths {
		hash, err := fileHash(path)
		if err != nil {
			return err
		}

		hashes = append(hashes, hash)
	}

	if len(hashes) > 1 {
		head, _ := slicer.Head(hashes...)
		tail, _ := slicer.Tail(hashes...)

		for _, hash := range tail {
			if head != hash {
				return fmt.Errorf(
					"hashes don't match: %s",
					stringer.Join(", ", file.SourcePaths...),
				)
			}
		}
	}

	file.Hash = hashes[0]
	return nil
}

func (this *Runtime) WalkSource() error {
	err := this.walkSource(func(path string, info fs.FileInfo, err error) error {
		found := this.State.Files.BySourcePath(path)
		var file *model.FileDescriptor

		switch len(found) {
		case 0:
			file = model.NewFileDescriptor()
			file.SourcePaths = []string{path}
			file.Id = uuid.New()
			this.State.Files = append(this.State.Files, file)
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

func (this *Runtime) FilesDuplicates() map[string][]*model.FileDescriptor {
	return this.State.Files.Duplicates()
}

func (this *Runtime) FilesAll() model.FileDescriptors {
	return this.State.Files
}

func (this *Runtime) FilesAllFiltered() model.FileDescriptors {
	if len(this.Labels.Selected) == 0 {
		return this.State.Files
	}

	var files model.FileDescriptors

	for _, file := range this.State.Files {
		if model.SimpleLabelsSame(file.Labels, this.Labels.Selected) {
			files = append(files, file)
		}
	}

	return files
}

func (this *Runtime) FileGet(id string) []*model.FileDescriptor {
	return this.State.Files.ById(id)
}
