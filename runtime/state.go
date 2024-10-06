package runtime

import (
	"encoding/json"
	"os"

	"github.com/boundedinfinity/go-commoner/idiomatic/pather"
)

func (this *Runtime) LoadState() error {
	this.mutex.Lock()
	defer this.mutex.Unlock()

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

	if _, err := this.Labels.Add(false, this.State.Labels...); err != nil {
		return err
	}

	for _, label := range this.Config.Labels {
		if _, err := this.Labels.Add(false, label); err != nil {
			return err
		}
	}

	for _, file := range this.State.Files {
		if _, err := this.Labels.Add(true, file.Labels...); err != nil {
			return err
		}
	}

	return nil
}

func (this *Runtime) SaveState() error {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	if _, err := pather.Dirs.EnsureErr(this.Config.RepositoryDir); err != nil {
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
