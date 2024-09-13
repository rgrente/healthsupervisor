package supervisor

// Function to evaluate local supervisor level
// This function is used to determine the level of the local supervisor by summing the weights of the healthy probes
func (s *Supervisor) evaluateLevel() {
	s.Level = 0
	for _, probe := range s.Prober.GetProbes() {
		if probe.GetProbeStatus().StatusOK {
			s.Level += probe.GetProbeWeight()
		}
	}
}
