dp-observation-extractor
================

* Consumes a Kafka message specifying a CSV file hosted on AWS S3
* Retrieves the file and produces a Kafka message for each row in the CSV

### Getting started

You may want vault to run this service:

- Run `brew install vault`
- Run `vault server -dev`

* Clone the repo `go get github.com/ONSdigital/dp-csv-splitter`
* Run the application `make debug`

#### Running in isolation
* run kafka consumer / producer apps
* run local S3 store?

### Kafka scripts

Scripts for updating and debugging Kafka can be found [here](https://github.com/ONSdigital/dp-data-tools)(dp-data-tools)

### Configuration

| Environment variable       | Default                 | Description
| ---------------------------| ----------------------- | ----------------------------------------------------
| AWS_REGION                 | "eu-west-1"             | The AWS region to use
| BIND_ADDR                  | ":21600"                | The port to bind to
| ENCRYPTION_DISABLED        | true                    | A boolean flag to identify if encryption of files is disabled or not
| ERROR_PRODUCER_TOPIC       | "report-events"         | The Kafka topic to send report event errors to
| FILE_CONSUMER_GROUP        | "dimensions-inserted"   | The Kafka consumer group to consume file messages from
| FILE_CONSUMER_TOPIC        | "dimensions-inserted"   | The Kafka topic to consume file messages from
| GRACEFUL_SHUTDOWN_TIMEOUT  | "5s"                    | The shutdown timeout in seconds
| KAFKA_ADDR                 | "http://localhost:9092" | The address of the Kafka instance
| OBSERVATION_PRODUCER_TOPIC | "observation-extracted" | The Kafka topic to send the observation messages to
| VAULT_ADDR                 | http://localhost:8200   | The vault address
| VAULT_TOKEN                | -                       | Vault token required for the client to talk to vault. (Use `make debug` to create a vault token)
| VAULT_PATH                 | secret/shared/psk       | The path where the psks will be stored in for vault
| AWS_ACCESS_KEY_ID          | -                       | The AWS access key credential for the observation extractor
| AWS_SECRET_ACCESS_KEY      | -                       | The AWS secret key credential for the observation extractor


### Contributing

See [CONTRIBUTING](CONTRIBUTING.md) for details.

### License

Copyright © 2016-2019, Office for National Statistics (https://www.ons.gov.uk)

Released under MIT license, see [LICENSE](LICENSE.md) for details.
