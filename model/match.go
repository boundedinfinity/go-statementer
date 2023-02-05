package model

import "github.com/oriser/regroup"

type MatcherContext struct {
	MatcherConfig
	OnExtract MatchHandler
	Matcher   *regroup.ReGroup
}

type MatcherConfig struct {
	Name        string   `yaml:"name"`
	Pattern     string   `yaml:"pattern"`
	RegexGroups []string `yaml:"regexGroups"`
}

func NewMatcherConfig(name, pattern string, regexGroups ...string) MatcherConfig {
	return MatcherConfig{
		Name:        name,
		Pattern:     pattern,
		RegexGroups: regexGroups,
	}
}
