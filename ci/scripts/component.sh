#!/bin/bash -eux

pushd dp-observation-extractor
  make test-component
popd
