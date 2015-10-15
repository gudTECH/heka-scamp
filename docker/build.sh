#!/usr/bin/env sh

set -xe

# assume heka and heka-scamp are siblings in a directory
docker build --no-cache=false --rm -t xrlx/heka_base ../heka

# prepare the scamp plugin using our heka_base build environment
docker build --no-cache=true --rm -t xrlx/heka_scamp_build ./docker

# run the build env CMD for copying deb and installing to minimal image
docker run --rm -v /var/run/docker.sock:/var/run/docker.sock -ti xrlx/heka_scamp_build

# clean up
# docker rmi xrlx/heka_scamp_build

docker tag -f xrlx/heka_scamp gcr.io/retailops-1/heka_scamp
