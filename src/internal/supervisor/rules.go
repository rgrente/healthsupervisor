package supervisor

import (
	"healthsupervisor/internal/rules"
	"log"
)

// Function to evaluate rules
// This function is used to evaluate the rules of the local supervisor
func (s *Supervisor) evaluateRules() {
	s.previousStatusOK = s.StatusOK
	s.StatusOK = true

	for _, rule := range s.Rules {
		switch rule.Kind {
		case "availability":
			s.evaluateHealthyRule(rule)
		case "remoteSupervisorEvaluation":
			s.evaluateRemoteSupervisorRule(rule)
		}
	}
	if s.StatusOK != s.previousStatusOK {
		log.Println("Supervisor status changed to", s.StatusOK)
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
		if !probe.GetProbeStatus().StatusOK {
			for _, ifNot := range rule.IfNot {
				if !*ifNot.SupervisorHealthy {
					s.StatusOK = false
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
		if probe.GetProbeStatus().StatusOK {
			if probe.GetProbeStatus().Level > s.Level {
				s.StatusOK = false
			}
		}
	}
}
