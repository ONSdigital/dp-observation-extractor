package schema

import (
	"github.com/ONSdigital/go-ns/avro"
)

// Request the Avro schema for request messages.
var request = `{
  "type": "record",
  "name": "dimensions-inserted",
  "namespace": "",
  "fields": [
    {"name": "file_url", "type": "string"},
    {"name": "instance_id", "type": "string"}
  ]
}`

var Request *avro.Schema = &avro.Schema{
	Definition: request,
}

// ObservationExtractedEvent the Avro schema for each observation extracted.
var observationExtractedEvent = `{
  "type": "record",
  "name": "observation-extracted",
  "fields": [
    {"name": "instance_id", "type": "string"},
    {"name": "row", "type": "string"}
  ]
}`

var ObservationExtractedEvent *avro.Schema = &avro.Schema{
	Definition: observationExtractedEvent,
}
