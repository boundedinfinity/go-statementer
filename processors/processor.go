package processors

type Processor interface {
	Completed() bool
	Process(line string) error
}
