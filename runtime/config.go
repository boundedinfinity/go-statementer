package runtime

import (
	"os"

	"github.com/boundedinfinity/go-commoner/environmenter"
	"github.com/boundedinfinity/go-commoner/extentioner"
	"gopkg.in/yaml.v3"
)

func (t *Runtime) LoadConfig(path string) error {
	bs, err := os.ReadFile(path)

	if err != nil {
		return err
	}

	if err := yaml.Unmarshal(bs, &t.config); err != nil {
		return err
	}

	t.normalize()

	return nil
}

func (t *Runtime) normalize() {
	t.config.InputPath = environmenter.Sub(t.config.InputPath)
	t.config.OutputPath = environmenter.Sub(t.config.OutputPath)
	t.config.WorkPath = environmenter.Sub(t.config.WorkPath)
	t.config.SumExt = extentioner.Normalize(t.config.SumExt)
	t.config.InputExt = extentioner.Normalize(t.config.InputExt)

	if t.config.IgnorePaths == nil {
		t.config.IgnorePaths = make([]string, 0)
	}

	for _, p := range t.config.IgnorePaths {
		t.config.IgnorePaths = append(t.config.IgnorePaths, environmenter.Sub(p))
	}
}
