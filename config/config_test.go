package config_test

import (
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
					AWSRegion:                "eu-west-1",
					BindAddr:                 ":21600",
					EncryptionDisabled:       false,
					ErrorProducerTopic:       "report-events",
					FileConsumerGroup:        "dimensions-inserted",
					FileConsumerTopic:        "dimensions-inserted",
					GracefulShutdownTimeout:  time.Second * 5,
					KafkaAddr:                []string{"localhost:9092"},
					KafkaVersion:             "1.0.2",
					ObservationProducerTopic: "observation-extracted",
					VaultAddr:                "http://localhost:8200",
					VaultPath:                "secret/shared/psk",
					VaultToken:               "",
					BucketNames:              []string{"dp-frontend-florence-file-uploads"},
					HealthCheckInterval:      30 * time.Second,
					HealthCriticalTimeout:    90 * time.Second,
					KafkaOffset:              true,
				})
			})
		})
	})
}
