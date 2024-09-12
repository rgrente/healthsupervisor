package prober

import (
	"healthsupervisor/internal/probes"
	"log"
	"time"
)

type Prober struct {
	probes        []probes.Probe
	statusChanges chan probes.Probe
}

func NewProber(configuredProbes []map[string]interface{}) *Prober {
	prober := &Prober{
		statusChanges: make(chan probes.Probe, 100), // Buffered channel to avoid blocking
	}
	for _, probeConfig := range configuredProbes {
		probe, err := probes.NewProbeFromConfig(probeConfig)
		if err != nil {
			log.Printf("Error creating probe: %v", err)
			continue
		}
		probe.SetStatusChangeChannel(prober.statusChanges) // Set the status change channel
		log.Println("Configure probe", probe.GetName())
		prober.probes = append(prober.probes, probe)
	}
	return prober
}

func (p *Prober) Run() {
	for _, probe := range p.probes {
		go func(pr probes.Probe) {
			ticker := time.NewTicker(time.Duration(pr.GetInterval()) * time.Second)
			defer ticker.Stop()
			for range ticker.C {
				pr.Run()
			}
		}(probe)
	}
}

func (p *Prober) GetProbes() []probes.Probe {
	return append([]probes.Probe{}, p.probes...)
}

func (p *Prober) GetProbesStatus() []probes.ProbeStatus {
	states := make([]probes.ProbeStatus, 0)
	for _, probe := range p.probes {
		states = append(states, probe.GetProbeStatus())
	}
	return states
}

func (p *Prober) GetProbeByName(name string) probes.Probe {
	for _, probe := range p.probes {
		if probe.GetName() == name {
			return probe
		}
	}
	return nil
}

func (p *Prober) StatusChanges() <-chan probes.Probe {
	return p.statusChanges
}
