package model

import (
	"errors"
	"fmt"

	"github.com/boundedinfinity/go-commoner/idiomatic/slicer"
	"github.com/boundedinfinity/go-commoner/idiomatic/stringer"
	"github.com/boundedinfinity/statementer/label"
	"github.com/google/uuid"
)

func NewFileDescriptor() *FileDescriptor {
	return &FileDescriptor{
		Labels: []*label.LabelViewModel{},
	}
}

type FileDescriptor struct {
	Id          uuid.UUID               `json:"id" yaml:"id"`
	Title       string                  `json:"title" yaml:"title"`
	SourcePaths []string                `json:"source-path" yaml:"source-path"`
	RepoPath    string                  `json:"repo-path" yaml:"repo-path"`
	Size        Size                    `json:"size" yaml:"size"`
	Extention   string                  `json:"extention" yaml:"extention"`
	Labels      []*label.LabelViewModel `json:"labels" yaml:"labels"`
	Hash        string                  `json:"hash" yaml:"hash"`
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

func FileIdFilter(id string) FileFilterFunc {
	return func(i int, file *FileDescriptor) bool { return file.Id.String() == id }
}

func FileTitleFilter(term string) FileFilterFunc {
	return func(i int, file *FileDescriptor) bool { return stringer.Contains(file.Title, term) }
}

func FileSourcePathFilter(term string) FileFilterFunc {
	return func(i int, file *FileDescriptor) bool {
		return slicer.ContainsFn(func(_ int, path string) bool {
			return stringer.Contains(path, term)
		}, file.SourcePaths...)
	}
}

func FileExtentionFilter(term string) FileFilterFunc {
	return func(i int, file *FileDescriptor) bool { return stringer.Contains(file.Extention, term) }
}

func FileLabelTermFilter(term string) FileFilterFunc {
	return func(i int, file *FileDescriptor) bool {
		return slicer.ContainsFn(label.ContainsFilter(term), file.Labels...)
	}
}

func FileTermFilter(term string) FileFilterFunc {
	title := FileTitleFilter(term)
	sourcePath := FileSourcePathFilter(term)
	ext := FileExtentionFilter(term)

	return func(i int, file *FileDescriptor) bool {
		return title(i, file) || sourcePath(i, file) || ext(i, file)
	}
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

type FileFilterFunc func(int, *FileDescriptor) bool

func (this FileDescriptors) Filter(fns ...FileFilterFunc) []*FileDescriptor {
	files := this[:]

	for _, fn := range fns {
		files = slicer.Filter(fn, files...)
	}

	return files
}

// =====================================================================================
// File Persistence Model
// =====================================================================================

type FilePersistenceModel struct {
	Id          uuid.UUID   `json:"id" yaml:"id"`
	Title       string      `json:"title" yaml:"title"`
	SourcePaths []string    `json:"source-path" yaml:"source-path"`
	RepoPath    string      `json:"repo-path" yaml:"repo-path"`
	Size        Size        `json:"size" yaml:"size"`
	Extention   string      `json:"extention" yaml:"extention"`
	Labels      []uuid.UUID `json:"labels" yaml:"labels"`
	Hash        string      `json:"hash" yaml:"hash"`
}

// =====================================================================================
// Companion
// =====================================================================================

var Files = files{}

type files struct{}

func (this files) Model2Persist(lm *label.LabelManager, files ...*FileDescriptor) []FilePersistenceModel {
	return slicer.Map(func(_ int, file *FileDescriptor) FilePersistenceModel {
		return FilePersistenceModel{
			Id:          file.Id,
			Title:       file.Title,
			SourcePaths: file.SourcePaths,
			RepoPath:    file.RepoPath,
			Size:        file.Size,
			Extention:   file.Extention,
			Labels:      lm.Ids(file.Labels),
			Hash:        file.Hash,
		}
	}, files...)
}

func (this files) Model2Persist1(lm *label.LabelManager, file *FileDescriptor) FilePersistenceModel {
	return this.Model2Persist(lm, file)[0]
}

func (this files) Persist2Model(lm *label.LabelManager, files ...FilePersistenceModel) []*FileDescriptor {
	var descriptors []*FileDescriptor

	for _, persist := range files {
		descriptor := &FileDescriptor{
			Id:          persist.Id,
			Title:       persist.Title,
			SourcePaths: persist.SourcePaths,
			RepoPath:    persist.RepoPath,
			Size:        persist.Size,
			Extention:   persist.Extention,
			Hash:        persist.Hash,
		}

		for _, id := range persist.Labels {
			if label, ok := lm.ById(id); ok {
				descriptor.Labels = append(descriptor.Labels, label)
			}
		}

		descriptors = append(descriptors, descriptor)
	}

	return descriptors
}
