package supervisor

import (
	"healthsupervisor/internal/rules"
	"log"
)

// Function to evaluate rules
// This function is used to evaluate the rules of the local supervisor
func (s *Supervisor) evaluateRules() {
	s.previousHealthy = s.Healthy
	s.Healthy = true

	for _, rule := range s.Rules {
		switch rule.Kind {
		case "availability":
			s.evaluateHealthyRule(rule)
		case "remoteSupervisorEvaluation":
			s.evaluateRemoteSupervisorRule(rule)
		}
	}
	if s.Healthy != s.previousHealthy {
		log.Println("Supervisor health changed to", s.Healthy)
	}
}

// Function to evaluate Healthy rule
// This rule is used to determine the health of the local supervisor by evaluating the health of the probes
func (s *Supervisor) evaluateHealthyRule(rule *rules.Rule) {
	for _, probeName := range rule.Probes {
		probe := s.Prober.GetProbeByName(probeName)
		if probe == nil {
			// Log an error or handle the case where the probe is not found
			continue
		}
		if !probe.GetProbeStatus().Healthy {
			for _, ifNot := range rule.IfNot {
				if !*ifNot.SupervisorHealthy {
					s.Healthy = false
				}
			}
			return
		}
	}
}

// Function to evaluate RemoteSupervisorEvaluation rule
// This rule is used to determine the health of the local supervisor by evaluating the health and level of the remote supervisor
func (s *Supervisor) evaluateRemoteSupervisorRule(rule *rules.Rule) {
	for _, probeName := range rule.Probes {
		probe := s.Prober.GetProbeByName(probeName)
		if probe == nil {
			// Log an error or handle the case where the probe is not found
			continue
		}
		if probe.GetProbeStatus().Healthy {
			if probe.GetProbeStatus().Level > s.Level {
				s.Healthy = false
			}
		}
	}
}
