package config

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestValidateKafkaValues(t *testing.T) {
	Convey("Given valid kafka configurations", t, func() {
		cfg := getDefaultConfig()

		Convey("When validateKafkaValues is called", func() {
			errs := cfg.KafkaConfig.validateKafkaValues()

			Convey("Then no error messages should be returned", func() {
				So(errs, ShouldBeEmpty)
			})
		})
	})

	Convey("Given an empty KAFKA_ADDR", t, func() {
		cfg := getDefaultConfig()
		cfg.KafkaConfig.BindAddr = []string{}

		Convey("When validateKafkaValues is called", func() {
			errs := cfg.KafkaConfig.validateKafkaValues()

			Convey("Then an error message should be returned", func() {
				So(errs, ShouldNotBeEmpty)
				So(errs, ShouldResemble, []string{"no KAFKA_ADDR given"})
			})
		})
	})

	Convey("Given an empty KAFKA_VERSION", t, func() {
		cfg := getDefaultConfig()
		cfg.KafkaConfig.Version = ""

		Convey("When validateKafkaValues is called", func() {
			errs := cfg.KafkaConfig.validateKafkaValues()

			Convey("Then an error message should be returned", func() {
				So(errs, ShouldNotBeEmpty)
				So(errs, ShouldResemble, []string{"no KAFKA_VERSION given"})
			})
		})
	})

	Convey("Given an invalid KAFKA_SEC_PROTO", t, func() {
		cfg := getDefaultConfig()
		cfg.KafkaConfig.SecProtocol = "invalid"

		Convey("When validateKafkaValues is called", func() {
			errs := cfg.KafkaConfig.validateKafkaValues()

			Convey("Then an error message should be returned", func() {
				So(errs, ShouldNotBeEmpty)
				So(errs, ShouldResemble, []string{"KAFKA_SEC_PROTO has invalid value"})
			})
		})
	})

	Convey("Given an empty KAFKA_SEC_CLIENT_CERT", t, func() {
		cfg := getDefaultConfig()
		cfg.KafkaConfig.SecClientCert = ""

		Convey("And KAFKA_SEC_CLIENT_KEY has been set", func() {
			cfg.KafkaConfig.SecClientKey = "test key"

			Convey("When validateKafkaValues is called", func() {
				errs := cfg.KafkaConfig.validateKafkaValues()

				Convey("Then an error message should be returned", func() {
					So(errs, ShouldNotBeEmpty)
					So(errs, ShouldResemble, []string{"no KAFKA_SEC_CLIENT_CERT given but got KAFKA_SEC_CLIENT_KEY"})
				})
			})
		})
	})

	Convey("Given an empty KAFKA_SEC_CLIENT_KEY", t, func() {
		cfg := getDefaultConfig()
		cfg.KafkaConfig.SecClientKey = ""

		Convey("And KAFKA_SEC_CLIENT_CERT has been set", func() {
			cfg.KafkaConfig.SecClientCert = "test cert"

			Convey("When validateKafkaValues is called", func() {
				errs := cfg.KafkaConfig.validateKafkaValues()

				Convey("Then an error message should be returned", func() {
					So(errs, ShouldNotBeEmpty)
					So(errs, ShouldResemble, []string{"no KAFKA_SEC_CLIENT_KEY given but got KAFKA_SEC_CLIENT_CERT"})
				})
			})
		})
	})

	Convey("Given more than one invalid kafka configuration", t, func() {
		cfg := getDefaultConfig()
		cfg.KafkaConfig.Version = ""
		cfg.KafkaConfig.SecProtocol = "invalid"

		Convey("When validateKafkaValues is called", func() {
			errs := cfg.KafkaConfig.validateKafkaValues()

			Convey("Then an error message should be returned", func() {
				So(errs, ShouldNotBeEmpty)
				So(errs, ShouldResemble, []string{"no KAFKA_VERSION given", "KAFKA_SEC_PROTO has invalid value"})
			})
		})
	})
}
