package cons

import (
	"encoding/json"
	"time"
)

const TimeFormat = time.RFC3339Nano

type Measures struct {
	recv time.Time
	recb time.Time
	dsrl time.Time
}

func NewMeasures() *Measures {
	return &Measures{}
}

func (m *Measures) ReceivedRequest() {
	m.recv = time.Now()
}

func (m *Measures) ReceivedBody() {
	m.recb = time.Now()
}

func (m *Measures) Desserialized() {
	m.dsrl = time.Now()
}

type jsonMeasures struct {
	Recv string `json:"recv"`
	Recb string `json:"recb"`
	Dsrl string `json:"dsrl"`
}

func (j jsonMeasures) asMeasures() (*Measures, error) {
	var err error

	m := NewMeasures()
	m.recv, err = time.Parse(TimeFormat, j.Recv)
	if err != nil {
		return nil, err
	}

	m.recb, err = time.Parse(TimeFormat, j.Recb)
	if err != nil {
		return nil, err
	}

	m.dsrl, err = time.Parse(TimeFormat, j.Dsrl)
	if err != nil {
		return nil, err
	}

	return m, nil
}

func (m *Measures) ToJSON() []byte {
	b, err := json.Marshal(m.toJSONMeasures())
	if err != nil {
		panic(err) // should never fail
	}
	return b
}

func (m *Measures) toJSONMeasures() jsonMeasures {
	return jsonMeasures{
		Recv: m.recv.Format(TimeFormat),
		Recb: m.recb.Format(TimeFormat),
		Dsrl: m.dsrl.Format(TimeFormat),
	}
}

func MeasuresFromJSON(b []byte) (*Measures, error) {
	j := jsonMeasures{}
	err := json.Unmarshal(b, &j)
	if err != nil {
		return nil, err
	}

	return j.asMeasures()
}
