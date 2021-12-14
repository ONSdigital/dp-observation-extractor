package config_test

import (
	"errors"
	"os"
	"testing"
	"time"

	"github.com/ONSdigital/dp-observation-extractor/config"
	. "github.com/smartystreets/goconvey/convey"
)

func TestSpec(t *testing.T) {

	Convey("Given an environment with no environment variables set", t, func() {

		os.Clearenv()
		cfg, err := config.Get()

		Convey("When the config values are retrieved", func() {

			Convey("There should be no error returned", func() {
				So(err, ShouldBeNil)
			})

			Convey("The values should be set to the expected defaults", func() {
				So(*cfg, ShouldResemble, config.Config{
					BindAddr:                ":21600",
					AWSRegion:               "eu-west-1",
					BucketNames:             []string{"dp-frontend-florence-file-uploads"},
					EncryptionDisabled:      false,
					GracefulShutdownTimeout: time.Second * 5,
					HealthCheckInterval:     30 * time.Second,
					HealthCriticalTimeout:   90 * time.Second,
					KafkaConfig: config.KafkaConfig{
						Brokers:                  []string{"localhost:9092", "localhost:9093", "localhost:9094"},
						Version:                  "1.0.2",
						OffsetOldest:             true,
						SecProtocol:              "",
						SecCACerts:               "",
						SecClientCert:            "",
						SecClientKey:             "",
						SecSkipVerify:            false,
						ErrorProducerTopic:       "report-events",
						FileConsumerGroup:        "dimensions-inserted",
						FileConsumerTopic:        "dimensions-inserted",
						ObservationProducerTopic: "observation-extracted",
					},
					VaultAddr:  "http://localhost:8200",
					VaultToken: "",
					VaultPath:  "secret/shared/psk",
				})
			})
		})

		Convey("When configuration is called with an invalid security setting", func() {
			defer os.Clearenv()
			os.Setenv("KAFKA_SEC_PROTO", "ssl")
			cfg, err := config.Get()

			Convey("Then an error should be returned", func() {
				So(cfg, ShouldBeNil)
				So(err, ShouldNotBeNil)
				So(err, ShouldResemble, errors.New("kafka config validation errors: KAFKA_SEC_PROTO has invalid value"))
			})
		})
	})
}

func TestString(t *testing.T) {
	Convey("Given the config values", t, func() {
		cfg, err := config.Get()
		So(err, ShouldBeNil)

		Convey("When String is called", func() {
			cfgStr := cfg.String()

			Convey("Then the string format of config should not contain any sensitive configurations", func() {
				So(cfgStr, ShouldNotContainSubstring, "Brokers")
				So(cfgStr, ShouldNotContainSubstring, "SecClientKey")
				So(cfgStr, ShouldNotContainSubstring, "VaultToken")
				So(cfgStr, ShouldNotContainSubstring, "ServiceAuthToken")
				So(cfgStr, ShouldNotContainSubstring, "BucketNames")

				Convey("And should contain all non-sensitive configurations", func() {
					So(cfgStr, ShouldContainSubstring, "BindAddr")
					So(cfgStr, ShouldContainSubstring, "AWSRegion")
					So(cfgStr, ShouldContainSubstring, "EncryptionDisabled")
					So(cfgStr, ShouldContainSubstring, "GracefulShutdownTimeout")
					So(cfgStr, ShouldContainSubstring, "HealthCheckInterval")
					So(cfgStr, ShouldContainSubstring, "HealthCriticalTimeout")

					So(cfgStr, ShouldContainSubstring, "KafkaConfig")
					So(cfgStr, ShouldContainSubstring, "Version")
					So(cfgStr, ShouldContainSubstring, "SecProtocol")
					So(cfgStr, ShouldContainSubstring, "SecCACerts")
					So(cfgStr, ShouldContainSubstring, "SecClientCert")
					So(cfgStr, ShouldContainSubstring, "SecSkipVerify")
					So(cfgStr, ShouldContainSubstring, "ErrorProducerTopic")
					So(cfgStr, ShouldContainSubstring, "FileConsumerGroup")
					So(cfgStr, ShouldContainSubstring, "FileConsumerTopic")
					So(cfgStr, ShouldContainSubstring, "ObservationProducerTopic")

					So(cfgStr, ShouldContainSubstring, "VaultAddr")
					So(cfgStr, ShouldContainSubstring, "VaultPath")
				})
			})
		})
	})
}
