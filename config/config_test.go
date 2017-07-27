package config_test

import (
	"github.com/ONSdigital/dp-observation-extractor/config"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestSpec(t *testing.T) {

	cfg, err := config.Get()

	Convey("Given an environment with no environment variables set", t, func() {

		Convey("When the config values are retrieved", func() {

			Convey("There should be no error returned", func() {
				So(err, ShouldBeNil)
			})

			Convey("The values should be set to the expected defaults", func() {
				So(cfg.BindAddr, ShouldEqual, ":21600")
				So(cfg.KafkaAddr, ShouldEqual, "http://localhost:9092")
				So(cfg.FileConsumerGroup, ShouldEqual, "dimensions-inserted")
				So(cfg.FileConsumerTopic, ShouldEqual, "dimensions-inserted")
				So(cfg.AWSRegion, ShouldEqual, "eu-west-1")
				So(cfg.ObservationProducerTopic, ShouldEqual, "observation-io-extracted")
			})
		})
	})
}
