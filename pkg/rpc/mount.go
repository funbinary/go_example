package rpc

// MountProcedure is the valid RPC calls for the mount service.
type MountProcedure uint32

// MountProcedure Codes
const (
	MountProcNull MountProcedure = iota
	MountProcMount
	MountProcDump
	MountProcUmnt
	MountProcUmntAll
	MountProcExport
)

func (m MountProcedure) String() string {
	switch m {
	case MountProcNull:
		return "Null"
	case MountProcMount:
		return "Mount"
	case MountProcDump:
		return "Dump"
	case MountProcUmnt:
		return "Umnt"
	case MountProcUmntAll:
		return "UmntAll"
	case MountProcExport:
		return "Export"
	default:
		return "Unknown"
	}
}

type MountStat3 uint32

const (
	MNT3_OK             MountStat3 = iota   // success
	MNT3ERR_PERM                            // Not owner
	MNT3ERR_NOENT                           // No such file or directory
	MNT3ERR_IO          MountStat3 = 5      // I/O error
	MNT3ERR_ACCES       MountStat3 = 12     // Permission denied
	MNT3ERR_NOTDIR      MountStat3 = 13     // Not a directory
	MNT3ERR_INVAL       MountStat3 = 22     // Invalid argument
	MNT3ERR_NAMETOOLONG MountStat3 = 63     // Filename too long
	MNT3ERR_NOTSUPP     MountStat3 = 10004  // Operation not supported
	MNT3ERR_SERVERFAULT MountStat3 = 100006 // A failure on the server

)
