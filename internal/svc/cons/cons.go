package cons

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/rafaelkperes/tcc/pkg/data"
	"github.com/rafaelkperes/tcc/pkg/measure"
	log "github.com/sirupsen/logrus"
)

type Config struct {
}

func NewConsumerServer(cfg *Config) http.Handler {
	return newSrv(cfg)
}

type srv struct {
	*http.ServeMux
}

func newSrv(cfg *Config) *srv {
	mux := http.NewServeMux()
	s := &srv{
		ServeMux: mux,
	}

	// register routes
	mux.HandleFunc("/", s.handleRoot)

	return s
}

func (s *srv) handleRoot(w http.ResponseWriter, r *http.Request) {
	logger := log.WithFields(log.Fields{
		"handler": "root",
		"time":    time.Now(),
	})

	m := measure.NewMeasure()
	m.Add("recv")

	payload, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logger.WithFields(log.Fields{"event": "readRequestBody", "error": err, "args": r.Body}).
			Error("failed to read body")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	m.Add("rbod")

	d := data.Strings{}

	// TODO: change to generic unmarshalling
	err = json.Unmarshal(payload, &d)
	if err != nil {
		logger.WithFields(log.Fields{"event": "unmarshalRequestBody", "error": err, "args": []interface{}{payload, &d}}).
			Error("failed to unmarshal body")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	m.Add("dsrl")

	obj := m.AsObject()
	log.WithFields(log.Fields{"event": "measured", "measures": obj}).
		Info("add measures to response")

	j, err := json.Marshal(obj)
	if err != nil {
		logger.WithFields(log.Fields{"event": "marshalMeasures", "error": err, "args": obj}).
			Error("failed to marshal measures")
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(j)
	if err != nil {
		log.WithFields(log.Fields{"event": "writeResponseBody", "error": err, "args": j}).
			Error("failed to write response body")
	}
}
