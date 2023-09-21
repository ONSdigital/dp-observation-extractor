# dp-observation-extractor

* Consumes a Kafka message specifying a CSV file hosted on AWS S3
* Retrieves the file and produces a Kafka message for each row in the CSV

## Getting started

You may want vault to run this service:

- Run `brew install vault`
- Run `vault server -dev`

* Clone the repo `go get github.com/ONSdigital/dp-csv-splitter`
* Run the application `make debug`

### Running in isolation
* run kafka consumer / producer apps
* run local S3 store?

## Kafka scripts

Scripts for updating and debugging Kafka can be found [here](https://github.com/ONSdigital/dp-data-tools)(dp-data-tools)

## Configuration

| Environment variable         | Default                             | Description
| ---------------------------- | ----------------------------------- | ----------------------------------------------------
| BIND_ADDR                    | ":21600"                            | The port to bind to
| AWS_REGION                   | "eu-west-1"                         | The AWS region to use
| BUCKET_NAMES                 | ons-dp-publishing-uploaded-datasets | The expected S3 bucket names where the CSV files will be obtained from
| ENCRYPTION_DISABLED          | true                                | A boolean flag to identify if encryption of files is disabled or not
| GRACEFUL_SHUTDOWN_TIMEOUT    | "5s"                                | The shutdown timeout in seconds
| HEALTHCHECK_INTERVAL         | 30s                                 | The period of time between health checks
| HEALTHCHECK_CRITICAL_TIMEOUT | 90s                                 | The period of time after which failing checks will result in critical global 
| KAFKA_ADDR                   | "localhost:9092"                    | The addresses of the Kafka brokers (comma-separated)
| KAFKA_VERSION                | "1.0.2"                             | The kafka version that this service expects to connect to
| KAFKA_OFFSET_OLDEST          | true                                | set kafka offset to be oldest if `true`
| KAFKA_SEC_PROTO              | _unset_                             | if set to `TLS`, kafka connections will use TLS [[1]](#notes_1)
| KAFKA_SEC_CA_CERTS           | _unset_                             | CA cert chain for the server cert [[1]](#notes_1)
| KAFKA_SEC_CLIENT_KEY         | _unset_                             | PEM for the client key [[1]](#notes_1)
| KAFKA_SEC_CLIENT_CERT        | _unset_                             | PEM for the client certificate [[1]](#notes_1)
| KAFKA_SEC_SKIP_VERIFY        | false                               | ignores server certificate issues if `true` [[1]](#notes_1)
| LOCALSTACK_HOST              | ""                                  | Localstack to connect to for local S3 functionality
| ERROR_PRODUCER_TOPIC         | "report-events"                     | The Kafka topic to send report event errors to
| FILE_CONSUMER_GROUP          | "dimensions-inserted"               | The Kafka consumer group to consume file messages from
| FILE_CONSUMER_TOPIC          | "dimensions-inserted"               | The Kafka topic to consume file messages from
| OBSERVATION_PRODUCER_TOPIC   | "observation-extracted"             | The Kafka topic to send the observation messages to
| VAULT_ADDR                   | http://localhost:8200               | The vault address
| VAULT_TOKEN                  | -                                   | Vault token required for the client to talk to vault. (Use `make debug` to create a vault token)
| VAULT_PATH                   | secret/shared/psk                   | The path where the psks will be stored in for vault
check status

**Notes:**

 	1. <a name="notes_1">For more info, see the [kafka TLS examples documentation](https://github.com/ONSdigital/dp-kafka/tree/main/examples#tls)</a>

## Contributing

See [CONTRIBUTING](CONTRIBUTING.md) for details.

## License

Copyright © 2016-2021, Office for National Statistics (https://www.ons.gov.uk)

Released under MIT license, see [LICENSE](LICENSE.md) for details.
