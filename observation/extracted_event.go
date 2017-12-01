package observation

// ExtractedEvent is the data that is output for each observation extracted.
type ExtractedEvent struct {
	RowIndex   int64  `avro:"row_index"`
	Row        string `avro:"row"`
	InstanceID string `avro:"instance_id"`
}
