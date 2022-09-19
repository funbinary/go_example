package nfsv2

// Procedure is the valid RPC calls for the nfs service.
type Procedure uint32

// NfsProcedure Codes
const (
	NFS2ProcedureNull Procedure = iota
	NFS2ProcedureGetAttr
	NFS2ProcedureSetAttr
	NFS2ProcedureRoot // 无用
	NFS2ProcedureLookUp
	NFS2ProcedureReadlink
	NFS2ProcedureRead
	NFS2ProcedureWriteCache // 无用
	NFS2ProcedureWrite
	NFS2ProcedureCreate
	NFS2ProcedureRemove
	NFS2ProcedureRename
	NFS2ProcedureLink
	NFS2ProcedureSymlink
	NFS2ProcedureMkDir
	NFS2ProcedureRmDir
	NFS2ProcedureReadDir
	NFS2ProcedureStatFs
)

func (n Procedure) String() string {
	switch n {
	case NFS2ProcedureNull:
		return "Null"
	case NFS2ProcedureGetAttr:
		return "GetAttr"
	case NFS2ProcedureSetAttr:
		return "SetAttr"
	case NFS2ProcedureRoot:
		return "Root"
	case NFS2ProcedureLookUp:
		return "LookUp"
	case NFS2ProcedureReadlink:
		return "ReadLink"
	case NFS2ProcedureRead:
		return "Read"
	case NFS2ProcedureWriteCache:
		return "WriteCache"
	case NFS2ProcedureWrite:
		return "Write"
	case NFS2ProcedureCreate:
		return "Create"
	case NFS2ProcedureRemove:
		return "Remove"
	case NFS2ProcedureRename:
		return "Rename"
	case NFS2ProcedureLink:
		return "Link"
	case NFS2ProcedureSymlink:
		return "Symlink"
	case NFS2ProcedureMkDir:
		return "Mkdir"
	case NFS2ProcedureRmDir:
		return "Rmdir"
	case NFS2ProcedureReadDir:
		return "ReadDir"
	case NFS2ProcedureStatFs:
		return "StatFs"
	default:
		return "Unknown"
	}
}
