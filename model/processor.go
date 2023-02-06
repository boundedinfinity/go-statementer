package model

import "github.com/oriser/regroup"

type AccountClassifier struct {
	Account string `yaml:"account"`
}

type Processor interface {
	Name() string
	Lines() []*LineDescriptor
	Extract(string) error
	Convert() error
	Print()
}

type FieldDescriptor struct {
	Name string
	Key  string
}

func NewFieldAndKey(name, key string) *FieldDescriptor {
	return &FieldDescriptor{
		Name: name,
		Key:  key,
	}
}

func NewField(name string) *FieldDescriptor {
	return NewFieldAndKey(name, name)
}

type LineDescriptor struct {
	Name    string
	Pattern string
	Fields  []*FieldDescriptor
	Regex   *regroup.ReGroup
}

func NewLine(name, pattern string, fields ...*FieldDescriptor) *LineDescriptor {
	return &LineDescriptor{
		Name:    name,
		Pattern: pattern,
		Fields:  fields,
	}
}

func NewLineWithField(name, pattern string) *LineDescriptor {
	return NewLine(name, pattern, NewField(name))
}

func NewLineWithFieldAndKey(name, key, pattern string) *LineDescriptor {
	return NewLine(name, pattern, NewFieldAndKey(name, key))
}
