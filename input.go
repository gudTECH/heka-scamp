package heka_scamp

import "errors"

import "github.com/gudtech/scamp-go/scamp"
import "github.com/mozilla-services/heka/pipeline"

type SCAMPInputPluginConfig struct {
	Service string `toml:"listen"`
	Name string `toml:"name"`
	Handlers map[string]SCAMPInputHandlerConfig `toml:"handler"`
}

type SCAMPInputHandlerConfig struct {
	Action string `toml:"action"`
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

			pack = <-ir.InChan()
			payload := string(req.Blob[:])
			scamp.Trace.Printf("payload: `%s`", payload)
			pack.Message.SetPayload(payload)
			ir.Deliver(pack)

			scamp.Trace.Printf("closing session")
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