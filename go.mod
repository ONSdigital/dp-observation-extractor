module github.com/ONSdigital/dp-observation-extractor

go 1.20

// to avoid 'sonatype-2021-4899' non-CVE Vulnerability
exclude github.com/gorilla/sessions v1.2.1

// to avoid 'sonatype-2020-1055' non-CVE vulnerability introduced by github.com/go-ldap/ldap/v3@v3.1.10
replace github.com/go-ldap/ldap/v3 v3.1.10 => github.com/go-ldap/ldap/v3 v3.4.3

replace github.com/prometheus/client_golang => github.com/prometheus/client_golang v1.11.1

//to avoid  [CVE-2022-29153] CWE-918: Server-Side Request Forgery (SSRF)
exclude github.com/hashicorp/consul/api v1.1.0

// to fix: [CVE-2023-32731] CWE-Other
replace google.golang.org/grpc => google.golang.org/grpc v1.55.0

require (
	github.com/ONSdigital/dp-healthcheck v1.6.1
	github.com/ONSdigital/dp-kafka/v2 v2.8.0
	github.com/ONSdigital/dp-reporter-client v1.1.0
	github.com/ONSdigital/dp-s3 v1.6.0
	github.com/ONSdigital/dp-vault v1.3.1
	github.com/ONSdigital/go-ns v0.0.0-20210916104633-ac1c1c52327e
	github.com/ONSdigital/log.go/v2 v2.4.1
	github.com/aws/aws-sdk-go v1.44.76
	github.com/golang/snappy v0.0.4 // indirect
	github.com/gopherjs/gopherjs v1.17.2 // indirect
	github.com/gorilla/mux v1.8.0
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-cleanhttp v0.5.2 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/hashicorp/go-retryablehttp v0.7.2 // indirect
	github.com/hashicorp/go-rootcerts v1.0.2 // indirect
	github.com/johnnadratowski/golang-neo4j-bolt-driver v0.0.0-20200323142034-807201386efa
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/smartystreets/assertions v1.13.1 // indirect
	github.com/smartystreets/goconvey v1.8.0
	golang.org/x/net v0.17.0 // indirect
	golang.org/x/sys v0.13.0 // indirect
	golang.org/x/time v0.3.0 // indirect
	gopkg.in/square/go-jose.v2 v2.6.0 // indirect
)

require (
	github.com/ONSdigital/dp-api-clients-go/v2 v2.254.0 // indirect
	github.com/ONSdigital/dp-net/v2 v2.9.1 // indirect
	github.com/ONSdigital/s3crypto v0.0.0-20180725145621-f8943119a487 // indirect
	github.com/Shopify/sarama v1.30.1 // indirect
	github.com/cenkalti/backoff/v3 v3.2.2 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/eapache/go-resiliency v1.2.0 // indirect
	github.com/eapache/go-xerial-snappy v0.0.0-20180814174437-776d5712da21 // indirect
	github.com/eapache/queue v1.1.0 // indirect
	github.com/fatih/color v1.15.0 // indirect
	github.com/go-avro/avro v0.0.0-20171219232920-444163702c11 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/hashicorp/go-secure-stdlib/parseutil v0.1.7 // indirect
	github.com/hashicorp/go-secure-stdlib/strutil v0.1.2 // indirect
	github.com/hashicorp/go-sockaddr v1.0.2 // indirect
	github.com/hashicorp/go-uuid v1.0.3 // indirect
	github.com/hashicorp/hcl v1.0.1-vault-5 // indirect
	github.com/hashicorp/vault/api v1.9.1 // indirect
	github.com/hokaccha/go-prettyjson v0.0.0-20211117102719-0474bc63780f // indirect
	github.com/jcmturner/aescts/v2 v2.0.0 // indirect
	github.com/jcmturner/dnsutils/v2 v2.0.0 // indirect
	github.com/jcmturner/gofork v1.0.0 // indirect
	github.com/jcmturner/gokrb5/v8 v8.4.2 // indirect
	github.com/jcmturner/rpc/v2 v2.0.3 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/jtolds/gls v4.20.0+incompatible // indirect
	github.com/justinas/alice v1.2.0 // indirect
	github.com/klauspost/compress v1.15.9 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.18 // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/pierrec/lz4 v2.6.1+incompatible // indirect
	github.com/rcrowley/go-metrics v0.0.0-20201227073835-cf1acfcdf475 // indirect
	github.com/ryanuber/go-glob v1.0.0 // indirect
	golang.org/x/crypto v0.14.0 // indirect
	golang.org/x/text v0.13.0 // indirect
)
