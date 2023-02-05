package runtime

import "github.com/boundedinfinity/docsorter/model"

func New() Runtime {
	return Runtime{
		imageExt: ".png",
		textExt:  ".txt",
	}
}

type Runtime struct {
	userConfig model.UserConfig
	imageExt   string
	textExt    string
}
