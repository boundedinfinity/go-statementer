package model

type State struct {
	Files          FileDescriptors `json:"files" yaml:"files"`
	SelectedLabels []*Label        `json:"selected-labels" yaml:"selected-labels"`
	Debug          bool            `json:"debug" yaml:"debug"`
}
