package config_test

import (
	"testing"
	"time"

	"github.com/ONSdigital/dp-observation-extractor/config"
	. "github.com/smartystreets/goconvey/convey"
)

func TestSpec(t *testing.T) {

	cfg, err := config.Get()

	Convey("Given an environment with no environment variables set", t, func() {

		Convey("When the config values are retrieved", func() {

			Convey("There should be no error returned", func() {
				So(err, ShouldBeNil)
			})

			Convey("The values should be set to the expected defaults", func() {
				So(*cfg, ShouldResemble, config.Config{
					AWSRegion:                   "eu-west-1",
					BindAddr:                    ":21600",
					EncryptionDisabled:          false,
					ErrorProducerTopic:          "report-events",
					FileConsumerGroup:           "dimensions-inserted",
					FileConsumerTopic:           "dimensions-inserted",
					GracefulShutdownTimeout:     time.Second * 5,
					KafkaAddr:                   []string{"localhost:9092"},
					ObservationProducerTopic:    "observation-extracted",
					VaultAddr:                   "http://localhost:8200",
					VaultPath:                   "secret/shared/psk",
					VaultToken:                  "",
					BucketNames:                 []string{"ons-dp-publishing-uploaded-datasets"},
					HealthCheckInterval:         10 * time.Second,
					HealthCheckRecoveryInterval: 1 * time.Minute,
				})
			})
		})
	})
}
