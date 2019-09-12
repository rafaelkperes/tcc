package cons

import (
	"encoding/json"
	"time"
)

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

type jsonMeaures struct {
	Recv string `json:"recv"`
	Recb string `json:"recb"`
	Dsrl string `json:"dsrl"`
}

func (j jsonMeaures) asMeasures() (*Measures, error) {
	var err error

	m := NewMeasures()
	m.recv, err = time.Parse(time.RFC3339, j.Recv)
	if err != nil {
		return nil, err
	}

	m.recb, err = time.Parse(time.RFC3339, j.Recb)
	if err != nil {
		return nil, err
	}

	m.dsrl, err = time.Parse(time.RFC3339, j.Dsrl)
	if err != nil {
		return nil, err
	}

	return m, nil
}

func (m *Measures) ToJSON() []byte {
	b, err := json.Marshal(m.toJSONMeasures)
	if err != nil {
		panic(err) // should never fail
	}
	return b
}

func (m *Measures) toJSONMeasures() jsonMeaures {
	return jsonMeaures{
		Recv: m.recv.Format(time.RFC3339),
		Recb: m.recb.Format(time.RFC3339),
		Dsrl: m.dsrl.Format(time.RFC3339),
	}
}

func MeasuresFromJSON(b []byte) (*Measures, error) {
	j := jsonMeaures{}
	err := json.Unmarshal(b, &j)
	if err != nil {
		return nil, err
	}

	return j.asMeasures()
}
