package sniffer

import (
	"fmt"
	"github.com/funbinary/go_example/pkg/bfile"
	log "github.com/funbinary/go_example/pkg/blog"
	"github.com/funbinary/go_example/pkg/errors"
	"github.com/google/gopacket"
	"github.com/google/gopacket/ip4defrag"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/google/gopacket/pcapgo"
	"go.uber.org/atomic"
	"io"
	"net"
	"runtime"
	"strings"
	"syscall"
	"time"
)

type Sniffer struct {
	file       string
	device     string
	bpfFilter  string
	dumpFile   string
	timeout    time.Duration
	snaplen    int
	promisc    bool
	bufferSize int
	state      atomic.Uint32
	writer     *pcapgo.Writer
	defrag     *ip4defrag.IPv4Defragmenter
}

type snifferHandle interface {
	gopacket.PacketDataSource
	LinkType() layers.LinkType
	Close()
}

// sniffer state values
const (
	snifferInactive = 0
	snifferClosing  = 1
	snifferActive   = 2
)

func New(opts ...Option) (*Sniffer, error) {
	s := &Sniffer{
		timeout:    -1 * time.Second,
		snaplen:    65536,
		promisc:    true,
		bufferSize: 1024 * 1024 * 100,
	}
	for _, opt := range opts {
		opt(s)
	}

	if s.file != "" {
		if s.bpfFilter != "" {
		}
		s.device = ""
	} else {
		if name, err := resolveDeviceName(s.device); err != nil {
			return nil, err
		} else {
			s.device = name
			if name == "any" && !deviceAnySupported {
				return nil, fmt.Errorf("any interface is not supported on %s", runtime.GOOS)
			}

		}
	}

	err := s.validation()
	if err != nil {
		return nil, err
	}

	return s, nil
}

func (s *Sniffer) Run() error {
	handle, err := s.open()
	if err != nil {
		return errors.Wrap(err, "failed to open sniffer ")
	}
	defer handle.Close()

	if s.dumpFile != "" {
		f, err := bfile.Create(s.dumpFile)
		if err != nil {
			return err
		}
		defer f.Close()

		s.writer = pcapgo.NewWriter(f)
		err = s.writer.WriteFileHeader(65535, handle.LinkType())
		if err != nil {
			return errors.Wrapf(err, "failed to write dump file %s header", s.dumpFile)
		}
	}

	if !s.state.CAS(snifferInactive, snifferActive) {
		return nil
	}
	defer s.state.Store(snifferInactive)
	//if err := decoder.Init(handle.LinkType()); err != nil {
	//	return errors.Wrap(err, "init decoder fail ")
	//}

	//source := gopacket.NewPacketSource(handle, handle.LinkType())
	s.defrag = ip4defrag.NewIPv4Defragmenter()

	//var packet gopacket.Packet
	for s.state.Load() == snifferActive {

		_, ci, err := handle.ReadPacketData()
		fmt.Println("now:", time.Now().Format("2006-01-02 15:04::05.0000"), "capture:", ci.Timestamp.Format("2006-01-02 15:04::05.0000"))
		if err == nil {
			//err = s.handlePacket(packet)
			//if err != nil {
			//	return errors.Wrapf(err, "handlePacket error")
			//}
			continue
		}

		// Immediately retry for temporary network errors
		if nerr, ok := err.(net.Error); ok && nerr.Temporary() {
			continue
		}

		// Immediately retry for EAGAIN
		if err == syscall.EAGAIN {
			continue
		}

		// Immediately break for known unrecoverable errors
		if err == io.EOF || err == io.ErrUnexpectedEOF ||
			err == io.ErrNoProgress || err == io.ErrClosedPipe || err == io.ErrShortBuffer ||
			err == syscall.EBADF ||
			strings.Contains(err.Error(), "use of closed file") {
			break
		}

		// Sleep briefly and try again
		//time.Sleep(time.Millisecond * time.Duration(5))

	}

	return nil

}

func (s *Sniffer) handlePacket(packet gopacket.Packet) error {

	if s.writer != nil {
		err := s.writer.WritePacket(packet.Metadata().CaptureInfo, packet.Data())
		if err != nil {
			return errors.Errorf("failed to write packet: %v", err)
		}
	}
	if packet == nil {
		return errors.New("packet is nil")
	}
	// 拼接IP层
	ip4Layer := packet.Layer(layers.LayerTypeIPv4)
	if ip4Layer == nil {
		return nil
	}
	ip4 := ip4Layer.(*layers.IPv4)
	l := ip4.Length
	newip4, err := s.defrag.DefragIPv4(ip4)
	if err != nil {
		log.Errorf("Error while de-fragmenting: %+v", err)
	} else if newip4 == nil {
		return nil
	}
	if newip4.Length != l {
		pb, ok := packet.(gopacket.PacketBuilder)
		if !ok {
			log.Errorf("Not a PacketBuilder")
		}
		nextDecoder := newip4.NextLayerType()
		nextDecoder.Decode(newip4.Payload, pb)
	}

	return nil
}

func (s *Sniffer) open() (snifferHandle, error) {
	if s.file != "" {
		return newFileHandler(s.file, true, 0)
	} else {
		return s.openPcap()
	}
}

func (s *Sniffer) openPcap() (snifferHandle, error) {
	ih, err := pcap.NewInactiveHandle(s.device)
	if err != nil {
		return nil, errors.Wrapf(err, "NewInactiveHandle %s failed", s.device)
	}
	ih.SetPromisc(s.promisc)
	ih.SetBufferSize(s.bufferSize)
	ih.SetSnapLen(s.snaplen)
	ih.SetTimeout(s.timeout)
	ih.SetImmediateMode(true)
	handler, err := ih.Activate()
	if err != nil {
		return nil, errors.Wrapf(err, "activate %s error", s.device)

	}

	err = handler.SetBPFFilter(s.bpfFilter)
	if err != nil {
		handler.Close()
		return nil, errors.Wrapf(err, "set BPF filter %s failed", s.bpfFilter)
	}
	log.Infof("open device %s bpfFilter:%s ", s.device, s.bpfFilter)

	return handler, nil
}

func (s *Sniffer) validation() error {
	if s.bpfFilter == "" {
		return nil
	}
	_, err := pcap.NewBPF(layers.LinkTypeEthernet, 65535, s.bpfFilter)
	return err
}
