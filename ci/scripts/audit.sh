#!/bin/bash -eux

export cwd=$(pwd)

pushd $cwd/dp-observation-extractor
  make audit
popd