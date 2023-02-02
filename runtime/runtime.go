package runtime

import "github.com/boundedinfinity/docsorter/model"

func New() Runtime {
	return Runtime{}
}

type Runtime struct {
	config model.Config
}
