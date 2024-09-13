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
	ExpectedStatusCode int
	httpClient         http.Client
}

func NewHTTPProbe(data map[string]interface{}) (*httpProbe, error) {
	var gp genericProbe
	err := mapstructure.Decode(data, &gp)
	if err != nil {
		return nil, err
	}

	p := &httpProbe{
		genericProbe:       gp,
		ExpectedStatusCode: data["expectedStatusCode"].(int),
		httpClient:         http.Client{},
	}
	return p, nil
}

func (p *httpProbe) Run() {
	previousHealth := p.StatusOK
	p.lastPoll = time.Now()

	r, err := p.httpClient.Get(p.Target)
	if err != nil {
		p.StatusOK = false
		p.Message = err.Error()
	} else if r.StatusCode != p.ExpectedStatusCode {
		p.StatusOK = false
		p.Message = fmt.Sprintf("unexpected status code: got %d, want %d", r.StatusCode, p.ExpectedStatusCode)
	} else {
		p.StatusOK = true
		p.Message = "OK"
	}

	if p.StatusOK != previousHealth && p.statusChangeChannel != nil {
		log.Println("Probe", p.GetName(), "status changed to", p.StatusOK)
		p.statusChangeChannel <- p
	}
}
