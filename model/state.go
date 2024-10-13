package model

import (
	"github.com/boundedinfinity/statementer/label"
	"github.com/google/uuid"
)

type StateDescriminator struct {
	Version string `json:"version" yaml:"version"`
}

type StateV1 struct {
	Version        string                  `json:"version" yaml:"version"`
	Files          FileDescriptors         `json:"files" yaml:"files"`
	Labels         []*label.LabelViewModel `json:"labels" yaml:"labels"`
	SelectedLabels []uuid.UUID             `json:"selected-labels" yaml:"selected-labels"`
	Debug          bool                    `json:"debug" yaml:"debug"`
}

type StateV2 struct {
	Version        string                                `json:"version" yaml:"version"`
	Files          []FilePersistenceModel                `json:"files" yaml:"files"`
	Labels         []label.SimpleLabelPersistenceModelV2 `json:"labels" yaml:"labels"`
	SelectedLabels []uuid.UUID                           `json:"selected-labels" yaml:"selected-labels"`
	Debug          bool                                  `json:"debug" yaml:"debug"`
}
