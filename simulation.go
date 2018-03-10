package devsim

type EvtChan interface {
	Evt() chan Event
}

type ObsChan interface {
	Obs() chan Observations
}

type Observer interface{
	Observe(chan Observations)
}
// We define an option type. It is a function that takes one argument, the Simulation we are operating on.
// Options are used to initialize simulation within a simulation
type Option func(Simulation)

// Simulation interface
type Simulation interface{
	EvtChan
	ObsChan
	Observer
	Run(chan int)
}


// Simulation implementation
//type BaseSimulation struct{
	//evtChan  chan Event
	//obsChan chan Observations
//}

//func (s *BaseSimulation) Evt() chan Event {
	//return s.evtChan
//}

//func (s *BaseSimulation) Obs() chan Observations {
	//return s.obsChan
//}

// Observe whatch simulation obsChan to feed simulator writerChan
//func (s *DevsSimulation) Observe(writerChan chan Observations) {
	//for {
		//obs := <-s.Obs()
		//writerChan <- obs
	//}
//}
