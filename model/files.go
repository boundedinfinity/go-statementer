package model

import (
	"errors"
	"fmt"

	"github.com/boundedinfinity/go-commoner/idiomatic/slicer"
	"github.com/boundedinfinity/go-commoner/idiomatic/stringer"
	"github.com/google/uuid"
)

func NewFileDescriptor() *FileDescriptor {
	return &FileDescriptor{
		Labels: []*SimpleLabel{},
	}
}

type FileDescriptor struct {
	Id          uuid.UUID      `json:"id" yaml:"id"`
	Title       string         `json:"title" yaml:"title"`
	SourcePaths []string       `json:"source-path" yaml:"source-path"`
	RepoPath    string         `json:"repo-path" yaml:"repo-path"`
	Size        Size           `json:"size" yaml:"size"`
	Extention   string         `json:"extention" yaml:"extention"`
	Labels      []*SimpleLabel `json:"labels" yaml:"labels"`
	Hash        string         `json:"hash" yaml:"hash"`
}

func (this *FileDescriptor) Merge(that *FileDescriptor) error {
	if this.Hash == "" {
		return &ErrFileDescriptorDetails{
			msg:   "missing hash",
			files: []*FileDescriptor{this},
		}
	}

	if that.Hash == "" {
		return &ErrFileDescriptorDetails{
			msg:   "missing hash",
			files: []*FileDescriptor{that},
		}
	}

	if this.Hash != that.Hash {
		return &ErrFileDescriptorDetails{
			msg:   "hashes do not match",
			files: []*FileDescriptor{this, that},
		}
	}

	this.SourcePaths = append(this.SourcePaths, that.SourcePaths...)

	return nil
}

func fileIdFilter(file *FileDescriptor, id string) bool {
	return file.Id.String() == id
}

func fileTitleFilter(file *FileDescriptor, text string) bool {
	return stringer.Contains(file.Title, text)
}

func fileSourcePathFilter(file *FileDescriptor, text string) bool {
	return slicer.ContainsFn(func(_ int, path string) bool {
		return stringer.Contains(path, text)
	}, file.SourcePaths...)
}

func fileExtentionFilter(file *FileDescriptor, text string) bool {
	return stringer.Contains(file.Extention, text)
}

func fileLabelTermFilter(file *FileDescriptor, text string) bool {
	return slicer.ContainsFn(func(_ int, label *SimpleLabel) bool {
		return labelNameFilter(label, text) || labelDescriptionFilter(label, text)
	}, file.Labels...)
}

// =====================================================================================
// Errors
// =====================================================================================

var (
	ErrFileDescriptorErr = errors.New("file descriptor error")
)

type ErrFileDescriptorDetails struct {
	msg   string
	files []*FileDescriptor
}

func (e ErrFileDescriptorDetails) Error() string {
	lines := []string{ErrFileDescriptorErr.Error(), e.msg}

	var names []string

	if len(e.files) > 0 {
		for _, file := range e.files {
			names = append(names, file.SourcePaths...)
		}

		lines = append(lines, fmt.Sprintf("files - %s", stringer.Join(", ", names...)))
	}

	return stringer.Join(" : ", lines...)
}

func (e ErrFileDescriptorDetails) Unwrap() error {
	return ErrFileDescriptorErr
}

// =====================================================================================
// FileDescriptors
// =====================================================================================

type FileDescriptors []*FileDescriptor

func (this FileDescriptors) Duplicates() map[string][]*FileDescriptor {
	found := map[string][]*FileDescriptor{}

	group := slicer.Group(func(_ int, file *FileDescriptor) string {
		return file.Hash
	}, this...)

	for hash, files := range group {
		if len(files) > 1 {
			found[hash] = files
		}
	}

	return found
}

func (this FileDescriptors) ById(id string) []*FileDescriptor {
	return this.filter(id, fileIdFilter)
}

func (this FileDescriptors) ByTerm(text string) []*FileDescriptor {
	return this.filter(text, fileTitleFilter, fileSourcePathFilter, fileExtentionFilter)
}

func (this FileDescriptors) BySourcePath(name string) []*FileDescriptor {
	return this.filter(name, fileSourcePathFilter)
}

func (this FileDescriptors) ByTitle(name string) []*FileDescriptor {
	return this.filter(name, fileTitleFilter)
}

func (this FileDescriptors) ByExtention(name string) []*FileDescriptor {
	return this.filter(name, fileExtentionFilter)
}

func (this FileDescriptors) ByLabel(name string) []*FileDescriptor {
	return this.filter(name, fileLabelTermFilter)
}

func (this FileDescriptors) filter(text string, fns ...func(*FileDescriptor, string) bool) []*FileDescriptor {
	var found []*FileDescriptor

	for _, file := range this {
		for _, fn := range fns {
			if fn(file, text) {
				found = append(found, file)
			}
		}
	}

	return found
}
