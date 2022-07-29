package main

import (
	"github.com/funbinary/go_example/pkg/bfile"
	"log"
	"text/template"
)

var temp = `
package layer
var LayerType{{.Name}}  = gopacket.RegisterLayerType({{.Value}}, gopacket.LayerTypeMetadata{Name: "{{.Name}} ", Decoder: gopacket.DecodeFunc(decode{{.Name}})})

func decode{{.Name}}(data []byte, p gopacket.PacketBuilder) error {
	t := &{{.Name}}{}
	err := t.DecodeFromBytes(data, p)
	if err != nil {
		return err
	}
	p.AddLayer(t)
	p.SetApplicationLayer(t)
	return nil
}


type {{.Name}} struct {
	layers.BaseLayer
}

func (self *{{.Name}} ) LayerType() gopacket.LayerType { return LayerType{{.Name}}  }

//@brief 解码逻辑
//@
func (self *{{.Name}} ) DecodeFromBytes(data []byte, df gopacket.DecodeFeedback) error {
	self.BaseLayer.Contents = data
	self.BaseLayer.Payload = nil

	return nil
}

func (self *{{.Name}}) Payload() []byte {
	return nil
}

func (self *{{.Name}}) SerializeTo(b gopacket.SerializeBuffer, opts gopacket.SerializeOptions) error {
	return nil
}
`

type Layer struct {
	Name  string
	Value int
}

func main() {
	l := Layer{
		Name:  "Component",
		Value: 115200,
	}
	t := template.New("test")
	t = template.Must(t.Parse(temp))
	f, err := bfile.OpenFile("test.go", bfile.O_RDWR|bfile.O_CREATE, 0666)
	if err != nil {
		log.Panicln(err)
	}
	t.Execute(f, l)
}
