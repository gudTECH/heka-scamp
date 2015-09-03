#!/usr/bin/env sh

set -xe

# prepare the scamp plugin using our heka_base build environment
docker build --rm -t xrlx/heka_scamp_build .

# run the build env CMD for copying deb and installing to minimal image
docker run --rm -v /var/run/docker.sock:/var/run/docker.sock -ti xrlx/heka_scamp_build

# clean up
docker rmi xrlx/heka_scamp_build
