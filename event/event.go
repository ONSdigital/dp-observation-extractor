package event

// Event is the structure of each event sent to the observation extractor.
type Event struct {
	FileURL    string `avro:"file_url"`
	InstanceID string `avro:"instance_id"`
}
