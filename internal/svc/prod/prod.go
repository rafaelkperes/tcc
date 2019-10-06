package prod

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"github.com/rafaelkperes/tcc/pkg/data"
	"github.com/rafaelkperes/tcc/pkg/measure"
	log "github.com/sirupsen/logrus"
)

type Producer struct {
	endpoint string
	client   *http.Client
}

func NewProducer(consumerEndpoint string) *Producer {
	return &Producer{
		endpoint: consumerEndpoint,
		client:   &http.Client{},
	}
}

func (p *Producer) Produce(payload data.Data, format data.Format, n int, interval time.Duration) {
	log.Debugf("requesting a total of %d requests every %d milliseconds to %s", n, interval.Milliseconds(), p.endpoint)

	loop := func(idx int) {
		log.WithField("event", "progress").Debugf("sending request %d/%d", idx+1, n)

		m := measure.New()
		m.Add("strt")

		b, typ, err := payload.Marshal(format)
		if err != nil {
			log.Errorf("failed to marshal data: %v", err)
		}

		m.Add("srld")

		req, err := http.NewRequest(http.MethodPost, p.endpoint, bytes.NewReader(b))
		if err != nil {
			log.Errorf("failed to create request: %v", err)
			return
		}

		req.Header.Add("Content-Type", string(format))
		req.Header.Add("x-data-type", string(typ))

		r, err := p.client.Do(req)
		if err != nil {
			log.Errorf("got error on request: %v", err)
			return
		}
		defer r.Body.Close()

		m.Add("sent")

		rb, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Errorf("got error while reading response body: %v", err)
			return
		}

		obj := make(map[string]interface{})
		err = json.Unmarshal(rb, &obj)
		if err != nil {
			log.Errorf("failed to unmarshal consumer measures: %v", err)
			return
		}

		cm, err := measure.FromObject(obj)
		if err != nil {
			log.Errorf("failed to parse consumer measures: %v", err)
			return
		}
		log.Debugf("got consumer measures: %v", cm.AsObject())

		m.AddMeasures(cm)
		log.WithFields(map[string]interface{}{"event": "measured", "measures": m, "typ": typ, "format": format}).
			Info("finished measures")
	}

	if interval <= 0 {
		// non-concurrent, w/o interval usage
		for idx := 0; idx < n; idx++ {
			loop(idx)
			time.Sleep(interval)
		}
	} else {
		// interval-based concurrent requests
		var wg sync.WaitGroup
		wg.Add(n)
		for idx := 0; idx < n; idx++ {
			go func(idx int) {
				loop(idx)
				wg.Done()
			}(idx)
			time.Sleep(interval)
		}
		wg.Wait()
	}
	log.Debug("finished requests")
}
