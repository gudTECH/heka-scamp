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

	announcer,err := scamp.NewDiscoveryAnnouncer()
	if err != nil {
		scamp.Error.Printf("failed to create announcer: `%s`", err)
		return
	}
	announcer.Track(sop.service)
	go announcer.AnnounceLoop()

	var handlerConfig SCAMPInputHandlerConfig
	for _,handlerConfig = range sop.conf.Handlers {
		scamp.Trace.Printf("registering handler: `%s`", handlerConfig)

		sop.service.Register(handlerConfig.Action, func(req scamp.Request, sess *scamp.Session) {
			var pack *pipeline.PipelinePack


			pack = <-ir.InChan()
			pack.Message.SetUuid(uuid.NewRandom())
			pack.Message.SetTimestamp(time.Now().UnixNano())
			pack.Message.SetType(handlerConfig.Type)
			pack.Message.SetPayload(string(req.Blob[:]))
			pack.Message.SetSeverity(int32(handlerConfig.Severity))
			pack.Message.SetLogger("heka-scamp") // TODO not sure what this means
			ir.Deliver(pack)

			err = sess.Send(scamp.Reply{Blob: []byte("{}")})
			if err != nil {
				return
			}
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