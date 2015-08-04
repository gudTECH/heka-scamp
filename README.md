# heka-scamp

Input/output plugin for talking on a scamp network.

## Building heka with scamp support

```
git clone git@github.com:mozilla-services/heka.git
cd heka/
echo "add_external_plugin(git https://github.com/gudtech/heka-scamp master)" > "./cmake/plugin_loader.cmake"
go get github.com/gudTECH/scamp-go/scamp
. ./env.sh
sh build.sh
```
