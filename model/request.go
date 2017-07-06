package model

// Request is the structure of each request sent to the observation extractor.
type Request struct {
	FileURL    string `avro:"file_url"`
	InstanceID string `avro:"instance_id"`
}
