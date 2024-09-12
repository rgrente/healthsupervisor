package probes

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/mitchellh/mapstructure"
)

type remoteSupervisor struct {
	genericProbe
	httpClient http.Client
}

type remoteSupervisorResponse struct {
	Status string `json:"status"`
	Level  int    `json:"level"`
}

func NewRemoteSupervisor(data map[string]interface{}) (*remoteSupervisor, error) {
	var gp genericProbe
	err := mapstructure.Decode(data, &gp)
	if err != nil {
		return nil, err
	}

	p := &remoteSupervisor{
		genericProbe: gp,
		httpClient:   http.Client{},
	}

	return p, nil
}

func (p *remoteSupervisor) Run() {
	previousHealth := p.Healthy
	previousLevel := p.Level
	p.lastPoll = time.Now()

	r, err := p.httpClient.Get(p.Target)
	if err != nil {
		p.Healthy = false
		p.Message = err.Error()
		return
	}

	defer r.Body.Close()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		p.Healthy = false
		p.Message = err.Error()
		return
	}

	var rss remoteSupervisorResponse
	err = json.Unmarshal(body, &rss)
	if err != nil {
		p.Healthy = false
		p.Message = err.Error()
		return
	}

	if rss.Status != "UP" {
		p.Healthy = false
		p.Message = fmt.Sprintf("unexpected status: got %s, want %s", rss.Status, "UP")
	} else {
		p.Healthy = true
		p.Message = "OK"
		p.Level = rss.Level
	}

	if (p.Healthy != previousHealth || p.Level != previousLevel) && p.statusChangeChannel != nil {
		log.Println("Probe", p.GetName(), "health changed to", p.Healthy, "with level", p.Level)
		p.statusChangeChannel <- p
	}
}
