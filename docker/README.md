 - [ ] Create `heka_base` image with latest code, push to registry (right now we're using `xrlx/heka_base` without any real versioning)
   - `git clone https://github.com/mozilla-services/heka`
   - `cd heka && docker build --rm -t xrlx/heka_base`
   - `docker push xrlx/heka_base`
 - [ ] Run `./build.sh` to prepare image
   - modify `hekad.toml` however you see fit
   - `./build.sh`
   - `docker push xrlx/heka_scamp`

If you want direct access to the goodies in `xrlx/heka_scamp` you should override the entrypoint on run:

    docker run -it --rm --entrypoint=/bin/bash xrlx/heka_scamp
    apt-get install -y vim
    vim /etc/heka/conf.d/00-hekad.toml
