package rpc

type Auth struct {
	Flavor uint32
	Body   []byte
}

var AuthNull Auth

type Base struct {
	Xid        uint32
	CurMsgType MessageType
}

type ReqHeader struct {
	Rpcvers   uint32      // rpc版本号，固定为2
	Program   ProgramType // 类型
	ProgVers  uint32      //
	Procedure uint32
	Cred      Auth
	Verf      Auth
}
type RepHeader struct {
	State  uint32
	Verf   Auth
	Status AcceptStatus
}

type MessageType uint32

const (
	RpcCall  MessageType = 0
	RpcReply MessageType = 1
)

func (self MessageType) String() string {
	switch self {
	case RpcCall:
		return "Call"
	case RpcReply:
		return "Reply"
	default:
		return "Unknown"
	}
}

type ProgramType uint32

const (
	PortmapID      ProgramType = 100000
	NfsServiceID   ProgramType = 100003
	MountServiceID ProgramType = 100005
)

func (p ProgramType) String() string {
	switch p {
	case PortmapID:
		return "portmap"
	case NfsServiceID:
		return "nfs"
	case MountServiceID:
		return "mount"
	default:
		return "Unknown"
	}
}

// PORTMAP
// RFC 1057 Section A.1

const (
	PmapPort = 111
	PmapVers = 2
)

type AcceptStatus uint32

const (
	Success AcceptStatus = iota
	ProgUnavail
	ProgMismatch
	ProcUnavail
	GarbageArgs
	SystemErr
)
