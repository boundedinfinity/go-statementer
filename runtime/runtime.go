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
		extCvs:   ".csv",
		logger:   logger,
	}
}

type Runtime struct {
	userConfig model.UserConfig
	logger     *logrus.Logger
	extPdf     string
	extImage   string
	extText    string
	extCvs     string
}
