module github.com/ONSdigital/dp-observation-extractor

go 1.23.2

// to avoid 'sonatype-2020-1055' non-CVE vulnerability introduced by github.com/go-ldap/ldap/v3@v3.1.10
replace github.com/go-ldap/ldap/v3 v3.1.10 => github.com/go-ldap/ldap/v3 v3.4.3

replace github.com/prometheus/client_golang => github.com/prometheus/client_golang v1.11.1

//to avoid  [CVE-2022-29153] CWE-918: Server-Side Request Forgery (SSRF)
exclude github.com/hashicorp/consul/api v1.1.0

// to fix: [CVE-2023-32731] CWE-Other
replace google.golang.org/grpc => google.golang.org/grpc v1.55.0

// [CVE-2024-24786] CWE-835: Loop with Unreachable Exit Condition ('Infinite Loop')
replace google.golang.org/protobuf => google.golang.org/protobuf v1.33.0

// [CVE-2023-48795] CWE-354: Improper Validation of Integrity Check Value
replace golang.org/x/crypto => golang.org/x/crypto v0.17.0

// [CVE-2024-6104] CWE-532: Information Exposure Through Log Files
replace github.com/hashicorp/go-retryablehttp => github.com/hashicorp/go-retryablehttp v0.7.7

require (
	github.com/ONSdigital/dp-healthcheck v1.6.3
	github.com/ONSdigital/dp-kafka/v2 v2.8.0
	github.com/ONSdigital/dp-reporter-client v1.2.0
	github.com/ONSdigital/dp-s3 v1.10.0
	github.com/ONSdigital/go-ns v0.0.0-20191104121206-f144c4ec2e58
	github.com/ONSdigital/log.go/v2 v2.4.3
	github.com/aws/aws-sdk-go v1.55.5
	github.com/golang/snappy v0.0.4 // indirect
	github.com/gopherjs/gopherjs v1.17.2 // indirect
	github.com/gorilla/mux v1.8.1
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-cleanhttp v0.5.2 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/hashicorp/go-rootcerts v1.0.2 // indirect
	github.com/johnnadratowski/golang-neo4j-bolt-driver v0.0.0-20200323142034-807201386efa
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/smartystreets/goconvey v1.8.1
	golang.org/x/net v0.19.0 // indirect
	golang.org/x/sys v0.20.0 // indirect
	golang.org/x/time v0.8.0 // indirect
)

require github.com/ONSdigital/dp-vault v1.3.1

require (
	github.com/ONSdigital/dp-api-clients-go/v2 v2.255.0 // indirect
	github.com/ONSdigital/dp-net/v2 v2.11.2 // indirect
	github.com/Shopify/sarama v1.38.1 // indirect
	github.com/cenkalti/backoff/v3 v3.2.2 // indirect
	github.com/cenkalti/backoff/v4 v4.3.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/eapache/go-resiliency v1.7.0 // indirect
	github.com/eapache/go-xerial-snappy v0.0.0-20230731223053-c322873962e3 // indirect
	github.com/eapache/queue v1.1.0 // indirect
	github.com/fatih/color v1.16.0 // indirect
	github.com/go-avro/avro v0.0.0-20171219232920-444163702c11 // indirect
	github.com/go-jose/go-jose/v4 v4.0.4 // indirect
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/hashicorp/go-retryablehttp v0.7.2 // indirect
	github.com/hashicorp/go-secure-stdlib/parseutil v0.1.8 // indirect
	github.com/hashicorp/go-secure-stdlib/strutil v0.1.2 // indirect
	github.com/hashicorp/go-sockaddr v1.0.7 // indirect
	github.com/hashicorp/go-uuid v1.0.3 // indirect
	github.com/hashicorp/hcl v1.0.1-vault-7 // indirect
	github.com/hashicorp/vault/api v1.9.1 // indirect
	github.com/hokaccha/go-prettyjson v0.0.0-20211117102719-0474bc63780f // indirect
	github.com/jcmturner/aescts/v2 v2.0.0 // indirect
	github.com/jcmturner/dnsutils/v2 v2.0.0 // indirect
	github.com/jcmturner/gofork v1.7.6 // indirect
	github.com/jcmturner/gokrb5/v8 v8.4.4 // indirect
	github.com/jcmturner/rpc/v2 v2.0.3 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/jtolds/gls v4.20.0+incompatible // indirect
	github.com/justinas/alice v1.2.0 // indirect
	github.com/klauspost/compress v1.17.11 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mgutz/ansi v0.0.0-20200706080929-d51e80ef957d // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/pierrec/lz4/v4 v4.1.21 // indirect
	github.com/rcrowley/go-metrics v0.0.0-20201227073835-cf1acfcdf475 // indirect
	github.com/ryanuber/go-glob v1.0.0 // indirect
	github.com/smarty/assertions v1.16.0 // indirect
	go.opentelemetry.io/otel v1.32.0 // indirect
	go.opentelemetry.io/otel/metric v1.32.0 // indirect
	go.opentelemetry.io/otel/trace v1.32.0 // indirect
	golang.org/x/crypto v0.29.0 // indirect
	golang.org/x/text v0.20.0 // indirect
	gopkg.in/square/go-jose.v2 v2.6.0 // indirect
)
