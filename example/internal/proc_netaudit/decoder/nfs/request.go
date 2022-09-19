package nfs

import (
	"fmt"
	"github.com/funbinary/go_example/example/internal/proc_netaudit/decoder/nfs/nfsv3"
	"github.com/funbinary/go_example/pkg/rpc"
)

type Nfs struct {
	rpc.Base
	rpc.ReqHeader
	rpc.RepHeader
	Call  interface{}
	Reply interface{}
}

func NewNfsHeader() *Nfs {
	return &Nfs{}
}

func (*Nfs) Protocol() string {
	return "nfs"
}

func (self *Nfs) String() string {
	switch self.Program {
	case rpc.PortmapID:
		// RPC Call #xid nfs.3.read
		return fmt.Sprintf("RPC %v #%d (%s.%s.%s)", self.CurMsgType, self.Xid, self.Program, self.ReqHeader.ProgVers, rpc.PortmapProcedure(self.Procedure))
	case rpc.NfsServiceID:
		switch rpc.NfsVersion(self.ProgVers) {
		//case NfsV2:
		//	return fmt.Sprintf("RPC %v #%d (%s.%s.%s)", MessageType(self.CurMsgType), self.Xid, ProgramType(self.Program), self.ReqHeader.ProgVers, nfsv3.Procedure(self.Procedure))
		case rpc.NfsV3:
			return fmt.Sprintf("RPC %v #%d (%s.%s.%s)", self.CurMsgType, self.Xid, self.Program, self.ReqHeader.ProgVers, nfsv3.Procedure(self.Procedure))
		//case NfsV4:
		//	return fmt.Sprintf("RPC %v #%d (%s.%s.%s)", MessageType(self.CurMsgType), self.Xid, ProgramType(self.Program), self.ReqHeader.ProgVers, nfsv3.Procedure(self.Procedure))
		default:
			return "Unknown"
		}
	case rpc.MountServiceID:
		return fmt.Sprintf("RPC %v #%d (%s.%s.%s)", self.CurMsgType, self.Xid, self.Program, self.ReqHeader.ProgVers, rpc.MountProcedure(self.Procedure))
	default:
		return fmt.Sprintf("RPC %v #%d (%s.%s.unknown)", self.CurMsgType, self.Xid, self.Program, self.ReqHeader.ProgVers)
	}
}
