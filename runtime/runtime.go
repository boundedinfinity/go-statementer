package runtime

import (
	"github.com/boundedinfinity/docsorter/model"
	"github.com/sirupsen/logrus"
)

func New(logger *logrus.Logger) Runtime {
	return Runtime{
		extPdf:   ".pdf",
		extImage: ".png",
		extText:  ".txt",
		extCsv:   ".csv",
		extYaml:  ".yaml",
		logger:   logger,
	}
}

type Runtime struct {
	UserConfig model.UserConfig
	logger     *logrus.Logger
	extPdf     string
	extImage   string
	extText    string
	extCsv     string
	extYaml    string
}
