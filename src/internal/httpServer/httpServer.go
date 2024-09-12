package httpServer

import (
	"encoding/json"
	"healthsupervisor/internal/prober"
	"healthsupervisor/internal/supervisor"
	"net/http"
	"time"
)

type supervisorHealthRender struct {
	Name   string        `json:"name"`
	Status string        `json:"status"`
	Level  int           `json:"level"`
	Probes []probeRender `json:"probes"`
}

type probeRender struct {
	Name     string       `json:"name"`
	Status   statusRender `json:"status"`
	LastPoll string       `json:"lastPoll"`
}

type statusRender struct {
	Healthy bool   `json:"healthy"`
	Message string `json:"message"`
}

func HealthHandler(w http.ResponseWriter, r *http.Request, supervisor *supervisor.Supervisor, prober *prober.Prober) {

	probes := prober.GetProbes()
	probeStates := make([]probeRender, len(probes))

	for i, probe := range probes {
		lastPoll := probe.GetLastPoll().Format(time.RFC3339)
		probeStates[i] = probeRender{
			Name: probe.GetName(),
			Status: statusRender{
				Healthy: probe.GetProbeStatus().Healthy,
				Message: probe.GetProbeStatus().Message,
			},
			LastPoll: lastPoll,
		}
	}

	var supervisorStatus string
	if supervisor.Healthy {
		supervisorStatus = "UP"
	} else {
		supervisorStatus = "DOWN"
	}

	supervisorHealthRender := supervisorHealthRender{
		Name:   supervisor.Name,
		Status: supervisorStatus,
		Level:  supervisor.Level,
		Probes: probeStates,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(supervisorHealthRender)
}
