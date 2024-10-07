package model

import "github.com/google/uuid"

type State struct {
	Files          FileDescriptors `json:"files" yaml:"files"`
	Labels         []*SimpleLabel  `json:"labels" yaml:"labels"`
	SelectedLabels []uuid.UUID     `json:"selected-labels" yaml:"selected-labels"`
	Debug          bool            `json:"debug" yaml:"debug"`
}
