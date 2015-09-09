package heka_scamp

import "errors"

import "github.com/gudtech/scamp-go/scamp"
import "github.com/mozilla-services/heka/pipeline"

import "time"
import "github.com/pborman/uuid"

type SCAMPInputPluginConfig struct {
	Service string `toml:"listen"`
	Name string `toml:"name"`
	Handlers map[string]SCAMPInputHandlerConfig `toml:"handler"`
}

type SCAMPInputHandlerConfig struct {
	Action string `toml:"action"`
	Type string `toml:"type"`
	Severity int `toml:"severity"`
}

type SCAMPInputPlugin struct {
	conf *SCAMPInputPluginConfig
	service *scamp.Service
}

func (sop *SCAMPInputPlugin) ConfigStruct() interface{} {
	return &SCAMPInputPluginConfig {
		Service: ":30101",
		Handlers: make(map[string]SCAMPInputHandlerConfig),
	}
}

func (sop *SCAMPInputPlugin) Init(config interface{}) (err error) {
	scamp.Initialize()
	sop.conf = config.(*SCAMPInputPluginConfig)

	if len(sop.conf.Handlers) < 1 {
		err = errors.New("must provide at least 1 handler")
		return
	}

	return
}

func (sop *SCAMPInputPlugin) Run(ir pipeline.InputRunner, h pipeline.PluginHelper) (err error) {
	sop.service,err = scamp.NewService(sop.conf.Service, sop.conf.Name)
	if err != nil {
		return
	}

	var handlerConfig SCAMPInputHandlerConfig
	for _,handlerConfig = range sop.conf.Handlers {
		scamp.Trace.Printf("registering handler: `%s`", handlerConfig)

		sop.service.Register(handlerConfig.Action, func(req scamp.Request, sess *scamp.Session) {
			var pack *pipeline.PipelinePack

			// Example code from the documentation: https://hekad.readthedocs.org/en/v0.9.2/developing/plugin.html
		    // pack := <-hi.ir.InChan()
		    // pack.Message.SetUuid(uuid.NewRandom())
		    // pack.Message.SetTimestamp(time.Now().UnixNano())
		    // pack.Message.SetType("heka.httpinput.error")
		    // pack.Message.SetPayload(err.Error())
		    // pack.Message.SetSeverity(hi.conf.ErrorSeverity)
		    // pack.Message.SetLogger(url)
		    // hi.ir.Deliver(pack)
		    // return

			pack = <-ir.InChan()
			pack.Message.SetUuid(uuid.NewRandom())
			pack.Message.SetTimestamp(time.Now().UnixNano())
			pack.Message.SetType(handlerConfig.Type)
			pack.Message.SetPayload(string(req.Blob[:]))
			pack.Message.SetSeverity(int32(handlerConfig.Severity))
			pack.Message.SetLogger("heka-scamp") // TODO not sure what this means
			ir.Deliver(pack)

			sess.CloseReply()
		})
	}

	sop.service.Run()
	return
}

func (sop *SCAMPInputPlugin) Stop() {
	sop.service.Stop()
	return
}

func init(){
	pipeline.RegisterPlugin("ScampInput", func() interface{} {
		return new(SCAMPInputPlugin)
	})
}