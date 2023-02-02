package processors

import "github.com/boundedinfinity/rfc3339date"

type lineProcessor struct {
	name       string
	key        string
	matchFns   []func(string) map[string]string
	containFn  func(map[string]string) bool
	cleanFns   []func(map[string]string)
	extractFns []func(map[string]string)
	completed  bool
}

var _ Processor = &lineProcessor{}

func (t *lineProcessor) Completed() bool {
	return t.completed
}

func (t *lineProcessor) Process(line string) error {
	for _, match := range t.matchFns {
		m := match(line)

		if !t.containFn(m) {
			continue
		}

		for _, clean := range t.cleanFns {
			clean(m)
		}

		for _, extract := range t.extractFns {
			extract(m)
		}

		t.completed = true
	}

	return nil
}

type lineProcessorBuilder struct {
	processor *lineProcessor
}

func BuildLineProcessor(n string) *lineProcessorBuilder {
	return &lineProcessorBuilder{
		processor: &lineProcessor{
			name: n,
		},
	}
}

func (t *lineProcessorBuilder) Name(n string) *lineProcessorBuilder {
	t.processor.name = n
	return t
}

func (t *lineProcessorBuilder) ExtractString(k string, f *string) *lineProcessorBuilder {
	t.processor.extractFns = append(t.processor.extractFns, extractStringFn(k, f))
	return t
}

func (t *lineProcessorBuilder) ExtractDate(k string, f *rfc3339date.Rfc3339Date) *lineProcessorBuilder {
	t.processor.extractFns = append(t.processor.extractFns, extractDateFn(k, f))
	return t
}

func (t *lineProcessorBuilder) ExtractFloat(k string, f *float32) *lineProcessorBuilder {
	t.processor.extractFns = append(t.processor.extractFns, extractFloatFn(k, f))
	return t
}

func (t *lineProcessorBuilder) Clean(k string, fn ...func(string) string) *lineProcessorBuilder {
	t.processor.cleanFns = append(t.processor.cleanFns, cleanFn(k, fn...))
	return t
}

func (t *lineProcessorBuilder) Contains(k string) *lineProcessorBuilder {
	t.processor.containFn = containsFn(k)
	return t
}

func (t *lineProcessorBuilder) Match(k string, patterns ...string) *lineProcessorBuilder {
	t.processor.matchFns = append(t.processor.matchFns, matchFn(k, patterns...))
	return t
}

func (t *lineProcessorBuilder) Done() *lineProcessor {
	return t.processor
}
