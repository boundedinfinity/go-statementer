// Package model the model
package model

type Config struct {
	SourceDir    string   `json:"source-dir" yaml:"source-dir"`
	ProcessedDir string   `json:"processed-dir" yaml:"processed-dir"`
	AllowedExts  []string `json:"allowed-exts" yaml:"allowed-exts"`
	Labels       []Label  `json:"labels" yaml:"labels"`
	DateLabels   []Label  `json:"date-labels" yaml:"date-labels"`
}

type State struct {
	Files          FileDescriptors `json:"files" yaml:"files"`
	Labels         LabelMap        `json:"labels" yaml:"labels"`
	DateLabels     LabelMap        `json:"date-labels" yaml:"date-labels"`
	SelectedLabels []string        `json:"selected-labels" yaml:"selected-labels"`
}
