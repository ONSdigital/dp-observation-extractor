package event

// DimensionsInserted is the structure of each event consumed by the observation extractor.
type DimensionsInserted struct {
	FileURL    string `avro:"file_url"`
	InstanceID string `avro:"instance_id"`
}
