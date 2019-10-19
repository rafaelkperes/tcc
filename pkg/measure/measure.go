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

func (m Measure) Get(key string) (time.Time, bool) {
	tt, ok := m[key]
	return tt, ok
}

func (m Measure) GetString(key string) (string, bool) {
	tt, ok := m[key]
	if !ok {
		return "", false
	}
	return tt.Format(timeFormat), ok
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
	b, _ := json.Marshal(m.AsObject())
	return b
}

func (m Measure) String() string {
	return string(m.AsJSON())
}
