package observation

// Reader provides an common interface to read Observations
type Reader interface {
	Read() (*Observation, error)
}
