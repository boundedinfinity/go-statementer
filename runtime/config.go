package runtime

import (
	"os"

	"github.com/boundedinfinity/go-commoner/environmenter"
	"github.com/boundedinfinity/go-commoner/extentioner"
	"github.com/boundedinfinity/go-commoner/slicer"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

func (t *Runtime) LoadUserConfig(path string) error {
	bs, err := os.ReadFile(path)

	if err != nil {
		return err
	}

	if err := yaml.Unmarshal(bs, &t.userConfig); err != nil {
		return err
	}

	t.normalizeUserConfig()

	t.logger.SetFormatter(&logrus.TextFormatter{
		DisableTimestamp: true,
	})

	switch t.userConfig.LogLevel {
	case "info":
		t.logger.SetLevel(logrus.InfoLevel)
	case "debug":
		t.logger.SetLevel(logrus.DebugLevel)
	case "trace":
		t.logger.SetLevel(logrus.TraceLevel)
	}

	return nil
}

func (t *Runtime) normalizeUserConfig() {
	t.userConfig.InputPaths = slicer.Map(t.userConfig.InputPaths, func(path string) string {
		return environmenter.Sub(path)
	})
	t.userConfig.OutputPath = environmenter.Sub(t.userConfig.OutputPath)
	t.userConfig.WorkPath = environmenter.Sub(t.userConfig.WorkPath)
	t.userConfig.SumExt = extentioner.Normalize(t.userConfig.SumExt)
	t.userConfig.InputExt = extentioner.Normalize(t.userConfig.InputExt)

	if t.userConfig.IgnorePaths == nil {
		t.userConfig.IgnorePaths = make([]string, 0)
	}

	for _, p := range t.userConfig.IgnorePaths {
		t.userConfig.IgnorePaths = append(t.userConfig.IgnorePaths, environmenter.Sub(p))
	}
}
