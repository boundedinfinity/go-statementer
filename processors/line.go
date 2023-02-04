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
	if t.Completed() {
		return nil
	}

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

/////////////////////////////////////////////////////////////////////

type lineProcessorBuilder2 struct {
	processor *lineProcessor
}

func BuildLineProcessor2(name string, key string) *lineProcessorBuilder2 {
	return &lineProcessorBuilder2{
		processor: &lineProcessor{
			name: name,
			key:  key,
		},
	}
}

func (t *lineProcessorBuilder2) Extract(fn func(matches map[string]string)) *lineProcessorBuilder2 {
	t.processor.extractFns = append(t.processor.extractFns, fn)
	return t
}

func (t *lineProcessorBuilder2) ExtractMap(m map[string]*string) *lineProcessorBuilder2 {
	fn := func(matches map[string]string) {
		for matchKey, matchVal := range matches {
			if v, ok := m[matchKey]; ok {
				*v = matchVal
			}
		}
	}

	return t.Extract(fn)
}

func (t *lineProcessorBuilder2) ExtractString(value *string) *lineProcessorBuilder2 {
	return t.Extract(extractStringFn(t.processor.key, value))
}

func (t *lineProcessorBuilder2) ExtractDate(layout string, value *rfc3339date.Rfc3339Date) *lineProcessorBuilder2 {
	return t.Extract(extractDateFn(t.processor.key, layout, value))
}

func (t *lineProcessorBuilder2) ExtractFloat(value *float32) *lineProcessorBuilder2 {
	return t.Extract(extractFloatFn(t.processor.key, value))
}

func (t *lineProcessorBuilder2) Clean(fn ...func(string) string) *lineProcessorBuilder2 {
	t.processor.cleanFns = append(t.processor.cleanFns, cleanFn(t.processor.key, fn...))
	return t
}

func (t *lineProcessorBuilder2) Match(patterns ...string) *lineProcessorBuilder2 {
	t.processor.matchFns = append(t.processor.matchFns, matchFn(t.processor.key, patterns...))
	t.processor.containFn = containsFn(t.processor.key)
	return t
}

func (t *lineProcessorBuilder2) Done() *lineProcessor {
	return t.processor
}
