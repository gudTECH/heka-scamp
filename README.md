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

## Sample ScampOutput config

```
[TailTestLog]
type = "LogstreamerInput"
log_directory = "/var/log"
file_match = 'authd\.log'
decoder = "SimpleDecoder"
	[SimpleDecoder]
	type = "SandboxDecoder"
	filename = "simple_decoder.lua"

[ScampOutput]
encoder = "PayloadEncoder"
action = "hello.helloworld"
message_matcher = "TRUE"
```

## Sample ScampInput config

This will reuse an instance of `scamp.Service` and register handlers which will feed in to the pipeline

```
[ScampInput]
[[ScampInput.handler]]
action = "sup.dude"
```