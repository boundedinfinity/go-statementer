package runtime

import (
	"encoding/json"
	"os"

	"github.com/boundedinfinity/go-commoner/idiomatic/pather"
	"github.com/boundedinfinity/statementer/model"
)

func (this *Runtime) LoadState() error {
	_, err := pather.Dirs.EnsureErr(this.Config.RepositoryDir)
	if err != nil {
		return err
	}

	// var state1 model.StateV1
	// data, err := os.ReadFile(this.Config.StatePath)
	// if err != nil {
	// 	return err
	// }
	// if err := json.Unmarshal(data, &state1); err != nil {
	// 	return err
	// }

	// var state2 model.StateV2
	// state2.Files = model.Files.Model2Persist(this.Labels, state1.Files...)
	// state2.Labels = model.Labels.M2P(state1.Labels...)
	// state2.SelectedLabels = state1.SelectedLabels

	// data, err = json.MarshalIndent(state2, "", "    ")
	// if err != nil {
	// 	return err
	// }

	// if err := os.WriteFile(this.Config.StatePath, data, os.FileMode(0644)); err != nil {
	// 	return err
	// }

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

	if err := this.Labels.Add(model.Labels.P2M(state.Labels...)...); err != nil {
		return err
	}

	for _, id := range state.SelectedLabels {
		this.Labels.AddSelected(id)
	}

	for _, label := range this.Config.Labels {
		if err := this.Labels.Add(&label); err != nil {
			return err
		}
	}

	if err := this.Labels.ResolveParents(); err != nil {
		return err
	}

	this.State.Files = model.Files.Persist2Model(this.Labels, state.Files...)

	for _, file := range this.State.Files {
		this.Labels.Count(file.Labels...)
	}

	return nil
}

func (this *Runtime) SaveState() error {
	if _, err := pather.Dirs.EnsureErr(this.Config.RepositoryDir); err != nil {
		return err
	}

	state := model.StateV2{
		Version:        "2",
		Labels:         model.Labels.M2P(this.Labels.All()...),
		Files:          model.Files.Model2Persist(this.Labels, this.State.Files...),
		SelectedLabels: this.Labels.Selected,
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
