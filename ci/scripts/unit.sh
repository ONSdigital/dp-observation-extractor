#!/bin/bash -eux

export GOPATH=$(pwd)/go

pushd $GOPATH/src/github.com/ONSdigital/dp-observation-extractor
  make test
popd