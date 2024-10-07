package runtime

import (
	"encoding/json"
	"os"

	"github.com/boundedinfinity/go-commoner/idiomatic/pather"
)

func (this *Runtime) LoadState() error {
	_, err := pather.Dirs.EnsureErr(this.Config.RepositoryDir)
	if err != nil {
		return err
	}

	ok, err := pather.Files.ExistsErr(this.Config.StatePath)
	if err != nil {
		return err
	}

	if ok {
		data, err := os.ReadFile(this.Config.StatePath)
		if err != nil {
			return err
		}

		if err := json.Unmarshal(data, &this.State); err != nil {
			return err
		}
	}

	if err := this.refreshLabels(); err != nil {
		return err
	}

	return nil
}

func (this *Runtime) SaveState() error {
	if _, err := pather.Dirs.EnsureErr(this.Config.RepositoryDir); err != nil {
		return err
	}

	if err := this.refreshLabels(); err != nil {
		return err
	}

	this.State.Labels = this.Labels.All()

	data, err := json.MarshalIndent(this.State, "", "    ")
	if err != nil {
		return err
	}

	if err := os.WriteFile(this.Config.StatePath, data, os.FileMode(0755)); err != nil {
		return err
	}

	return nil
}
