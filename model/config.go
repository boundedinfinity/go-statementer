// Package model the model
package model

import "github.com/boundedinfinity/statementer/label"

type Config struct {
	ConfigPath    string              `json:"-" yaml:"-"`
	StatePath     string              `json:"-" yaml:"-"`
	SourceDir     string              `json:"source-dir" yaml:"source-dir"`
	RepositoryDir string              `json:"repository-dir" yaml:"repository-dir"`
	AllowedExts   []string            `json:"allowed-exts" yaml:"allowed-exts"`
	Labels        []label.SimpleLabel `json:"labels" yaml:"labels"`
	DateLabels    []label.SimpleLabel `json:"date-labels" yaml:"date-labels"`
	Debug         bool                `json:"debug" yaml:"debug"`
}
