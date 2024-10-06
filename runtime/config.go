package runtime

import (
	"encoding/json"
	"os"

	"github.com/boundedinfinity/go-commoner/idiomatic/environmenter"
	"github.com/boundedinfinity/go-commoner/idiomatic/pather"
)

var (
	_ENV_CONFIG_PATH = "STATEMENTER_CONFIG_PATH"
	_CONFIG_PATHS    = []string{
		"./statementer.config.json",
		"$HOME/.config/statementer/statementer.config.json",
		"$HOME/.config/statementer/config.json",
		"$HOME/statementer.config.json",
	}
)

func (this *Runtime) LoadConfig(path string) error {
	ev := environmenter.New()
	if err := ev.Process(); err != nil {
		return err
	}

	if err := this.findConfig(ev, path); err != nil {
		return err
	}

	data, err := os.ReadFile(this.configPath)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(data, &this.Config); err != nil {
		return err
	}

	this.Config.SourceDir = ev.Substitue(this.Config.SourceDir)
	this.Config.RepositoryDir = ev.Substitue(this.Config.RepositoryDir)
	this.statePath = pather.Paths.Join(this.Config.RepositoryDir, "state.json")

	return nil
}

func (this *Runtime) findConfig(ev *environmenter.Environmenter, path string) error {
	process := func(path string) bool {
		if path != "" {
			path = ev.Substitue(path)
			if pather.Files.Exists(path) {
				this.configPath = path
				return true
			}
		}

		return false
	}

	paths := []string{
		path,
		os.Getenv(_ENV_CONFIG_PATH),
	}

	paths = append(paths, _CONFIG_PATHS...)

	for _, path = range paths {
		if process(path) {
			break
		}
	}

	return nil
}
