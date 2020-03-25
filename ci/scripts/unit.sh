#!/bin/bash -eux

cwd=$(pwd)

pushd $cwd/dp-observation-extractor
  make test
popd