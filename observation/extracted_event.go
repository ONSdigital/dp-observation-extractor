package observation

// ExtractedEvent is the data that is output for each observation extracted.
type ExtractedEvent struct {
	Row        string `avro:"row"`
	InstanceID string `avro:"instance_id"`
}
