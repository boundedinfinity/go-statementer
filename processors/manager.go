package processors

import (
	"github.com/boundedinfinity/docsorter/model"
	"github.com/sirupsen/logrus"
)

type ProcessManager struct {
	userConfig model.UserConfig
	logger     *logrus.Logger
	ocr        *model.OcrContext
	processors map[string]model.Processor
}

func NewManager(logger *logrus.Logger, userConfig model.UserConfig, ocr *model.OcrContext) *ProcessManager {
	return &ProcessManager{
		logger:     logger,
		userConfig: userConfig,
		ocr:        ocr,
	}
}
