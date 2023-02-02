package processors

type StatementProcessor struct {
	Name       string
	processors []Processor
}

var _ Processor = &StatementProcessor{}

func (p *StatementProcessor) Process(line string) error {
	for _, processor := range p.processors {
		if !processor.Completed() {
			if err := processor.Process(line); err != nil {
				return err
			}
		}
	}

	return nil
}

func (p *StatementProcessor) Completed() bool {
	for _, processor := range p.processors {
		if !processor.Completed() {
			return false
		}
	}

	return true
}
