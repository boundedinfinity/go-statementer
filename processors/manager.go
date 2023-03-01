package processors

import (
	"github.com/boundedinfinity/docsorter/model"
	"github.com/sirupsen/logrus"
)

type ProcessManager struct {
	userConfig model.UserConfig
	logger     *logrus.Logger
	ocr        *model.ProcessContext
}

func NewManager(logger *logrus.Logger, userConfig model.UserConfig, ocr *model.ProcessContext) *ProcessManager {
	return &ProcessManager{
		logger:     logger,
		userConfig: userConfig,
		ocr:        ocr,
	}
}
