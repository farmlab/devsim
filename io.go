package devsim

// Reader interface
type Reader interface {
	Read(*Simulator)
}

// Writer interface
type Writer interface {
	Write(*Simulator)
}

type Value interface{}

// Input event 
type Event map[string]Value

// Observable ouptut
type Observations struct {
	data	map[string]Value
	Keys	[]string
}

// State variables:
// State variables are variables that describes the conditions of system components at any time.
// The state variables change with time in dynamic system models (Wallach D, 2014)
type State map[string]Value

// Explanatory Variables and Parameters:
// Parameters are quantities that are unknown and are not measured directly, they must be estimated.
// Explanatory variables include the variables that are measured directly or known. This varailbe may vary between
// situation but do not change during simulation (Wallach D, 2014)
type Parameter map[string]Value



