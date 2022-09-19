package nfs

import (
	"bytes"
	"github.com/funbinary/go_example/example/internal/proc_netaudit/decoder/nfs/nfsv2"
	"github.com/funbinary/go_example/example/internal/proc_netaudit/decoder/nfs/nfsv3"
	"github.com/funbinary/go_example/example/internal/proc_netaudit/hook"
	log "github.com/funbinary/go_example/pkg/blog"
	"github.com/funbinary/go_example/pkg/errors"
	"github.com/funbinary/go_example/pkg/rpc"
	"github.com/funbinary/go_example/pkg/xdr"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"io"
	"sync"
)

type Decoder struct {
	name  string
	cache sync.Map //xid
}

func NewDecoder() *Decoder {
	return &Decoder{
		name: "nfs",
	}
}

func (self *Decoder) Name() string { return self.name }

func (self *Decoder) CanDecode(data []byte, t gopacket.TransportLayer) (canDecode bool) {
	reader := bytes.NewReader(data)

	return self.parseRpc(reader, t) != nil
}

func (self *Decoder) DetectPacketLength(data []byte, t gopacket.TransportLayer) (length uint64, off uint64, err error) {
	if !self.CanDecode(data, t) {
		return 0, 0, errors.New("can'not decode")
	}
	reader := bytes.NewReader(data)
	var reqLen uint32
	if t.LayerType() == layers.LayerTypeTCP {
		fragment, err := xdr.ReadUint32(reader)
		if err != nil {
			return 0, 0, errors.Wrap(err, "parse fragment error")
		}
		if fragment&(1<<31) == 0 {
			return 0, 0, errors.New("fragment dont for zero")
		}
		reqLen = fragment - uint32(1<<31)

	} else if t.LayerType() == layers.LayerTypeUDP {
		udp := t.(*layers.UDP)
		reqLen = uint32(udp.Length)
	} else {
		return 0, 0, errors.New("unknown transport type")
	}
	return uint64(reqLen), 4, nil
}

func (self *Decoder) Decode(data []byte, t gopacket.TransportLayer) (hook.Protocol, error) {
	r := bytes.NewReader(data)
	header := self.parseRpc(r, t)
	if header == nil {
		return header, errors.New("解析RPC协议失败")
	}
	log.Infof("decoder nfs")
	switch header.Program {
	case rpc.PortmapID:
		return header, nil
	case rpc.NfsServiceID:
		return header, self.parseNfs(header, r)
	case rpc.MountServiceID:
		return header, nil
	}
	return header, errors.New("Program Type当前不支持解析")
}

func (self *Decoder) String() string { return self.Name() }

func (self *Decoder) parseRpc(r io.Reader, t gopacket.TransportLayer) *Nfs {
	// 解析长度
	var reqLen uint32
	if t.LayerType() == layers.LayerTypeTCP {
		fragment, err := xdr.ReadUint32(r)
		if err != nil {
			return nil
		}
		if fragment&(1<<31) == 0 {
			return nil
		}
		reqLen = fragment - uint32(1<<31)

	} else if t.LayerType() == layers.LayerTypeUDP {
		udp := t.(*layers.UDP)
		reqLen = uint32(udp.Length)
	} else {
		return nil
	}

	if reqLen > 0xffffff {
		return nil
	}

	var err error
	lr := &io.LimitedReader{R: r, N: int64(reqLen)}
	xid, err := xdr.ReadUint32(lr)
	if err != nil {
		return nil
	}
	var nh *Nfs
	if v, ok := self.cache.Load(xid); ok {
		nh = v.(*Nfs)
	} else {
		nh = NewNfsHeader()
		nh.Xid = xid
		self.cache.Store(xid, nh)
	}

	msgType, err := xdr.ReadUint32(lr)
	if err != nil {
		return nil
	}
	nh.CurMsgType = rpc.MessageType(msgType)
	switch nh.CurMsgType {
	case rpc.RpcCall:
		if err = xdr.Read(lr, &nh.ReqHeader); err != nil {
			return nil
		}
	case rpc.RpcReply:
		if err = xdr.Read(lr, &nh.RepHeader); err != nil {
			return nil
		}
	default:
		return nil
	}

	return nh
}

func (self *Decoder) parseNfs(h *Nfs, r io.Reader) (err error) {
	switch rpc.NfsVersion(h.ProgVers) {
	case rpc.NfsV2:
		return self.parseNfsV2(h, r)
	case rpc.NfsV3:
		return self.parseNfsV3(h, r)
	default:
		return errors.Errorf("不支持解析Nfs版本:%d", h.ProgVers)
	}
}
func (self *Decoder) parseNfsV2(h *Nfs, r io.Reader) (err error) {
	log.Infof("decoder nfsv2 xid:%x dir:%v produce:%v", h.Xid, h.CurMsgType, nfsv2.Procedure(h.Procedure))
	switch nfsv2.Procedure(h.Procedure) {
	case nfsv2.NFS2ProcedureNull:
		return nil
	case nfsv2.NFS2ProcedureGetAttr:
		handlerNfsv2GetAttr(h, r)
	case nfsv2.NFS2ProcedureSetAttr:
		handlerNfsv2SetAttr(h, r)
	case nfsv2.NFS2ProcedureRoot:
		return nil
	case nfsv2.NFS2ProcedureLookUp:
		handlerNfsv2LookUp(h, r)
	case nfsv2.NFS2ProcedureReadlink:
		handlerNfsv2ReadLink(h, r)
	case nfsv2.NFS2ProcedureRead:
		handlerNfsv2Read(h, r)
	case nfsv2.NFS2ProcedureWriteCache:
		return nil
	case nfsv2.NFS2ProcedureWrite:
		handlerNfsv2Write(h, r)
	case nfsv2.NFS2ProcedureCreate:
		handlerNfsv2Create(h, r)
	case nfsv2.NFS2ProcedureRemove:
		handlerNfsv2Remove(h, r)
	case nfsv2.NFS2ProcedureRename:
		handlerNfsv2Rename(h, r)
	case nfsv2.NFS2ProcedureLink:
		handlerNfsv2Link(h, r)
	case nfsv2.NFS2ProcedureSymlink:
		handlerNfsv2SymLink(h, r)
	case nfsv2.NFS2ProcedureMkDir:
		handlerNfsv2Mkdir(h, r)
	case nfsv2.NFS2ProcedureRmDir:
		handlerNfsv2Rmdir(h, r)
	case nfsv2.NFS2ProcedureReadDir:
		handlerNfsv2ReadDir(h, r)
	case nfsv2.NFS2ProcedureStatFs:
		handlerNfsv2StatFs(h, r)
	}
	return nil
}
func (self *Decoder) parseNfsV3(h *Nfs, r io.Reader) (err error) {
	log.Infof("decoder nfsv3 xid:%x dir:%v produce:%v", h.Xid, h.CurMsgType, nfsv3.Procedure(h.Procedure))
	switch nfsv3.Procedure(h.Procedure) {
	case nfsv3.NFS3ProcedureNull:
		return nil
	case nfsv3.NFS3ProcedureGetAttr:
		return handlerGetAttr(h, r)
	case nfsv3.NFS3ProcedureSetAttr:
		return handleSetAttr(h, r)
	case nfsv3.NFS3ProcedureLookup:
		return handleLookUp(h, r)
	case nfsv3.NFS3ProcedureAccess:
		return handleAccess(h, r)
	case nfsv3.NFS3ProcedureReadlink:
		return handleReadLink(h, r)
	case nfsv3.NFS3ProcedureRead:
		return handleRead(h, r)
	case nfsv3.NFS3ProcedureWrite:
		return handleWrite(h, r)
	case nfsv3.NFS3ProcedureCreate:
		return handleCreate(h, r)
	case nfsv3.NFS3ProcedureMkDir:
		return handleMkdir(h, r)
	case nfsv3.NFS3ProcedureSymlink:
		return handleSymLink(h, r)
	case nfsv3.NFS3ProcedureMkNod:
		return handleMknod(h, r)
	case nfsv3.NFS3ProcedureRemove:
		return handleRemove(h, r)
	case nfsv3.NFS3ProcedureRmDir:
		return handleRmdir(h, r)
	case nfsv3.NFS3ProcedureRename:
		return handleRename(h, r)
	case nfsv3.NFS3ProcedureLink:
		return handleReadLink(h, r)
	case nfsv3.NFS3ProcedureReadDir:
		return handleReadDir(h, r)
	case nfsv3.NFS3ProcedureReadDirPlus:
		return handleReadDirPlus(h, r)
	case nfsv3.NFS3ProcedureFSStat:
		return handleFsStat(h, r)
	case nfsv3.NFS3ProcedureFSInfo:
		return handleFsInfo(h, r)
	case nfsv3.NFS3ProcedurePathConf:
		return handlePathConf(h, r)
	case nfsv3.NFS3ProcedureCommit:
		return handleCommit(h, r)
	default:
		return nil
		return errors.Errorf("不支持解析NFS produce:%d", h.Procedure)
	}
}
