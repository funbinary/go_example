package decoder

import (
	"github.com/funbinary/go_example/example/internal/proc_netaudit/hook"
	log "github.com/funbinary/go_example/pkg/blog"
	"github.com/funbinary/go_example/pkg/errors"
	"github.com/google/gopacket"
)

type Decoder interface {
	Name() string
	CanDecode(data []byte, t gopacket.TransportLayer) (canDecode bool)
	DetectPacketLength(data []byte, t gopacket.TransportLayer) (length uint64, off uint64, err error)
	Decode(data []byte, t gopacket.TransportLayer) (app hook.Protocol, err error)
}

func GetDecoder(p ApplicationPortocol) Decoder {
	return decoderMap[p]
}

var idm = make(map[string]ApplicationPortocol)

func Decode(id string, data []byte, tl gopacket.TransportLayer) (err error) {
	var proto ApplicationPortocol
	var ok bool
	if proto, ok = idm[id]; !ok {
		proto = DetectProtocol(data, tl)
		idm[id] = proto
	}
	decoder := GetDecoder(proto)
	app, err := decoder.Decode(data, tl)
	if err != nil {
		return err
	}
	switch proto {
	case NfsProtocol:
		go hook.Fire(id, app)

	default:
		err = errors.Wrap(err, "Unknown protocol")
	}
	return err
}

func DecodeUseProto(id string, data []byte, tl gopacket.TransportLayer, proto ApplicationPortocol) (err error) {
	decoder := GetDecoder(proto)
	app, err := decoder.Decode(data, tl)
	log.Errorf("%+v", err)
	if err != nil {
		return err
	}
	switch proto {
	case NfsProtocol:
		go hook.Fire(id, app)

	default:
		err = errors.Wrap(err, "Unknown protocol")
	}
	return err
}
