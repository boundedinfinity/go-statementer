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

	if err := yaml.Unmarshal(bs, &t.UserConfig); err != nil {
		return err
	}

	t.normalizeUserConfig()

	t.logger.SetFormatter(&logrus.TextFormatter{
		DisableTimestamp: true,
	})

	switch t.UserConfig.LogLevel {
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
	t.UserConfig.InputPaths = slicer.Map(t.UserConfig.InputPaths, func(path string) string {
		return environmenter.Sub(path)
	})

	t.UserConfig.OutputPath = environmenter.Sub(t.UserConfig.OutputPath)
	t.UserConfig.WorkPath = environmenter.Sub(t.UserConfig.WorkPath)
	t.UserConfig.SumExt = extentioner.Normalize(t.UserConfig.SumExt)
	t.UserConfig.InputExt = extentioner.Normalize(t.UserConfig.InputExt)

	if t.UserConfig.IgnorePaths == nil {
		t.UserConfig.IgnorePaths = make([]string, 0)
	}

	for _, p := range t.UserConfig.IgnorePaths {
		t.UserConfig.IgnorePaths = append(t.UserConfig.IgnorePaths, environmenter.Sub(p))
	}
}
