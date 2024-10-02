// Package runtime the runtime
package runtime

import (
	"github.com/boundedinfinity/statementer/model"
	"github.com/sirupsen/logrus"
)

func New(logger *logrus.Logger, configPath string) *Runtime {
	return &Runtime{
		logger:     logger,
		configPath: configPath,
	}
}

type Runtime struct {
	configPath string
	config     model.Config
	statePath  string
	state      model.State
	logger     *logrus.Logger
}
