package probes

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/mitchellh/mapstructure"
)

type httpProbe struct {
	genericProbe
	httpClient http.Client
}

func NewHTTPProbe(data map[string]interface{}) (*httpProbe, error) {
	var gp genericProbe
	err := mapstructure.Decode(data, &gp)
	if err != nil {
		return nil, err
	}

	p := &httpProbe{
		genericProbe: gp,
		httpClient:   http.Client{},
	}
	return p, nil
}

func (p *httpProbe) Run() {
	previousHealth := p.Healthy
	p.lastPoll = time.Now()

	r, err := p.httpClient.Get(p.Target)
	if err != nil {
		p.Healthy = false
		p.Message = err.Error()
	} else if r.StatusCode != p.ExpectedStatusCode {
		p.Healthy = false
		p.Message = fmt.Sprintf("unexpected status code: got %d, want %d", r.StatusCode, p.ExpectedStatusCode)
	} else {
		p.Healthy = true
		p.Message = "OK"
	}

	if p.Healthy != previousHealth && p.statusChangeChannel != nil {
		log.Println("Probe", p.GetName(), "health changed to", p.Healthy)
		p.statusChangeChannel <- p
	}
}
