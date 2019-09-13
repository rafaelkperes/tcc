package cons

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/rafaelkperes/tcc/pkg/data"
)

type Config struct {
}

func NewConsumerServer(cfg *Config) http.Handler {
	return newSrv(cfg)
}

type srv struct {
	*http.ServeMux

	cfg *Config
}

func newSrv(cfg *Config) *srv {
	mux := http.NewServeMux()
	s := &srv{
		ServeMux: mux,
		cfg:      cfg,
	}

	// register routes
	mux.HandleFunc("/", s.handleRoot)

	return s
}

func (s *srv) handleRoot(w http.ResponseWriter, r *http.Request) {
	log.Print("handleRoot")

	m := NewMeasures()
	m.ReceivedRequest()

	payload, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("failed to read body: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	m.ReceivedBody()

	d := data.Strings{}

	// TODO: change to generic unmarshalling
	err = json.Unmarshal(payload, &d)
	if err != nil {
		log.Printf("failed to unmarshal: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	m.Desserialized()

	log.Printf("respond with measures: %s", string(m.ToJSON()))
	w.WriteHeader(http.StatusOK)
	w.Write(m.ToJSON())
}
