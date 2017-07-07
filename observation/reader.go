package observation

// Reader provides an interface to read Observations
type Reader interface {
	Read() (*Observation, error)
}
