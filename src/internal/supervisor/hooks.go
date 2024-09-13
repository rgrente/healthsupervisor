package supervisor

import (
	"healthsupervisor/internal/hooks"
	"log"
)

func (s *Supervisor) checkAndExecuteHooks() {
	for _, hook := range s.Hooks {
		if s.hookConditionsMet(hook) {
			log.Println("Hook conditions met for", hook.Name)
			s.executeHookActions(hook)
		}
	}
}

func (s *Supervisor) hookConditionsMet(hook *hooks.Hook) bool {
	for _, condition := range hook.Conditions {
		if condition.SupervisorHealthy == s.StatusOK && s.StatusOK != s.previousStatusOK {
			return true
		}
		// Add more condition checks as needed
	}
	return false
}

func (s *Supervisor) executeHookActions(hook *hooks.Hook) {
	for _, action := range hook.Actions {
		switch action.Kind {
		case "http":
			// Execute HTTP action
			// You'll need to implement this method in the hooks package
			log.Println("Executing HTTP action", action.Name)
			if err := hooks.ExecuteHTTPAction(action); err != nil {
				log.Println("Error executing HTTP action", action.Name, err)
			} else {
				log.Println("HTTP action", action.Name, "executed successfully")
			}
			// Add more action types as needed
		}
	}
}
