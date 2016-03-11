#!/usr/bin/env sh
set -xe

# assume heka and heka-scamp are siblings in a directory
# docker build --no-cache=false --rm -t xrlx/heka_base ../heka

# we run most of the heka build one time so it's quick to do our plugin build
# docker build --rm -t xrlx/heka_base -f docker/Dockerfile.heka .

# prepare the scamp plugin using our heka_base build environment
# docker build --no-cache=true --rm -t xrlx/heka_scamp_build ./docker
docker build --no-cache=true --rm -t xrlx/heka_scamp_build ./docker

# copy the .deb in to a fresh build and tag for release
docker run --rm -v /var/run/docker.sock:/var/run/docker.sock -ti xrlx/heka_scamp_build

docker tag -f xrlx/heka_scamp gcr.io/retailops-1/heka_scamp
docker tag -f gcr.io/retailops-1/heka_scamp:latest gcr.io/retailops-1/heka_scamp:test-build

if [[ "x$1" == "x--push" ]]; then
  gcloud docker push gcr.io/retailops-1/heka_scamp
  gcloud docker push gcr.io/retailops-1/heka_scamp:test-build
fi