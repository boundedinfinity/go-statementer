package model

type State struct {
	Files          FileDescriptors `json:"files" yaml:"files"`
	Labels         LabelMap        `json:"labels" yaml:"labels"`
	DateLabels     LabelMap        `json:"date-labels" yaml:"date-labels"`
	SelectedLabels []string        `json:"selected-labels" yaml:"selected-labels"`
}
