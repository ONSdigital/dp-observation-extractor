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
				So(cfg.AWSRegion, ShouldEqual, "eu-west-1")
				So(cfg.BindAddr, ShouldEqual, ":21600")
				So(cfg.EncryptionDisabled, ShouldEqual, true)
				So(cfg.ErrorProducerTopic, ShouldEqual, "report-events")
				So(cfg.FileConsumerGroup, ShouldEqual, "dimensions-inserted")
				So(cfg.FileConsumerTopic, ShouldEqual, "dimensions-inserted")
				So(cfg.GracefulShutdownTimeout, ShouldEqual, time.Second*5)
				So(cfg.KafkaAddr, ShouldResemble, []string{"localhost:9092"})
				So(cfg.ObservationProducerTopic, ShouldEqual, "observation-extracted")
				So(cfg.VaultAddr, ShouldEqual, "http://localhost:8200")
				So(cfg.VaultPath, ShouldEqual, "secret/shared/psk")
				So(cfg.VaultToken, ShouldEqual, "")
			})
		})
	})
}
