package rpc

type NfsVersion uint32

const (
	NfsV2 NfsVersion = 2
	NfsV3 NfsVersion = 3
	NfsV4 NfsVersion = 4
)

func (self NfsVersion) String() string {
	switch self {
	case NfsV2:
		return "nfs2"
	case NfsV3:
		return "nfs3"
	case NfsV4:
		return "nfs4"
	default:
		return "Unknown"
	}
}
