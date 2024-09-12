package supervisor

import (
	"healthsupervisor/internal/hooks"
	"healthsupervisor/internal/prober"
	"healthsupervisor/internal/rules"
)

type Supervisor struct {
	Name            string
	Healthy         bool
	previousHealthy bool
	Level           int
	Interval        int
	Prober          *prober.Prober
	Rules           []*rules.Rule
	Hooks           []*hooks.Hook
}

func NewSupervisor(name string, interval int, remoteSupervisors []map[string]interface{}, prober *prober.Prober, rules []*rules.Rule, hooksConfig []map[string]interface{}) (*Supervisor, error) {
	supervisor := &Supervisor{
		Name:            name,
		Healthy:         false,
		previousHealthy: false,
		Level:           0,
		Interval:        interval,
		Prober:          prober,
		Rules:           rules,
	}

	for _, hook := range hooksConfig {
		h, err := hooks.NewHook(hook)
		if err != nil {
			return nil, err
		}
		supervisor.Hooks = append(supervisor.Hooks, h)
	}
	return supervisor, nil
}

func (s *Supervisor) Run() {
	// // Run each remote supervisor check in a separate goroutine
	// for _, rs := range s.RemoteSupervisors {
	// 	go func(rs *remoteSupervisor) {
	// 		ticker := time.NewTicker(time.Duration(rs.Interval) * time.Second)
	// 		defer ticker.Stop()
	// 		for range ticker.C {
	// 			rs.run()
	// 			s.evaluateStatus()
	// 			s.evaluateLevel()
	// 		}
	// 	}(rs)
	// }

	// Listen for probe status changes
	go func() {
		for range s.Prober.StatusChanges() {
			s.evaluateRules()
			s.evaluateLevel()
			s.checkAndExecuteHooks()
		}
	}()

	// Periodic evaluation (you may want to keep this for consistency)
	// go func() {
	// 	ticker := time.NewTicker(time.Duration(s.Interval) * time.Second)
	// 	defer ticker.Stop()
	// 	for range ticker.C {
	// 		s.evaluateStatus()
	// 		s.evaluateLevel()
	// 	}
	// }()

	// // Add a new goroutine to check and execute hooks
	// go func() {
	// 	ticker := time.NewTicker(time.Second) // Check hooks every second
	// 	defer ticker.Stop()
	// 	for range ticker.C {
	// 		s.checkAndExecuteHooks()
	// 	}
	// }()
}
