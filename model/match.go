package model

import "github.com/oriser/regroup"

type MatchConfig struct {
	Pattern string
	Extract func(map[string]string)
}

type MatchContext struct {
	Pattern string
	Extract func(map[string]string)
	Matcher *regroup.ReGroup
}
