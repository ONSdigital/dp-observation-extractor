#!/bin/bash -eux

pushd dp-observation-extractor
  go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.64.6
  make lint
popd
