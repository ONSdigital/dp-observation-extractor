package schema

import (
	"github.com/ONSdigital/go-ns/avro"
)

var dimensionsInsertedEvent = `{
  "type": "record",
  "name": "dimensions-inserted",
  "namespace": "",
  "fields": [
    {"name": "file_url", "type": "string"},
    {"name": "instance_id", "type": "string"}
  ]
}`

// DimensionsInsertedEvent the Avro schema for dimensionsInsertedEvent messages.
var DimensionsInsertedEvent *avro.Schema = &avro.Schema{
	Definition: dimensionsInsertedEvent,
}

var observationExtractedEvent = `{
  "type": "record",
  "name": "observation-extracted",
  "fields": [
    {"name": "instance_id", "type": "string"},
    {"name": "row", "type": "string"}
  ]
}`

// ObservationExtractedEvent is the Avro schema for each observation extracted.
var ObservationExtractedEvent *avro.Schema = &avro.Schema{
	Definition: observationExtractedEvent,
}
