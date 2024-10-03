// Package model the model
package model

type Config struct {
	SourceDir    string   `json:"source-dir" yaml:"source-dir"`
	ProcessedDir string   `json:"processed-dir" yaml:"processed-dir"`
	AllowedExts  []string `json:"allowed-exts" yaml:"allowed-exts"`
	Labels       []Label  `json:"labels" yaml:"labels"`
	DateLabels   []Label  `json:"date-labels" yaml:"date-labels"`
}
