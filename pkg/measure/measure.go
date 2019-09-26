package measure

import (
	"fmt"
	"time"
)

const (
	prefix     = "msr"
	timeFormat = time.RFC3339Nano
)

type Measure struct {
	t map[string]time.Time
}

func NewMeasure() *Measure {
	return &Measure{
		t: map[string]time.Time{fmt.Sprintf("%s.new", prefix): time.Now().UTC()},
	}
}

func (m *Measure) Add(key string) {
	m.t[fmt.Sprintf("%s.%s", prefix, key)] = time.Now().UTC()
}

func (m *Measure) AsObject() map[string]interface{} {
	o := make(map[string]interface{}, len(m.t))
	for k, v := range m.t {
		o[k] = v.Format(timeFormat)
	}
	return o
}

func FromObject(o map[string]interface{}) (*Measure, error) {
	t := make(map[string]time.Time, len(o))

	var err error
	var ok bool
	var s string
	for k, v := range o {
		if s, ok = v.(string); !ok {
			return nil, fmt.Errorf("on field '%s': expected string but found %T", k, v)
		}

		t[k], err = time.Parse(timeFormat, s)
		if err != nil {
			return nil, fmt.Errorf("on field '%s': %w", k, err)
		}
	}

	return &Measure{t: t}, nil
}
