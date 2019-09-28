package measure

import (
	"encoding/json"
	"fmt"
	"time"
)

const (
	timeFormat = time.RFC3339Nano
)

type Measure map[string]time.Time

func New() Measure {
	return map[string]time.Time{"new": time.Now().UTC()}
}

func (m Measure) Add(key string) {
	m[key] = time.Now().UTC()
}

func (m Measure) AddMeasures(other Measure) {
	for k, v := range other {
		m[k] = v
	}
}

func (m Measure) AsObject() map[string]interface{} {
	o := make(map[string]interface{}, len(m))
	for k, v := range m {
		o[k] = v.Format(timeFormat)
	}
	return o
}

func FromObject(o map[string]interface{}) (Measure, error) {
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

	return t, nil
}

func (m Measure) AsJSON() []byte {
	b, _ := json.Marshal(m)
	return b
}
