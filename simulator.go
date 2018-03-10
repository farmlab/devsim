package devsim

import (
	log "github.com/Sirupsen/logrus"
	"github.com/nu7hatch/gouuid"
	"sync"
)


// Simulator goroutine waitgroup
var swg sync.WaitGroup 
// Wait for all simulator to be done, only if simulator are running
func WaitSimulators() {
	swg.Wait()
}

func InitSimulator(r Reader, w Writer, s Simulation, opts ...Option) *Simulator{
	// Simulation initializion 
	for _, opt := range opts {
		opt(s)
	}

	// Simulation	
	simulator := newSimulator(s)
	simulator.Attach(r, w)

	// Register simulator to be retrieved later
	RegisterSimulator(simulator) // TODO: removed simulation from regisery

	simulator.Log().Info("Ready")
	return simulator
}



// Possible states of simulator
const (
	Ready       = 0 // Simulator is ready to proceed the simulation, just after initialization
	Active      = 1 // Simulator is running
	Paused      = 2 // SImulator is paused
	Interrupted = 3 // Closed the simulator and all related goroutine

	Finished   = 4 // Simulation reach end time
	Processing = 5 // Simulation is running
	EndInput   = 6 // No more input, the InputChan was closed
	Fatal      = 7 // Simulation error
)


type Simulator struct{
	
	simulation Simulation
	
	stateChan chan int // Chan to manage the state of the simulator, trigered by the user
	doneChan  chan int // Chan used to stop the goroutine, used to close the simulator, trigger by the simulation
	state     int      // Current state of the simualtion
	lock      sync.RWMutex
	err       error

	Uid *uuid.UUID // Unique identifier used for the registery

	readerChan chan Event // Chan to manage inputs/events, from sourceReader to simulator
	reader Reader

	writerChan chan Observations // Chan to manage outputs, from simulator to targetWriter
	writer  Writer
}

func newSimulator(sim Simulation) *Simulator {
	s := &Simulator{
		simulation: sim,
		stateChan: make(chan int, Ready),
		state:     Ready,
		doneChan:  make(chan int),
		readerChan:    make(chan Event),
		writerChan:   make(chan Observations),
	}
	// Set uid
	u4, err := uuid.NewV4()
	if err != nil {
		log.Fatal(err)
	}
	s.Uid = u4
	return s
}


//This function bind to a simulator an input source
func (s *Simulator) AttachReader(r Reader) {
	s.reader = r
	go s.reader.Read(s)
}

//This function bind to a simulator an output source
func (s *Simulator) AttachWriter(w Writer) {
	s.writer = w
	go s.writer.Write(s)
}

//This function bind to a simulator an output source and an input source
func (s *Simulator) Attach(r Reader, w Writer) {
	s.AttachReader(r)
	s.AttachWriter(w)
}

// Stop tells the goroutine to stop.
func (t *Simulator) Stop() {
	t.stateChan <- Interrupted
}

// Start tells the goroutine to start after pausing.
func (s Simulator) Start() {
	s.lock.RLock()
	s.Log().Info(s.state)
	switch s.state {
	case Interrupted, Finished, Fatal:
		s.Log().Info("Simulator already Interrupted | Fatal | Finished")
		break
	case Active:
		s.Log().Info("Simulator already Active")
		break
	default:
		s.stateChan <- Active
	}
	s.lock.RUnlock()
}

// Pause tells the simulator goroutine to pause.
func (s Simulator) Pause() {
	s.lock.RLock()
	switch s.state {
	case Interrupted, Finished, Fatal:
		s.Log().Info("Simulator already Interrupted | Fatal | Finished")
		break
	case Paused:
		s.Log().Info("Simulator already Paused")
		break
	default:
		s.stateChan <- Paused
	}
	s.lock.RUnlock()
}

// RunState return the current state of the simulator goroutine.
func (t *Simulator) State() int {
	t.lock.RLock()
	state := t.state
	t.lock.RUnlock()
	return state
}

// Err gets the error returned by the simulator goroutine.
func (t *Simulator) Err() error {
	t.lock.RLock()
	err := t.err
	t.lock.RUnlock()
	return err
}

// log context
func (s *Simulator) Log() *log.Entry {
	return log.WithFields(log.Fields{
		"uid":       s.Uid,
		//"simulation": s.simulation.GetName(),
	})
}

// Run the simulator in a go routine which run, observe and manage one simulation
func (s *Simulator) Run() {
	swg.Add(1)
	go func() {
		defer swg.Done()
		//s.Execute()
		s.Log().Info("Started")
		// Run the simulation in a goroutine, it wait for inputs (see coordiante)
		go s.simulation.Run(s.doneChan)
		// Get simulation output
		go s.simulation.Observe(s.writerChan)
		s.coordinate()
	}()
	//return s.err
}

// Manage simulator state, triggered either by th user or by the simulation
func (s *Simulator) coordinate() error {
	var state = Active
	simulationLoop:
	for {
		select {
			// Manage simulator states triggered from user, i.e. Activeed, Paused, Interrupted
		case state = <-s.stateChan:
			switch state {
			case Interrupted:
				s.state = Interrupted
				s.Log().Info("Stopped")
				break simulationLoop
			case Paused:
				s.state = Paused
				s.Log().Info("Paused")
			case Active:
				s.state = Active
				s.Log().Info("Started")
			}
		default:
			if state == Paused {
				break
			}

			select { // Manage simulation inputs
			case events, ok := <-s.readerChan:
				if !ok { // if s.InputChan was closed by the reader, we close the simulation chan to end the simulation
					close(s.simulation.Evt())
				} else {
					s.simulation.Evt() <- events // Send input to the simulation goroutine
				}

				// Wait for simulation processing, return simulation status
				switch <-s.doneChan { // Manage simulation state triggered by the simulation, i.e Processing, Finished
				case Finished:
					s.state = Finished
					s.Log().Info("Finished")
					break simulationLoop
				case EndInput:
					s.state = EndInput
					s.Log().Info("End of inputs")
					break simulationLoop
				case Fatal:
					s.state = Fatal
					s.Log().Info("Fatal error")
					break simulationLoop
				}

			default:
				break
			}
		}
	}
	return s.err
}




