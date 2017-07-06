build:
	go build -o build/dp-observation-extractor
debug: build default-env
	 HUMAN_LOG=1 ./build/dp-observation-extractor

.PHONY: build debug default-env
