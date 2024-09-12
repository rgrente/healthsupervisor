package probes

import (
	"time"

	"fmt"

	"github.com/mitchellh/mapstructure"
)

type Probe interface {
	Run()
	GetProbeStatus() ProbeStatus
	GetProbeWeight() int
	GetName() string
	GetInterval() int
	GetLastPoll() time.Time
	IsMandatory() bool
	SetStatusChangeChannel(chan<- Probe)
}

type genericProbe struct {
	Name               string "yaml:\"name\""
	Target             string "yaml:\"target\""
	Kind               string "yaml:\"kind\""
	ExpectedStatusCode int    "yaml:\"expectedStatusCode\""
	Weight             int    "yaml:\"weight\""
	Interval           int    "yaml:\"interval\""
	Mandatory          bool   "yaml:\"mandatory\""
	lastPoll           time.Time
	ProbeStatus
	statusChangeChannel chan<- Probe
}

type ProbeStatus struct {
	Healthy         bool
	PreviousHealthy bool
	Message         string
	Level           int
}

func NewGenericProbe(data map[string]interface{}) (*genericProbe, error) {
	var probe genericProbe
	err := mapstructure.Decode(data, &probe)
	if err != nil {
		return nil, err
	}
	return &probe, nil
}

func (p *genericProbe) GetName() string {
	return p.Name
}

func (p *genericProbe) GetInterval() int {
	return p.Interval
}

func (p *genericProbe) GetLastPoll() time.Time {
	return p.lastPoll
}

func (p *genericProbe) GetProbeStatus() ProbeStatus {
	return p.ProbeStatus
}

func (p *genericProbe) GetProbeWeight() int {
	return p.Weight
}

func (p *genericProbe) IsMandatory() bool {
	return p.Mandatory
}

func (p *genericProbe) SetStatusChangeChannel(ch chan<- Probe) {
	p.statusChangeChannel = ch
}

func NewProbeFromConfig(data map[string]interface{}) (Probe, error) {
	var genericProbe genericProbe
	err := mapstructure.Decode(data, &genericProbe)
	if err != nil {
		return nil, err
	}

	switch genericProbe.Kind {
	case "http":
		return NewHTTPProbe(data)
	case "remoteSupervisor":
		return NewRemoteSupervisor(data)
	// Add cases for other probe types here
	default:
		return nil, fmt.Errorf("unknown probe kind: %s", genericProbe.Kind)
	}
}
