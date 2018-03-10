package devsim


import (
	"fmt"
)

// Set up the simulation registration services
var Simulators = []*Simulator{} // These are are registered Simulation

// Register a Simulator
func RegisterSimulator(s *Simulator) {
	Simulators = append(Simulators, s)
}

// Get a simulation from the id
func GetSimulator(uuid_str string) (*Simulator, error) {
	for _, s := range Simulators {
		if s.Uid.String() == uuid_str {
			return s, nil
		}
	}
	var nils *Simulator
	return nils, fmt.Errorf("Simulator with uuid '%s' not found", uuid_str)
}
