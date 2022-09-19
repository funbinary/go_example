package nfsv3

import (
	"encoding/hex"
	"fmt"
	"github.com/funbinary/go_example/pkg/charset"
	"github.com/funbinary/go_example/pkg/xdr"
	"golang.org/x/text/encoding/simplifiedchinese"
	"io"
	"strings"
	"time"
)

type Fattr struct {
	Type FileType
	//todo:展示详细信息S_ISUID、S_ISGID、S_ISVTX、S_IRUSR、S_IWUSR、S_IXUSR、S_IRGRP、S_IWGRP、S_IXGRP、S_IROTH、S_IWOTH、S_IXOTH
	FileMode            uint32
	Nlink               uint32
	UID                 uint32
	GID                 uint32
	Filesize            uint64
	Used                uint64
	SpecData            [2]uint32
	FSID                uint64
	Fileid              uint64
	Atime, Mtime, Ctime NFS3Time
}

func (self *Fattr) String() string {
	var s strings.Builder
	s.WriteString(fmt.Sprintln("    type:", self.Type))
	s.WriteString(fmt.Sprintf("    mode:%o\n", self.FileMode))
	s.WriteString(fmt.Sprintln("    nlink:", self.Nlink))
	s.WriteString(fmt.Sprintln("    uid:", self.UID))
	s.WriteString(fmt.Sprintln("    gid:", self.GID))
	s.WriteString(fmt.Sprintln("    size:", self.Filesize))
	s.WriteString(fmt.Sprintln("    used:", self.Used))
	s.WriteString(fmt.Sprintln("    rdev:", self.SpecData))
	s.WriteString(fmt.Sprintln("    fsid:", self.FSID))
	s.WriteString(fmt.Sprintln("    fileid:", self.Fileid))
	s.WriteString(fmt.Sprintln("    atime:", self.Atime))
	s.WriteString(fmt.Sprintln("    mtime:", self.Mtime))
	s.WriteString(fmt.Sprintln("    ctime:", self.Ctime))
	return s.String()

}

func ReadAttributes(r io.Reader) (*Fattr, error) {
	hasAttr, err := xdr.ReadUint32(r)
	if err != nil {
		return nil, err
	}
	if hasAttr == 1 {
		var attr Fattr
		err = xdr.Read(r, &attr)
		if err != nil {
			return nil, err
		}

		return &attr, err
	}
	return nil, nil

}

type NFSStatus uint32

// NFSStatus codes
const (
	NFSStatusOk          NFSStatus = 0
	NFSStatusPerm        NFSStatus = 1
	NFSStatusNoEnt       NFSStatus = 2
	NFSStatusIO          NFSStatus = 5
	NFSStatusNXIO        NFSStatus = 6
	NFSStatusAccess      NFSStatus = 13
	NFSStatusExist       NFSStatus = 17
	NFSStatusXDev        NFSStatus = 18
	NFSStatusNoDev       NFSStatus = 19
	NFSStatusNotDir      NFSStatus = 20
	NFSStatusIsDir       NFSStatus = 21
	NFSStatusInval       NFSStatus = 22
	NFSStatusFBig        NFSStatus = 27
	NFSStatusNoSPC       NFSStatus = 28
	NFSStatusROFS        NFSStatus = 30
	NFSStatusMlink       NFSStatus = 31
	NFSStatusNameTooLong NFSStatus = 63
	NFSStatusNotEmpty    NFSStatus = 66
	NFSStatusDQuot       NFSStatus = 69
	NFSStatusStale       NFSStatus = 70
	NFSStatusRemote      NFSStatus = 71
	NFSStatusBadHandle   NFSStatus = 10001
	NFSStatusNotSync     NFSStatus = 10002
	NFSStatusBadCookie   NFSStatus = 10003
	NFSStatusNotSupp     NFSStatus = 10004
	NFSStatusTooSmall    NFSStatus = 10005
	NFSStatusServerFault NFSStatus = 10006
	NFSStatusBadType     NFSStatus = 10007
	NFSStatusJukebox     NFSStatus = 10008
)

func (s NFSStatus) String() string {
	switch s {
	case NFSStatusOk:
		return "Call Completed Successfull"
	case NFSStatusPerm:
		return "Not Owner"
	case NFSStatusNoEnt:
		return "No such file or directory"
	case NFSStatusIO:
		return "I/O error"
	case NFSStatusNXIO:
		return "I/O error: No such device"
	case NFSStatusAccess:
		return "Permission denied"
	case NFSStatusExist:
		return "File exists"
	case NFSStatusXDev:
		return "Attempt to do a cross device hard link"
	case NFSStatusNoDev:
		return "No such device"
	case NFSStatusNotDir:
		return "Not a directory"
	case NFSStatusIsDir:
		return "Is a directory"
	case NFSStatusInval:
		return "Invalid argument"
	case NFSStatusFBig:
		return "File too large"
	case NFSStatusNoSPC:
		return "No space left on device"
	case NFSStatusROFS:
		return "Read only file system"
	case NFSStatusMlink:
		return "Too many hard links"
	case NFSStatusNameTooLong:
		return "Name too long"
	case NFSStatusNotEmpty:
		return "Not empty"
	case NFSStatusDQuot:
		return "Resource quota exceeded"
	case NFSStatusStale:
		return "Invalid file handle"
	case NFSStatusRemote:
		return "Too many levels of remote in path"
	case NFSStatusBadHandle:
		return "Illegal NFS file handle"
	case NFSStatusNotSync:
		return "Synchronization mismatch"
	case NFSStatusBadCookie:
		return "Cookie is Stale"
	case NFSStatusNotSupp:
		return "Operation not supported"
	case NFSStatusTooSmall:
		return "Buffer or request too small"
	case NFSStatusServerFault:
		return "Unmapped error (EIO)"
	case NFSStatusBadType:
		return "Type not supported"
	case NFSStatusJukebox:
		return "Initiated, but too slow. Try again with new txn"
	default:
		return "unknown"
	}
}

type FileType uint32

// Enumeration of NFS FileTypes
const (
	FileTypeRegular FileType = iota + 1
	FileTypeDirectory
	FileTypeBlock
	FileTypeCharacter
	FileTypeLink
	FileTypeSocket
	FileTypeFIFO
)

func (f FileType) String() string {
	switch f {
	case FileTypeRegular:
		return "Regular"
	case FileTypeDirectory:
		return "Directory"
	case FileTypeBlock:
		return "Block Device"
	case FileTypeCharacter:
		return "Character Device"
	case FileTypeLink:
		return "Symbolic Link"
	case FileTypeSocket:
		return "Socket"
	case FileTypeFIFO:
		return "FIFO"
	default:
		return "Unknown"
	}
}

// FileTime is the NFS wire time format
// This is equivalent to go-nfs-client/nfs.NFS3Time
type NFS3Time struct {
	Seconds  uint32
	Nseconds uint32
}

// ToNFSTime generates the nfs 64bit time format from a golang time.
func ToNFSTime(t time.Time) NFS3Time {
	return NFS3Time{
		Seconds:  uint32(t.Unix()),
		Nseconds: uint32(t.UnixNano()) % uint32(time.Second),
	}
}

func (t NFS3Time) String() string {
	ts := time.Unix(int64(t.Seconds), int64(t.Nseconds))
	return ts.Format("2006-01-02 15:04:05.000000000")
}

// Native generates a golang time from an nfs time spec
func (t NFS3Time) Native() *time.Time {
	ts := time.Unix(int64(t.Seconds), int64(t.Nseconds))
	return &ts
}

// EqualTimespec returns if this time is equal to a local time spec
func (t NFS3Time) EqualTimespec(sec int64, nsec int64) bool {
	// TODO: bounds check on sec/nsec overflow
	return t.Nseconds == uint32(nsec) && t.Seconds == uint32(sec)
}

type Sattr3 struct {
	Mode  SetMode
	UID   SetID
	GID   SetID
	Size  SetSize
	Atime SetTime
	Mtime SetTime
}

func (self Sattr3) HasSet() bool {
	return self.Mode.SetIt || self.UID.SetIt || self.GID.SetIt || self.Size.SetIt || self.Atime.SetIt != DontChange || self.Mtime.SetIt != DontChange
}

func (self Sattr3) String() string {
	if !self.HasSet() {
		return ""
	}
	var s strings.Builder
	if self.Mode.SetIt {
		s.WriteString(fmt.Sprintf("Mode:%o\n", self.Mode.Mode))
	}
	if self.UID.SetIt {
		s.WriteString(fmt.Sprintln("UID:", self.UID.ID))
	}
	if self.GID.SetIt {
		s.WriteString(fmt.Sprintln("GID:", self.GID.ID))
	}
	if self.Size.SetIt {
		s.WriteString(fmt.Sprintln("Size:", self.Size.Size))
	}
	if self.Atime.SetIt != DontChange {
		s.WriteString(fmt.Sprintln("Atime:", self.Atime.Time))
	}
	if self.Mtime.SetIt != DontChange {
		s.WriteString(fmt.Sprintln("Mtime:", self.Mtime.Time))
	}

	return s.String()
}

type SetMode struct {
	SetIt bool   `xdr:"union"`
	Mode  uint32 `xdr:"unioncase=1"`
}

type SetID struct {
	SetIt bool   `xdr:"union"`
	ID    uint32 `xdr:"unioncase=1"`
}

type SetSize struct {
	SetIt bool   `xdr:"union"`
	Size  uint64 `xdr:"unioncase=1"`
}

type SetTime struct {
	SetIt TimeHow  `xdr:"union"`
	Time  NFS3Time `xdr:"unioncase=2"` //SetToClientTime
}

// TimeHow
// DONT_CHANGE        = 0
// SET_TO_SERVER_TIME = 1
// SET_TO_CLIENT_TIME = 2
type TimeHow int

const (
	DontChange TimeHow = iota
	SetToServerTime
	SetToClientTime
)

type Sattrguard3 struct {
	Check bool     `xdr:"union"`
	Time  NFS3Time `xdr:"unioncase=1"` //SetToClientTime
}

func (self Sattrguard3) HasCheck() bool {
	return self.Check
}

func (self Sattrguard3) String() string {
	return self.Time.String()
}

type FileCacheAttribute struct {
	Filesize     uint64
	Mtime, Ctime NFS3Time
}

func (self *FileCacheAttribute) String() string {
	var s strings.Builder
	s.WriteString(fmt.Sprintln("    filesize:", self.Filesize))
	s.WriteString(fmt.Sprintln("    mtime:", self.Mtime))
	s.WriteString(fmt.Sprintln("    ctime:", self.Ctime))
	return s.String()
}

type Wcc struct {
	Before PreOpAttr
	After  PostOpAttr
}

func (self Wcc) HasSet() bool {
	return self.Before.IsSet || self.After.IsSet
}

func (self *Wcc) String() string {
	var s strings.Builder
	if self.Before.IsSet {
		s.WriteString(fmt.Sprintln("before:"))
		s.WriteString(self.Before.String())
	}
	if self.After.IsSet {
		s.WriteString(fmt.Sprintln("after:"))
		s.WriteString(self.After.String())
	}
	return s.String()
}

type DirOpArg struct {
	FH       []byte
	Filename []byte
}

func (self DirOpArg) FileName() string {
	encoding := charset.GetStrCoding(self.Filename)
	switch encoding {
	case charset.GBK:
		utf8Data, _ := simplifiedchinese.GBK.NewDecoder().Bytes(self.Filename) //将gbk再转换为utf-8
		return string(utf8Data)
	}
	return string(self.Filename)
}

func (self DirOpArg) String() string {
	var s strings.Builder
	s.WriteString(fmt.Sprintln("handle:", hex.EncodeToString(self.FH)))
	encoding := charset.GetStrCoding(self.Filename)
	switch encoding {
	case charset.GBK:
		utf8Data, _ := simplifiedchinese.GBK.NewDecoder().Bytes(self.Filename) //将gbk再转换为utf-8
		s.WriteString(fmt.Sprintln("name:", string(utf8Data)))
	default:
		s.WriteString(fmt.Sprintln("name:", string(self.Filename)))
	}
	return s.String()
}

type Access uint32

type SymLinkData3 struct {
	Sattr3
	Nfs3Path []byte
}

func (self *SymLinkData3) String() string {
	var s strings.Builder
	s.WriteString(fmt.Sprintln("\nset file attributes:\n", self.Sattr3))
	s.WriteString(fmt.Sprintln("name:", string(self.Nfs3Path)))
	return s.String()
}

type PostOpFH3 struct {
	IsSet      bool   `xdr:"union"`
	FileHandle []byte `xdr:"unioncase=1"`
}

func (self PostOpFH3) String() string {
	if !self.IsSet {
		return ""
	}
	return hex.EncodeToString(self.FileHandle)
}

type PreOpAttr struct {
	IsSet bool     `xdr:"union"`
	Size  uint64   `xdr:"unioncase=1"`
	MTime NFS3Time `xdr:"unioncase=1"`
	CTime NFS3Time `xdr:"unioncase=1"`
}

func (self PreOpAttr) String() string {
	if !self.IsSet {
		return ""
	}
	var s strings.Builder
	s.WriteString(fmt.Sprintln("	size:", self.Size))
	s.WriteString(fmt.Sprintln("	Mtime:", self.MTime))
	s.WriteString(fmt.Sprintln("	CTime:", self.CTime))
	return s.String()
}

type PostOpAttr struct {
	IsSet bool  `xdr:"union"`
	Attr  Fattr `xdr:"unioncase=1"`
}

func (self PostOpAttr) HasSet() bool {
	return self.IsSet
}

func (self PostOpAttr) String() string {
	if !self.IsSet {
		return ""
	}

	return self.Attr.String()
}
