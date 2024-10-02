package model

import (
	"github.com/boundedinfinity/go-commoner/idiomatic/slicer"
	"github.com/boundedinfinity/go-commoner/idiomatic/stringer"
	"github.com/google/uuid"
)

func NewFileDescriptor() *FileDescriptor {
	return &FileDescriptor{
		Labels:     Labels{},
		DateLabels: []DateLabel{},
	}
}

type FileDescriptor struct {
	Id         uuid.UUID   `json:"id" yaml:"id"`
	Title      string      `json:"title" yaml:"title"`
	SourcePath string      `json:"source-path" yaml:"source-path"`
	RepoPath   string      `json:"repo-path" yaml:"repo-path"`
	Size       Size        `json:"size" yaml:"size"`
	Extention  string      `json:"extention" yaml:"extention"`
	Labels     Labels      `json:"labels" yaml:"labels"`
	DateLabels []DateLabel `json:"date-labels" yaml:"date-labels"`
	Hash       string      `json:"hash" yaml:"hash"`
}

func fileTitleFilter(file *FileDescriptor, text string) bool {
	return stringer.Contains(file.Title, text)
}

func fileSourcePathFilter(file *FileDescriptor, text string) bool {
	return stringer.Contains(file.SourcePath, text)
}

func fileExtentionFilter(file *FileDescriptor, text string) bool {
	return stringer.Contains(file.Extention, text)
}

func fileLabelFilter(file *FileDescriptor, text string) bool {
	return file.Labels.contains(text, labelNameFilter, labelDescriptionFilter)
}

// FileDescriptors

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
	return this.filter(name, fileLabelFilter)
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
