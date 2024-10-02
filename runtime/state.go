package runtime

import (
	"encoding/json"
	"os"

	"github.com/boundedinfinity/go-commoner/idiomatic/pather"
)

func (this *Runtime) LoadState() error {
	_, err := pather.Dirs.EnsureErr(this.config.ProcessedDir)
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

		if err := json.Unmarshal(data, &this.state); err != nil {
			return err
		}
	}

	this.ResetLabels()

	return nil
}

func (this *Runtime) SaveState() error {
	if _, err := pather.Dirs.EnsureErr(this.config.ProcessedDir); err != nil {
		return err
	}

	data, err := json.MarshalIndent(this.state, "", "    ")
	if err != nil {
		return err
	}

	if err := os.WriteFile(this.statePath, data, os.FileMode(0755)); err != nil {
		return err
	}

	return nil
}
