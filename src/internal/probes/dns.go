package probes

import (
	"log"
	"net"
	"time"

	"github.com/mitchellh/mapstructure"
)

type dnsProbe struct {
	genericProbe
	ExpectedIPs []interface{}
}

func NewDNSProbe(data map[string]interface{}) (*dnsProbe, error) {
	var gp genericProbe
	err := mapstructure.Decode(data, &gp)
	if err != nil {
		return nil, err
	}

	p := &dnsProbe{
		genericProbe: gp,
		ExpectedIPs:  data["expectedIPs"].([]interface{}),
	}
	return p, nil
}

func (p *dnsProbe) Run() {
	previousHealth := p.StatusOK
	p.lastPoll = time.Now()

	ips, err := net.LookupIP(p.Target)
	if err != nil {
		p.StatusOK = false
		p.Message = err.Error()
	} else {
		p.StatusOK = false
		for _, expectedIP := range p.ExpectedIPs {
			for _, ip := range ips {
				if ip.String() == expectedIP.(string) {
					p.StatusOK = true
					break
				}
			}
			if p.StatusOK {
				break
			}
		}

		if p.StatusOK {
			p.Message = "DNS resolution successful"
		} else {
			p.Message = "DNS resolution successful, but no matching IP found"
		}
	}

	if p.StatusOK != previousHealth && p.statusChangeChannel != nil {
		log.Println("Probe", p.GetName(), "status changed to", p.StatusOK)
		p.statusChangeChannel <- p
	}
}
