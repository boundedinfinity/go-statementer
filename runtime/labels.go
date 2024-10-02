package runtime

import "github.com/boundedinfinity/statementer/model"

func (this *Runtime) ResetLabels() error {
	this.state.Labels.Reset()
	this.state.DateLabels.Reset()

	for _, file := range this.state.Files {
		this.state.Labels.Update(file.Labels...)
		this.state.DateLabels.Update(model.DateLabels2Labels(file.DateLabels)...)
	}

	return nil
}
