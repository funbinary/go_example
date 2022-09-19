package nfs

import (
	"github.com/funbinary/go_example/example/internal/proc_netaudit/decoder/nfs/nfsv2"
	log "github.com/funbinary/go_example/pkg/blog"
	"github.com/funbinary/go_example/pkg/errors"
	"github.com/funbinary/go_example/pkg/rpc"
	"github.com/funbinary/go_example/pkg/xdr"
	"io"
)

func handlerNfsv2GetAttr(h *Nfs, r io.Reader) (err error) {
	dir := h.CurMsgType
	if dir == rpc.RpcCall {
		call := &nfsv2.GetAttrCall{}
		h.Call = call
		err = xdr.Read(r, call)
		if err != nil {
			return errors.Wrap(err, "handlerNfsv2GetAttrCall parse error")
		}
		log.Infof("GetAttrCall:%v", h.Call)
	} else if dir == rpc.RpcReply {
		reply := &nfsv2.GetAttrReply{}
		h.Reply = reply

		err = xdr.Read(r, reply)
		if err != nil {
			return errors.Wrap(err, "handlerNfsv2GetAttrReply parse failed")
		}

		log.Infof("GetAttrReply:%v", h.Reply)
	} else {
		return errors.Errorf("解析Message Type:%d失败", dir)
	}
	return nil
}

func handlerNfsv2SetAttr(h *Nfs, r io.Reader) (err error) {
	dir := h.CurMsgType
	if dir == rpc.RpcCall {
		call := &nfsv2.SetAttrCall{}
		h.Call = call
		err = xdr.Read(r, call)
		if err != nil {
			return errors.Wrap(err, "handlerNfsv2SetAttrCall parse error")
		}
		log.Infof("SetAttrCall:%v", h.Call)
	} else if dir == rpc.RpcReply {
		reply := &nfsv2.SetAttrReply{}
		h.Reply = reply

		err = xdr.Read(r, reply)
		if err != nil {
			return errors.Wrap(err, "handlerNfsv2SetAttrReply parse failed")
		}

		log.Infof("SetAttrReply:%v", h.Reply)
	} else {
		return errors.Errorf("解析Message Type:%d失败", dir)
	}
	return nil
}

func handlerNfsv2LookUp(h *Nfs, r io.Reader) (err error) {
	dir := h.CurMsgType
	if dir == rpc.RpcCall {
		call := &nfsv2.LookUpCall{}
		h.Call = call
		err = xdr.Read(r, call)
		if err != nil {
			return errors.Wrap(err, "handlerNfsv2LookUpCall parse error")
		}
		log.Infof("LookUpCall:%v", h.Call)
	} else if dir == rpc.RpcReply {
		reply := &nfsv2.LookUpReply{}
		h.Reply = reply

		err = xdr.Read(r, reply)
		if err != nil {
			return errors.Wrap(err, "handlerNfsv2LookUpReply parse failed")
		}

		log.Infof("LookUpReply:%v", h.Reply)
	} else {
		return errors.Errorf("解析Message Type:%d失败", dir)
	}
	return nil
}

func handlerNfsv2ReadLink(h *Nfs, r io.Reader) (err error) {
	dir := h.CurMsgType
	if dir == rpc.RpcCall {
		call := &nfsv2.ReadLinkCall{}
		h.Call = call
		err = xdr.Read(r, call)
		if err != nil {
			return errors.Wrap(err, "handlerNfsv2ReadLinkCall parse error")
		}
		log.Infof("ReadLinkCall:%v", h.Call)
	} else if dir == rpc.RpcReply {
		reply := &nfsv2.ReadLinkReply{}
		h.Reply = reply

		err = xdr.Read(r, reply)
		if err != nil {
			return errors.Wrap(err, "handlerNfsv2ReadLinkReply parse failed")
		}

		log.Infof("ReadLinkReply:%v", h.Reply)
	} else {
		return errors.Errorf("解析Message Type:%d失败", dir)
	}
	return nil
}

func handlerNfsv2Read(h *Nfs, r io.Reader) (err error) {
	dir := h.CurMsgType
	if dir == rpc.RpcCall {
		call := &nfsv2.ReadCall{}
		h.Call = call
		err = xdr.Read(r, call)
		if err != nil {
			return errors.Wrap(err, "handlerNfsv2ReadCall parse error")
		}
		log.Infof("ReadCall:%v", h.Call)
	} else if dir == rpc.RpcReply {
		reply := &nfsv2.ReadReply{}
		h.Reply = reply

		err = xdr.Read(r, reply)
		if err != nil {
			return errors.Wrap(err, "handlerNfsv2ReadReply parse failed")
		}

		log.Infof("ReadReply:%v", h.Reply)
	} else {
		return errors.Errorf("解析Message Type:%d失败", dir)
	}
	return nil
}

func handlerNfsv2Write(h *Nfs, r io.Reader) (err error) {
	dir := h.CurMsgType
	if dir == rpc.RpcCall {
		call := &nfsv2.WriteCall{}
		h.Call = call
		err = xdr.Read(r, call)
		if err != nil {
			return errors.Wrap(err, "handlerNfsv2WriteCall parse error")
		}
		log.Infof("WriteCall:%v", h.Call)
	} else if dir == rpc.RpcReply {
		reply := &nfsv2.WriteReply{}
		h.Reply = reply

		err = xdr.Read(r, reply)
		if err != nil {
			return errors.Wrap(err, "handlerNfsv2WriteReply parse failed")
		}

		log.Infof("WriteReply:%v", h.Reply)
	} else {
		return errors.Errorf("解析Message Type:%d失败", dir)
	}
	return nil
}

func handlerNfsv2Create(h *Nfs, r io.Reader) (err error) {
	dir := h.CurMsgType
	if dir == rpc.RpcCall {
		call := &nfsv2.CreateCall{}
		h.Call = call
		err = xdr.Read(r, call)
		if err != nil {
			return errors.Wrap(err, "handlerNfsv2CreateCall parse error")
		}
		log.Infof("CreateCall:%v", h.Call)
	} else if dir == rpc.RpcReply {
		reply := &nfsv2.CreateReply{}
		h.Reply = reply

		err = xdr.Read(r, reply)
		if err != nil {
			return errors.Wrap(err, "handlerNfsv2CreateReply parse failed")
		}

		log.Infof("CreateReply:%v", h.Reply)
	} else {
		return errors.Errorf("解析Message Type:%d失败", dir)
	}
	return nil
}

func handlerNfsv2Remove(h *Nfs, r io.Reader) (err error) {
	dir := h.CurMsgType
	if dir == rpc.RpcCall {
		call := &nfsv2.RemoveCall{}
		h.Call = call
		err = xdr.Read(r, call)
		log.Errorf("RemoveCall parse success")

		if err != nil {
			log.Errorf("RemoveCall error,%+v", err)
			return errors.Wrap(err, "handlerNfsv2RemoveCall parse error")
		}
		log.Infof("RemoveCall:%v", h.Call)
	} else if dir == rpc.RpcReply {
		reply := &nfsv2.RemoveReply{}
		h.Reply = reply

		err = xdr.Read(r, reply)
		if err != nil {
			return errors.Wrap(err, "handlerNfsv2RemoveReply parse failed")
		}

		log.Infof("RemoveReply:%v", h.Reply)
	} else {
		return errors.Errorf("解析Message Type:%d失败", dir)
	}
	return nil
}

func handlerNfsv2Rename(h *Nfs, r io.Reader) (err error) {
	dir := h.CurMsgType
	if dir == rpc.RpcCall {
		call := &nfsv2.RenameCall{}
		h.Call = call
		err = xdr.Read(r, call)
		if err != nil {
			return errors.Wrap(err, "handlerNfsv2RenameCall parse error")
		}
		log.Infof("RenameCall:%v", h.Call)
	} else if dir == rpc.RpcReply {
		reply := &nfsv2.RenameReply{}
		h.Reply = reply

		err = xdr.Read(r, reply)
		if err != nil {
			return errors.Wrap(err, "handlerNfsv2RenameReply parse failed")
		}

		log.Infof("RenameReply:%v", h.Reply)
	} else {
		return errors.Errorf("解析Message Type:%d失败", dir)
	}
	return nil
}

func handlerNfsv2Link(h *Nfs, r io.Reader) (err error) {
	dir := h.CurMsgType
	if dir == rpc.RpcCall {
		call := &nfsv2.LinkCall{}
		h.Call = call
		err = xdr.Read(r, call)
		if err != nil {
			return errors.Wrap(err, "handlerNfsv2LinkCall parse error")
		}
		log.Infof("LinkCall:%v", h.Call)
	} else if dir == rpc.RpcReply {
		reply := &nfsv2.LinkReply{}
		h.Reply = reply

		err = xdr.Read(r, reply)
		if err != nil {
			return errors.Wrap(err, "handlerNfsv2LinkReply parse failed")
		}

		log.Infof("LinkReply:%v", h.Reply)
	} else {
		return errors.Errorf("解析Message Type:%d失败", dir)
	}
	return nil
}

func handlerNfsv2SymLink(h *Nfs, r io.Reader) (err error) {
	dir := h.CurMsgType
	if dir == rpc.RpcCall {
		call := &nfsv2.SymlinkCall{}
		h.Call = call
		err = xdr.Read(r, call)
		if err != nil {
			return errors.Wrap(err, "handlerNfsv2SymLinkCall parse error")
		}
		log.Infof("SymLinkCall:%v", h.Call)
	} else if dir == rpc.RpcReply {
		reply := &nfsv2.SymlinkReply{}
		h.Reply = reply

		err = xdr.Read(r, reply)
		if err != nil {
			return errors.Wrap(err, "handlerNfsv2SymLinkReply parse failed")
		}

		log.Infof("SymLinkReply:%v", h.Reply)
	} else {
		return errors.Errorf("解析Message Type:%d失败", dir)
	}
	return nil
}

func handlerNfsv2Mkdir(h *Nfs, r io.Reader) (err error) {
	dir := h.CurMsgType
	if dir == rpc.RpcCall {
		call := &nfsv2.MkdirCall{}
		h.Call = call
		err = xdr.Read(r, call)
		if err != nil {
			return errors.Wrap(err, "handlerNfsv2MkdirCall parse error")
		}
		log.Infof("MkdirCall:%v", h.Call)
	} else if dir == rpc.RpcReply {
		reply := &nfsv2.MkdirReply{}
		h.Reply = reply

		err = xdr.Read(r, reply)
		if err != nil {
			return errors.Wrap(err, "handlerNfsv2MkdirReply parse failed")
		}

		log.Infof("MkdirReply:%v", h.Reply)
	} else {
		return errors.Errorf("解析Message Type:%d失败", dir)
	}
	return nil
}

func handlerNfsv2Rmdir(h *Nfs, r io.Reader) (err error) {
	dir := h.CurMsgType
	if dir == rpc.RpcCall {
		call := &nfsv2.RmdirCall{}
		h.Call = call
		err = xdr.Read(r, call)
		if err != nil {
			return errors.Wrap(err, "handlerNfsv2RmdirCall parse error")
		}
		log.Infof("RmdirCall:%v", h.Call)
	} else if dir == rpc.RpcReply {
		reply := &nfsv2.RmdirReply{}
		h.Reply = reply

		err = xdr.Read(r, reply)
		if err != nil {
			return errors.Wrap(err, "handlerNfsv2RmdirReply parse failed")
		}

		log.Infof("RmdirReply:%v", h.Reply)
	} else {
		return errors.Errorf("解析Message Type:%d失败", dir)
	}
	return nil
}

func handlerNfsv2ReadDir(h *Nfs, r io.Reader) (err error) {
	dir := h.CurMsgType
	if dir == rpc.RpcCall {
		call := &nfsv2.ReadDirCall{}
		h.Call = call
		err = xdr.Read(r, call)
		if err != nil {
			return errors.Wrap(err, "handlerNfsv2ReadDirCall parse error")
		}
		log.Infof("ReadDirCall:%v", h.Call)
	} else if dir == rpc.RpcReply {
		reply := &nfsv2.ReadDirReply{}
		h.Reply = reply

		err = xdr.Read(r, reply)
		if err != nil {
			return errors.Wrap(err, "handlerNfsv2ReadDirReply parse failed")
		}

		log.Infof("ReadDirReply:%v", h.Reply)
	} else {
		return errors.Errorf("解析Message Type:%d失败", dir)
	}
	return nil
}

func handlerNfsv2StatFs(h *Nfs, r io.Reader) (err error) {
	dir := h.CurMsgType
	if dir == rpc.RpcCall {
		call := &nfsv2.StatFsCall{}
		h.Call = call
		err = xdr.Read(r, call)
		if err != nil {
			return errors.Wrap(err, "handlerNfsv2StatFsCall parse error")
		}
		log.Infof("StatFsCall:%v", h.Call)
	} else if dir == rpc.RpcReply {
		reply := &nfsv2.StatFsReply{}
		h.Reply = reply

		err = xdr.Read(r, reply)
		if err != nil {
			return errors.Wrap(err, "handlerNfsv2StatFsReply parse failed")
		}

		log.Infof("StatFsReply:%v", h.Reply)
	} else {
		return errors.Errorf("解析Message Type:%d失败", dir)
	}
	return nil
}
