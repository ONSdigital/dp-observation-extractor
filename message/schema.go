package message

// RequestSchema the Avro schema for request messages.
var RequestSchema = `{
  "type": "record",
  "name": "dimensions-inserted",
  "namespace": "",
  "fields": [
    {"name": "file_url", "type": "string"},
    {"name": "instance_id", "type": "string"}
  ]
}`

// ObservationExtractedSchema the Avro schema for each observation extracted.
var ObservationExtractedSchema = `{
  "type": "record",
  "name": "observation-extracted",
  "fields": [
    {"name": "instance_id", "type": "string"},
    {"name": "row", "type": "string"}
  ]
}`
