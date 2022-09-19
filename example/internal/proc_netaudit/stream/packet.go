package stream

import (
	"bytes"
	"github.com/funbinary/go_example/example/internal/proc_netaudit/decoder"
	log "github.com/funbinary/go_example/pkg/blog"
	"github.com/google/gopacket/layers"
	"io"
)

type Packet struct {
	proto  decoder.ApplicationPortocol
	t      *layers.TCP
	len    int
	buf    bytes.Buffer
	offset int64
}

func NewPacket(proto decoder.ApplicationPortocol, t *layers.TCP) *Packet {
	return &Packet{
		proto:  proto,
		t:      t,
		len:    0,
		offset: 0,
	}
}

func (p *Packet) Append(data []byte) int {
	r := bytes.NewReader(data)
	if p.len == 0 {
		io.Copy(&p.buf, r)
		d := decoder.GetDecoder(p.proto)
		pkgLen, off, err := d.DetectPacketLength(p.buf.Bytes(), p.t)
		log.Debugf("parse pkg len:%d", pkgLen)
		if err != nil {
			log.Errorf("detect packet length %+v ", err)
		}
		p.len = int(pkgLen)
		p.offset = int64(off)
		return len(data)
	}

	remain := p.len - p.buf.Len() + int(p.offset)
	if remain >= len(data) {
		// 全部拷贝
		remain = len(data)
	}
	// 拷贝remain
	io.CopyN(&p.buf, r, int64(remain))

	return remain

}

func (p *Packet) Completed() bool {
	log.Debugf("len:%d p.buf:%d", p.len, p.buf.Len())
	return p.len > 0 && p.len == (p.buf.Len()-int(p.offset))

}

func (p *Packet) Data() []byte {
	return p.buf.Bytes()
}
