package runtime

func (this *Runtime) refreshLabels() error {
	this.Labels.Reset()

	if err := this.Labels.Add(this.State.Labels...); err != nil {
		return err
	}

	for _, label := range this.Config.Labels {
		if err := this.Labels.Add(&label); err != nil {
			return err
		}
	}

	for _, file := range this.State.Files {
		if err := this.Labels.Add(file.Labels...); err != nil {
			return err
		}

		this.Labels.Count(file.Labels...)
	}

	return nil
}
