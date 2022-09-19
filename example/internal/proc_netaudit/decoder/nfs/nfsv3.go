package nfs

import (
	nfsv32 "github.com/funbinary/go_example/example/internal/proc_netaudit/decoder/nfs/nfsv3"
	log "github.com/funbinary/go_example/pkg/blog"
	"github.com/funbinary/go_example/pkg/errors"
	"github.com/funbinary/go_example/pkg/rpc"
	"github.com/funbinary/go_example/pkg/xdr"
	"io"
)

func handlerGetAttr(h *Nfs, r io.Reader) (err error) {
	dir := h.CurMsgType
	if dir == rpc.RpcCall {
		call := &nfsv32.GetAttrCall{}
		h.Call = call
		err = xdr.Read(r, call)
		if err != nil {
			return errors.Wrap(err, "handlerGetAttr parse error")
		}
		log.Infof("GetAttrCall:%v", h.Call)
	} else if dir == rpc.RpcReply {
		reply := &nfsv32.GetAttrReply{}
		h.Reply = reply

		err = xdr.Read(r, reply)
		if err != nil {
			return errors.Wrap(err, "handlerGetAttrReply parse failed")
		}

		log.Infof("GetAttrReply:%v", h.Reply)
	} else {
		return errors.Errorf("解析Message Type:%d失败", dir)
	}
	return nil
}

func handleSetAttr(h *Nfs, r io.Reader) (err error) {
	dir := h.CurMsgType
	if dir == rpc.RpcCall {
		call := &nfsv32.SetAttrCall{}
		h.Call = call
		err = xdr.Read(r, call)
		if err != nil {
			return errors.Wrap(err, "handleSetAttr parse failed")
		}
		log.Infof("SetAttrCall:%v", h.Call)
	} else if dir == rpc.RpcReply {
		reply := &nfsv32.SetAttrReply{}
		h.Reply = reply

		err = xdr.Read(r, reply)
		if err != nil {
			return errors.Wrap(err, "handlerSetAttr parse failed")
		}
		log.Infof("SetAttrReply:%v", h.Reply)
	} else {
		return errors.Errorf("解析Message Type:%d失败", dir)
	}
	return nil
}

func handleLookUp(h *Nfs, r io.Reader) (err error) {
	dir := h.CurMsgType
	if dir == rpc.RpcCall {
		call := &nfsv32.LookUpCall{}
		h.Call = call
		err = xdr.Read(r, call)
		if err != nil {
			return errors.Wrap(err, "handleLookUpCall parse  failed")
		}
		log.Infof("LookupCall:%v", h.Call)
	} else if dir == rpc.RpcReply {
		reply := &nfsv32.LookUpReply{}
		h.Reply = reply
		err = xdr.Read(r, &reply.NFSStatus)
		if err != nil {
			return errors.Wrap(err, "handleLookUpReply parse status failed")
		}
		switch reply.NFSStatus {
		case nfsv32.NFSStatusOk:
			{
				err = xdr.Read(r, &reply.FH)
				if err != nil {
					return errors.Wrap(err, "handleLookUpReply parse object failed")
				}
				err = xdr.Read(r, &reply.Attr)
				if err != nil {
					return errors.Wrap(err, "handleLookUpReply parse attributes failed")
				}
				err = xdr.Read(r, &reply.DirAttr)
				if err != nil {
					return errors.Wrap(err, "handleLookUpReply parse dir attributes failed")
				}
			}

		default:
			err = xdr.Read(r, &reply.Attr)
			if err != nil {
				return errors.Wrap(err, "handleLookUpReply parse attr failed")
			}
		}
		log.Infof("LookUpReply:%v", h.Reply)
	} else {
		return errors.Errorf("解析Message Type:%d失败", dir)
	}
	return nil
}

func handleAccess(h *Nfs, r io.Reader) (err error) {
	dir := h.CurMsgType
	if dir == rpc.RpcCall {
		call := &nfsv32.AccessCall{}
		h.Call = call
		err = xdr.Read(r, call)
		if err != nil {
			return errors.Wrap(err, "handleAccessCall parse failed")
		}
		log.Infof("AccessCall:%v", h.Call)
	} else if dir == rpc.RpcReply {
		reply := &nfsv32.AccessReply{}
		h.Reply = reply
		err = xdr.Read(r, reply)
		if err != nil {
			return errors.Wrap(err, "handleAccessReply parse failed")
		}
		log.Infof("AccessReply:%v", h.Reply)
	} else {
		return errors.Errorf("解析Message Type:%d失败", dir)
	}
	return nil
}

func handleReadLink(h *Nfs, r io.Reader) (err error) {
	dir := h.CurMsgType
	if dir == rpc.RpcCall {
		call := &nfsv32.ReadLinkCall{}
		h.Call = call
		err = xdr.Read(r, call)
		if err != nil {
			return errors.Wrap(err, "handleReadLink parse failed")
		}
		log.Infof("ReadLinkCall:%v", h.Call)
	} else if dir == rpc.RpcReply {
		reply := &nfsv32.ReadLinkReply{}
		h.Reply = reply
		err = xdr.Read(r, reply)
		if err != nil {
			return errors.Wrap(err, "handleReadLinkReply parse failed")
		}
		log.Infof("ReadLinkReply:%v", h.Reply)
	} else {
		return errors.Errorf("解析Message Type:%d失败", dir)
	}
	return nil
}

func handleRead(h *Nfs, r io.Reader) (err error) {
	dir := h.CurMsgType
	if dir == rpc.RpcCall {
		call := &nfsv32.ReadCall{}
		h.Call = call
		err = xdr.Read(r, &call)
		if err != nil {
			return errors.Wrap(err, "handleReadCall parse  failed")
		}

		log.Infof("ReadCall:%v", h.Call)
	} else if dir == rpc.RpcReply {
		reply := &nfsv32.ReadReply{}
		h.Reply = reply
		err = xdr.Read(r, reply)
		if err != nil {
			return errors.Wrap(err, "handleReadReply parse failed")
		}
		log.Infof("ReadReply:%v", h.Reply)
	} else {
		return errors.Errorf("解析Message Type:%d失败", dir)
	}
	return nil
}

func handleWrite(h *Nfs, r io.Reader) (err error) {
	dir := h.CurMsgType
	if dir == rpc.RpcCall {
		call := &nfsv32.WriteCall{}
		h.Call = call
		err = xdr.Read(r, &call)
		if err != nil {
			return errors.Wrap(err, "handleWriteCall parse failed")
		}
		log.Infof("WriteCall:%v", h.Call)
	} else if dir == rpc.RpcReply {
		reply := &nfsv32.WriteReply{}
		h.Reply = reply

		err = xdr.Read(r, reply)
		if err != nil {
			return errors.Wrap(err, "handleWriteReply parse failed")
		}

		log.Infof("WriteReply:%v", h.Reply)

	} else {
		return errors.Errorf("解析Message Type:%d失败", dir)
	}
	return nil
}

func handleCreate(h *Nfs, r io.Reader) (err error) {
	dir := h.CurMsgType
	if dir == rpc.RpcCall {
		call := &nfsv32.CreateCall{}
		h.Call = call
		err = xdr.Read(r, call)
		if err != nil {
			return errors.Wrap(err, "handleCreateCall parse failed")
		}
		log.Infof("CreateCall:%v", h.Call)
	} else if dir == rpc.RpcReply {
		reply := &nfsv32.CreateReply{}
		h.Reply = reply
		err = xdr.Read(r, reply)
		if err != nil {
			return errors.Wrap(err, "handleCreateReply parse failed")
		}
		log.Infof("CreateReply:%v", h.Reply)

	} else {
		return errors.Errorf("解析Message Type:%d失败", dir)
	}
	return nil
}

func handleMkdir(h *Nfs, r io.Reader) (err error) {
	dir := h.CurMsgType
	if dir == rpc.RpcCall {
		call := &nfsv32.MkdirCall{}
		h.Call = call
		err = xdr.Read(r, call)
		if err != nil {
			return errors.Wrap(err, "handleMkdirCall parse failed")
		}
		log.Infof("MkdirCall:%v", h.Call)
	} else if dir == rpc.RpcReply {
		reply := &nfsv32.MkdirReply{}
		h.Reply = reply
		err = xdr.Read(r, reply)
		if err != nil {
			return errors.Wrap(err, "handleMkdirReply parse failed")
		}
		log.Infof("MkdirReply:%v", h.Reply)

	} else {
		return errors.Errorf("解析Message Type:%d失败", dir)
	}
	return nil
}

func handleSymLink(h *Nfs, r io.Reader) (err error) {
	dir := h.CurMsgType
	if dir == rpc.RpcCall {
		call := &nfsv32.SymlinkCall{}
		h.Call = call
		err = xdr.Read(r, call)
		if err != nil {
			return errors.Wrap(err, "handleSymlinkCall parse failed")
		}
		log.Infof("SymlinkCall:%v", h.Call)
	} else if dir == rpc.RpcReply {
		reply := &nfsv32.SymlinkReply{}
		h.Reply = reply

		err = xdr.Read(r, reply)
		if err != nil {
			return errors.Wrap(err, "handleSymLinkReply parse failed")
		}
		log.Infof("SymlinkReply:%v", h.Reply)

	} else {
		return errors.Errorf("解析Message Type:%d失败", dir)
	}
	return nil
}

func handleMknod(h *Nfs, r io.Reader) (err error) {
	dir := h.CurMsgType
	if dir == rpc.RpcCall {
		call := &nfsv32.MknodCall{}
		h.Call = call
		err = xdr.Read(r, call)
		if err != nil {
			return errors.Wrap(err, "handleMknodCall parse failed")
		}
		log.Infof("MknodCall:%v", h.Call)
	} else if dir == rpc.RpcReply {
		reply := &nfsv32.MknodReply{}
		h.Reply = reply

		err = xdr.Read(r, reply)
		if err != nil {
			return errors.Wrap(err, "handleMknodReply parse failed")
		}
		log.Infof("MknodReply:%v", h.Reply)

	} else {
		return errors.Errorf("解析Message Type:%d失败", dir)
	}
	return nil
}

func handleRemove(h *Nfs, r io.Reader) (err error) {
	dir := h.CurMsgType
	if dir == rpc.RpcCall {
		call := &nfsv32.RemoveCall{}
		h.Call = call
		err = xdr.Read(r, call)
		if err != nil {
			return errors.Wrap(err, "handleRemoveCall parse failed")
		}
		log.Infof("RemoveCall:%v", h.Call)
	} else if dir == rpc.RpcReply {
		reply := &nfsv32.RemoveReply{}
		h.Reply = reply
		err = xdr.Read(r, reply)
		if err != nil {
			return errors.Wrap(err, "handleRemoveReply parse failed")
		}
		log.Infof("RemoveReply:%v", h.Reply)

	} else {
		return errors.Errorf("解析Message Type:%d失败", dir)
	}
	return nil
}

func handleRmdir(h *Nfs, r io.Reader) (err error) {
	dir := h.CurMsgType
	if dir == rpc.RpcCall {
		call := &nfsv32.RmdirCall{}
		h.Call = call
		err = xdr.Read(r, call)
		if err != nil {
			return errors.Wrap(err, "handleRmdirCall parse failed")
		}
		log.Infof("RmdirCall:%v", h.Call)
	} else if dir == rpc.RpcReply {
		reply := &nfsv32.RemoveReply{}
		h.Reply = reply

		err = xdr.Read(r, reply)
		if err != nil {
			return errors.Wrap(err, "handleRmdirReply parse failed")
		}
		log.Infof("RmdirReply:%v", h.Reply)

	} else {
		return errors.Errorf("解析Message Type:%d失败", dir)
	}
	return nil
}

func handleRename(h *Nfs, r io.Reader) (err error) {
	dir := h.CurMsgType
	if dir == rpc.RpcCall {
		call := &nfsv32.RenameCall{}
		h.Call = call
		err = xdr.Read(r, call)
		if err != nil {
			return errors.Wrap(err, "handleRenameCall parse failed")
		}
		log.Infof("RenameCall:%v", h.Call)
	} else if dir == rpc.RpcReply {
		reply := &nfsv32.RenameReply{}
		h.Reply = reply

		err = xdr.Read(r, reply)
		if err != nil {
			return errors.Wrap(err, "handleRenameReply parse failed")
		}
		log.Infof("RenameReply:%v", h.Reply)

	} else {
		return errors.Errorf("解析Message Type:%d失败", dir)
	}
	return nil
}

func handleLink(h *Nfs, r io.Reader) (err error) {
	dir := h.CurMsgType
	if dir == rpc.RpcCall {
		call := &nfsv32.LinkCall{}
		h.Call = call
		err = xdr.Read(r, call)
		if err != nil {
			return errors.Wrap(err, "handleLinkCall  failed")
		}
		log.Infof("LinkCall:%v", h.Call)
	} else if dir == rpc.RpcReply {
		reply := &nfsv32.LinkReply{}
		h.Reply = reply
		err = xdr.Read(r, reply)
		if err != nil {
			return errors.Wrap(err, "handleLinkReply  failed")
		}
		log.Infof("LinkReply:%v", h.Reply)

	} else {
		return errors.Errorf("解析Message Type:%d失败", dir)
	}
	return nil
}

func handleReadDir(h *Nfs, r io.Reader) (err error) {
	dir := h.CurMsgType
	if dir == rpc.RpcCall {
		call := &nfsv32.ReadDirCall{}
		h.Call = call
		err = xdr.Read(r, call)
		if err != nil {
			return errors.Wrap(err, "handleReadDir  failed")
		}
		log.Infof("ReadDirCall:%v", h.Call)
	} else if dir == rpc.RpcReply {
		reply := &nfsv32.ReadDirReply{}
		h.Reply = reply
		err = xdr.Read(r, &reply.ReadDirArg)
		if err != nil {
			return errors.Wrap(err, "handleReadDirReply parse arg failed")
		}
		for {
			var item nfsv32.Entry3
			if err = xdr.Read(r, &item); err != nil {
				return errors.Wrap(err, "handleReadDirReply parse entry failed")
			}

			if !item.IsSet {
				break
			}
			reply.Entries = append(reply.Entries, &item.Entry)
		}
		log.Infof("ReadDirReply:%v", h.Reply)

	} else {
		return errors.Errorf("解析Message Type:%d失败", dir)
	}
	return nil
}

func handleReadDirPlus(h *Nfs, r io.Reader) (err error) {
	dir := h.CurMsgType
	if dir == rpc.RpcCall {
		call := &nfsv32.ReadDirPlusCall{}
		h.Call = call
		err = xdr.Read(r, call)
		if err != nil {
			return errors.Wrap(err, "handleReadDirPlusCall parse  failed")
		}
		log.Infof("ReadDirPlusCall:%v", h.Call)
	} else if dir == rpc.RpcReply {
		reply := &nfsv32.ReadDirPlusReply{}
		h.Reply = reply
		err = xdr.Read(r, &reply.ReaDirPlusArg)
		if err != nil {
			return errors.Wrap(err, "handleReadDirPlusReply parse arg failed")
		}
		for {
			var item nfsv32.EntryPlus3
			if err = xdr.Read(r, &item); err != nil {
				return errors.Wrap(err, "handleReadDirPlusReply parse entry plus3 failed")
			}
			log.Infof("%v", item)
			if !item.IsSet {
				break
			}

			reply.Entries = append(reply.Entries, &item.Entry)
		}
		log.Infof("ReadDirPlusReply:%v", h.Reply)

	} else {
		return errors.Errorf("解析Message Type:%d失败", dir)
	}
	return nil
}

func handleFsStat(h *Nfs, r io.Reader) (err error) {
	dir := h.CurMsgType
	if dir == rpc.RpcCall {
		call := &nfsv32.FsStatCall{}
		h.Call = call
		err = xdr.Read(r, call)
		if err != nil {
			return errors.Wrap(err, "handleFsStatCall parse  failed")
		}
		log.Infof("FsStatCall:%v", h.Call)
	} else if dir == rpc.RpcReply {
		reply := &nfsv32.FsStatReply{}
		h.Reply = reply
		err = xdr.Read(r, reply)
		if err != nil {
			return errors.Wrap(err, "handleFsStatReply  failed")
		}
		log.Infof("FsStatReply:%v", h.Reply)

	} else {
		return errors.Errorf("解析Message Type:%d失败", dir)
	}
	return nil
}

func handleFsInfo(h *Nfs, r io.Reader) (err error) {
	dir := h.CurMsgType
	if dir == rpc.RpcCall {
		call := &nfsv32.FsInfoCall{}
		h.Call = call
		err = xdr.Read(r, call)
		if err != nil {
			return errors.Wrap(err, "handleFsInfoCall parse failed")
		}

		log.Infof("FsInfoCall:%v", h.Call)
	} else if dir == rpc.RpcReply {
		reply := &nfsv32.FsInfoReply{}
		h.Reply = reply
		err = xdr.Read(r, reply)
		if err != nil {
			return errors.Wrap(err, "handleFsInfoReply parse  failed")
		}
		log.Infof("FsInfoReply:%v", h.Reply)

	} else {
		return errors.Errorf("解析Message Type:%d失败", dir)
	}
	return nil
}

func handlePathConf(h *Nfs, r io.Reader) (err error) {
	dir := h.CurMsgType
	if dir == rpc.RpcCall {
		call := &nfsv32.PathConfCall{}
		h.Call = call
		err = xdr.Read(r, call)
		if err != nil {
			return errors.Wrap(err, "handlePathConf parse failed")
		}

		log.Infof("PathConfCall:%v", h.Call)
	} else if dir == rpc.RpcReply {
		reply := &nfsv32.PathConfReply{}
		h.Reply = reply
		err = xdr.Read(r, reply)
		if err != nil {
			return errors.Wrap(err, "handlePathConfReply  failed")
		}
		log.Infof("PathConfReply:%v", h.Reply)

	} else {
		return errors.Errorf("解析Message Type:%d失败", dir)
	}
	return nil
}

func handleCommit(h *Nfs, r io.Reader) (err error) {
	dir := h.CurMsgType
	if dir == rpc.RpcCall {
		call := &nfsv32.CommitCall{}
		h.Call = call
		err = xdr.Read(r, call)
		if err != nil {
			return errors.Wrap(err, "handleCommitCall parse failed")
		}
		log.Infof("CommitCall:%v", h.Call)
	} else if dir == rpc.RpcReply {
		reply := &nfsv32.CommitReply{}
		h.Reply = reply
		err = xdr.Read(r, reply)
		if err != nil {
			return errors.Wrap(err, "handleCommitReply  failed")
		}
		log.Infof("CommitReply:%v", h.Reply)

	} else {
		return errors.Errorf("解析Message Type:%d失败", dir)
	}
	return nil
}
