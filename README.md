dp-observation-extractor
================

* Consumes a Kafka message specifying a CSV file hosted on AWS S3
* Retrieves the file and produces a Kafka message for each row in the CSV

### Getting started

* Clone the repo `go get github.com/ONSdigital/dp-csv-splitter`
* Run the application `make debug`

#### Running in isolation
* run kafka consumer / producer apps 
* run local S3 store?

### Configuration

| Environment variable       | Default                 | Description
| ---------------------------| ----------------------- | ----------------------------------------------------
| BIND_ADDR                  | ":21600"                | The port to bind to
| KAFKA_ADDR                 | "http://localhost:9092" | The address of the Kafka instance
| FILE_CONSUMER_GROUP        | "dimensions-inserted"   | The Kafka consumer group to consume file messages from
| FILE_CONSUMER_TOPIC        | "dimensions-inserted"   | The Kafka topic to consume file messages from
| AWS_REGION                 | "eu-west-1"             | The AWS region to use
| OBSERVATION_PRODUCER_TOPIC | "observation-extracted" | The Kafka topic to send the observation messages to
| ERROR_PRODUCER_TOPIC       | "report-events"         | The kafka topic to send report event errors to
| GRACEFUL_SHUTDOWN_TIMEOUT  | "5s"                    | The shutdown timeout in seconds


### Contributing

See [CONTRIBUTING](CONTRIBUTING.md) for details.

### License

Copyright © 2016-2017, Office for National Statistics (https://www.ons.gov.uk)

Released under MIT license, see [LICENSE](LICENSE.md) for details.
