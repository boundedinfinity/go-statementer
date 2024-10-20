package runtime

import (
	"encoding/json"
	"log"
	"os"

	"github.com/boundedinfinity/go-commoner/idiomatic/pather"
	"github.com/boundedinfinity/go-commoner/idiomatic/slicer"
	"github.com/boundedinfinity/statementer/label"
	"github.com/boundedinfinity/statementer/model"
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

	var state model.StateV2

	if ok {
		data, err := os.ReadFile(this.Config.StatePath)
		if err != nil {
			return err
		}

		if err := json.Unmarshal(data, &state); err != nil {
			return err
		}
	}

	if err := this.Labels.Add(label.Labels.P2M(state.Labels...)...); err != nil {
		return err
	}

	for _, id := range state.SelectedLabels {
		this.Labels.Select(true, id.String())
	}

	for _, label := range this.Config.Labels {
		if err := this.Labels.Add(&label); err != nil {
			return err
		}
	}

	if err := this.Labels.ResolveInit(); err != nil {
		return err
	}

	this.State.Files = model.Files.P2M(this.Labels, state.Files...)

	for _, file := range this.State.Files {
		this.Labels.Count(file.Labels...)
	}

	return nil
}

func (this *Runtime) SaveState() error {
	log.Println("save state: begin")
	defer log.Println("save state: end")

	if _, err := pather.Dirs.EnsureErr(this.Config.RepositoryDir); err != nil {
		return err
	}

	selectedLabels := slicer.Filter(label.SelectedFilter, this.Labels.List()...)

	state := model.StateV2{
		Version:        "2",
		Labels:         label.Labels.M2P(this.Labels.List()...),
		Files:          model.Files.M2P(this.State.Files...),
		SelectedLabels: slicer.Map(label.IdExtract, selectedLabels...),
	}

	data, err := json.MarshalIndent(state, "", "    ")
	if err != nil {
		return err
	}

	if err := os.WriteFile(this.Config.StatePath, data, os.FileMode(0644)); err != nil {
		return err
	}

	return nil
}
