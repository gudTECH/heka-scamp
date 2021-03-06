package heka_scamp

import (
	"fmt"
	"errors"
)
import "github.com/gudtech/scamp-go/scamp"
import "github.com/mozilla-services/heka/pipeline"

type SCAMPOutputPluginConfig struct {
	Service string `toml:"service"`
	Action string `toml:"action"`
}

type SCAMPOutputPlugin struct {
	conf *SCAMPOutputPluginConfig
	client *scamp.Client
}

func (sop *SCAMPOutputPlugin) ConfigStruct() interface{} {
	fmt.Println("ConfigStruct")
	return &SCAMPOutputPluginConfig {
		Service: ":30101",
		Action: "Test.test", // TODO no smart default for this
	}
}

func (sop *SCAMPOutputPlugin) Init(config interface{}) (err error) {
	scamp.Initialize(scamp.DefaultConfigPath)

	sop.conf = config.(*SCAMPOutputPluginConfig)

	scamp.Info.Printf( "Connecting to %s\n", sop.conf.Service )
	sop.client,err = scamp.Dial(sop.conf.Service)
	if err != nil {
		return
	}

	return
}

func (sop *SCAMPOutputPlugin) Run(or pipeline.OutputRunner, h pipeline.PluginHelper) (err error){
	var pack *pipeline.PipelinePack

	// We have no default encoder
	if or.Encoder() == nil {
		return errors.New("Encoder required.")
	}

	for pack = range or.InChan() {
		scamp.Info.Printf("received pipeline pack")
		encoded,err := or.Encode(pack) // pack.Message.GetPayload()

		if err == nil {
			scamp.Info.Printf("payload: %s", encoded)
			msg := scamp.NewMessage()
			msg.SetEnvelope(scamp.ENVELOPE_JSON)
			msg.SetAction(sop.conf.Action)
			msg.SetVersion(1)
			msg.Write(encoded)
		}

		pack.Recycle(err)
	}
	fmt.Println("sup from end of for loop in Run")
	return
}

func (sop *SCAMPOutputPlugin) CleanUp() {
	fmt.Println("sup from cleanup")
	return
}

func init(){
	pipeline.RegisterPlugin("ScampOutput", func() interface{} {
		return new(SCAMPOutputPlugin)
	})
}