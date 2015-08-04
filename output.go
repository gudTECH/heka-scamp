package heka_scamp

import "fmt"
import "github.com/gudtech/scamp-go/scamp"
import "github.com/mozilla-services/heka/pipeline"

type SCAMPOutputPluginConfig struct {
	Service string `toml:"service"`
}

type SCAMPOutputPlugin struct {
	conf *SCAMPOutputPluginConfig
	conn *scamp.Connection
}

func (sop *SCAMPOutputPlugin) ConfigStruct() interface{} {
		fmt.Println("ConfigStruct")
	return &SCAMPOutputPluginConfig {
		Service: ":30101",
	}
}

func (sop *SCAMPOutputPlugin) Init(config interface{}) (err error) {
	sop.conf = config.(*SCAMPOutputPluginConfig)

	fmt.Printf("Connecting to %s\n", sop.conf.Service)
	sop.conn,err = scamp.Connect(sop.conf.Service)
	if err != nil {
		return
	}

	return
}

func (sop *SCAMPOutputPlugin) Run(or pipeline.OutputRunner, h pipeline.PluginHelper) (err error){
	fmt.Errorf("running")
	for pack := range or.InChan() {
		fmt.Errorf("got a payload")
		payload := pack.Message.Payload
		fmt.Printf("payload: %s", payload)
	}
	return
}

func (sop *SCAMPOutputPlugin) CleanUp() {
	return
}

func init(){
	pipeline.RegisterPlugin("ScampOutput", func() interface{} {
		fmt.Println("RegisterPlugin/SCAMP")
		return new(SCAMPOutputPlugin)
	})
}