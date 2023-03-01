package processors

import (
	"github.com/boundedinfinity/docsorter/model"
	"github.com/sirupsen/logrus"
)

type ProcessManager struct {
	userConfig model.UserConfig
	logger     *logrus.Logger
	pc         *model.ProcessContext
}

func NewManager(logger *logrus.Logger, userConfig model.UserConfig, pc *model.ProcessContext) *ProcessManager {
	return &ProcessManager{
		logger:     logger,
		userConfig: userConfig,
		pc:         pc,
	}
}
