package runtime

import (
	"encoding/json"
	"os"

	"github.com/boundedinfinity/go-commoner/idiomatic/pather"
)

func (this *Runtime) LoadState() error {
	_, err := pather.Dirs.EnsureErr(this.Config.ProcessedDir)
	if err != nil {
		return err
	}

	ok, err := pather.Files.ExistsErr(this.statePath)
	if err != nil {
		return err
	}

	if ok {
		data, err := os.ReadFile(this.statePath)
		if err != nil {
			return err
		}

		if err := json.Unmarshal(data, &this.State); err != nil {
			return err
		}
	}

	for _, file := range this.State.Files {
		if _, err := this.Labels.Add(file.Labels...); err != nil {
			return err
		}
	}

	for _, label := range this.Config.Labels {
		if _, err := this.Labels.Add(&label); err != nil {
			return err
		}
	}

	return nil
}

func (this *Runtime) SaveState() error {
	if _, err := pather.Dirs.EnsureErr(this.Config.ProcessedDir); err != nil {
		return err
	}

	data, err := json.MarshalIndent(this.State, "", "    ")
	if err != nil {
		return err
	}

	if err := os.WriteFile(this.statePath, data, os.FileMode(0755)); err != nil {
		return err
	}

	return nil
}
