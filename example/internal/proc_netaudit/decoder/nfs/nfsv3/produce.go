package nfsv3

// Procedure is the valid RPC calls for the nfs service.
type Procedure uint32

// NfsProcedure Codes
const (
	NFS3ProcedureNull Procedure = iota
	NFS3ProcedureGetAttr
	NFS3ProcedureSetAttr
	NFS3ProcedureLookup
	NFS3ProcedureAccess
	NFS3ProcedureReadlink
	NFS3ProcedureRead
	NFS3ProcedureWrite
	NFS3ProcedureCreate
	NFS3ProcedureMkDir
	NFS3ProcedureSymlink
	NFS3ProcedureMkNod
	NFS3ProcedureRemove
	NFS3ProcedureRmDir
	NFS3ProcedureRename
	NFS3ProcedureLink
	NFS3ProcedureReadDir
	NFS3ProcedureReadDirPlus
	NFS3ProcedureFSStat
	NFS3ProcedureFSInfo
	NFS3ProcedurePathConf
	NFS3ProcedureCommit
)

func (n Procedure) String() string {
	switch n {
	case NFS3ProcedureNull:
		return "Null"
	case NFS3ProcedureGetAttr:
		return "GetAttr"
	case NFS3ProcedureSetAttr:
		return "SetAttr"
	case NFS3ProcedureLookup:
		return "Lookup"
	case NFS3ProcedureAccess:
		return "Access"
	case NFS3ProcedureReadlink:
		return "ReadLink"
	case NFS3ProcedureRead:
		return "Read"
	case NFS3ProcedureWrite:
		return "Write"
	case NFS3ProcedureCreate:
		return "Create"
	case NFS3ProcedureMkDir:
		return "Mkdir"
	case NFS3ProcedureSymlink:
		return "Symlink"
	case NFS3ProcedureMkNod:
		return "Mknod"
	case NFS3ProcedureRemove:
		return "Remove"
	case NFS3ProcedureRmDir:
		return "Rmdir"
	case NFS3ProcedureRename:
		return "Rename"
	case NFS3ProcedureLink:
		return "Link"
	case NFS3ProcedureReadDir:
		return "ReadDir"
	case NFS3ProcedureReadDirPlus:
		return "ReadDirPlus"
	case NFS3ProcedureFSStat:
		return "FSStat"
	case NFS3ProcedureFSInfo:
		return "FSInfo"
	case NFS3ProcedurePathConf:
		return "PathConf"
	case NFS3ProcedureCommit:
		return "Commit"
	default:
		return "Unknown"
	}
}
