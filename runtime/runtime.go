package runtime

import (
	"github.com/boundedinfinity/docsorter/model"
	"github.com/sirupsen/logrus"
)

func New(logger *logrus.Logger) Runtime {
	return Runtime{
		imageExt: ".png",
		textExt:  ".txt",
		logger:   logger,
	}
}

type Runtime struct {
	userConfig model.UserConfig
	imageExt   string
	textExt    string
	logger     *logrus.Logger
}
