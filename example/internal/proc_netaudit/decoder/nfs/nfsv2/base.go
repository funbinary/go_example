package nfsv2

import (
	"encoding/hex"
	"fmt"
	"github.com/funbinary/go_example/pkg/charset"
	"golang.org/x/text/encoding/simplifiedchinese"
	"strings"
)

type Nfs2Status uint32

const (
	NFS2StatusOK Nfs2Status = iota
	NFS2StatusPerm
	NFS2StatusNoEnt
	NFS2StatusIO
	NFS2StatusNXIO
	NFS2StatusAccess
	NFS2StatusExist
	NFS2StatusNoDev
	NFS2StatusNotDir
	NFS2StatusIsDir
	NFS2StatusFBig
	NFS2StatusNoSPC
	NFS2StatusROFS
	NFS2StatusTooLong
	NFS2StatusNotEmpty
	NFS2StatusDQuot
	NFS2StatusStale
	NFS2StatusWFlush
)

func (s Nfs2Status) String() string {
	switch s {
	case NFS2StatusOK:
		return "Call Completed Successfull"
	case NFS2StatusPerm:
		return "Not Owner"
	case NFS2StatusNoEnt:
		return "No such file or directory"
	case NFS2StatusIO:
		return "I/O error"
	case NFS2StatusNXIO:
		return "I/O error: No such device"
	case NFS2StatusAccess:
		return "Permission denied"
	case NFS2StatusExist:
		return "File exists"
	case NFS2StatusNoDev:
		return "No such device"
	case NFS2StatusNotDir:
		return "Not a directory"
	case NFS2StatusIsDir:
		return "Is a directory"
	case NFS2StatusFBig:
		return "File too large"
	case NFS2StatusNoSPC:
		return "No space left on device"
	case NFS2StatusROFS:
		return "Read only file system"
	case NFS2StatusTooLong:
		return "Name too long"
	case NFS2StatusNotEmpty:
		return "Not empty"
	case NFS2StatusDQuot:
		return "Resource quota exceeded"
	case NFS2StatusStale:
		return "Invalid file handle"
	case NFS2StatusWFlush:
		return "flush cache error"
	default:
		return "unknown"
	}
}

type Nfs2Ftype uint32

const (
	Nfs2FtypeNFNON Nfs2Ftype = iota
	Nfs2FtypeNFREG
	Nfs2FtypeNFDIR
	Nfs2FtypeNFBLK
	Nfs2FtypeNFCHR
	Nfs2FtypeNFLNK
)

func (f Nfs2Ftype) String() string {
	switch f {
	case Nfs2FtypeNFNON:
		return "Not file"
	case Nfs2FtypeNFREG:
		return "regular file"
	case Nfs2FtypeNFDIR:
		return "Directory"
	case Nfs2FtypeNFBLK:
		return "Block Device"
	case Nfs2FtypeNFCHR:
		return "Character Device"
	case Nfs2FtypeNFLNK:
		return "Link File"
	default:
		return "Unknown"
	}
}

type TimeVal struct {
	Second   uint32 // 秒
	USeconds uint32 // 微妙
}

type Fattr struct {
	Type                Nfs2Ftype //文件类型
	Mode                uint32    //mode
	NLink               uint32    //文件的硬链接数
	UID                 uint32    //文件的所有者的用户标识号码
	GID                 uint32    //文件的组的组标识号码
	Size                uint32    //以字节数计算的大小
	BlockSize           uint32
	Rdev                uint32
	Blocks              uint32
	FsId                uint32
	FileId              uint32
	Atime, Mtime, Ctime TimeVal
}

func (self Fattr) String() string {
	var s strings.Builder
	s.WriteString(fmt.Sprintln("    type:", self.Type))
	s.WriteString(fmt.Sprintf("    mode:%o\n", self.Mode))
	s.WriteString(fmt.Sprintln("    nlink:", self.NLink))
	s.WriteString(fmt.Sprintln("    uid:", self.UID))
	s.WriteString(fmt.Sprintln("    gid:", self.GID))
	s.WriteString(fmt.Sprintln("    size:", self.Size))
	s.WriteString(fmt.Sprintln("    block size:", self.BlockSize))
	s.WriteString(fmt.Sprintln("    rdev:", self.Rdev))
	s.WriteString(fmt.Sprintln("    fsid:", self.FsId))
	s.WriteString(fmt.Sprintln("    fileid:", self.FileId))
	s.WriteString(fmt.Sprintln("    atime:", self.Atime))
	s.WriteString(fmt.Sprintln("    mtime:", self.Mtime))
	s.WriteString(fmt.Sprintln("    ctime:", self.Ctime))
	return s.String()

}

type Sattr struct {
	Mode         uint32 //mode
	UID          uint32 //文件的所有者的用户标识号码
	GID          uint32 //文件的组的组标识号码
	Size         uint32 //以字节数计算的大小
	Atime, Mtime TimeVal
}

func (self Sattr) String() string {
	var s strings.Builder
	s.WriteString(fmt.Sprintf("    mode:%o\n", self.Mode))
	s.WriteString(fmt.Sprintln("    uid:", self.UID))
	s.WriteString(fmt.Sprintln("    gid:", self.GID))
	s.WriteString(fmt.Sprintln("    size:", self.Size))
	s.WriteString(fmt.Sprintln("    atime:", self.Atime))
	s.WriteString(fmt.Sprintln("    mtime:", self.Mtime))
	return s.String()

}

type DirOpArgs struct {
	FileHandle []byte `xdr:"limit=32"`
	FileName   []byte
}

func (self DirOpArgs) FileNameStr() string {
	encoding := charset.GetStrCoding(self.FileName)
	switch encoding {
	case charset.GBK:
		utf8Data, _ := simplifiedchinese.GBK.NewDecoder().Bytes(self.FileName) //将gbk再转换为utf-8
		return string(utf8Data)
	}
	return string(self.FileName)
}

func (self DirOpArgs) String() string {
	var s strings.Builder
	s.WriteString(fmt.Sprintln("handle:", hex.EncodeToString(self.FileHandle)))
	encoding := charset.GetStrCoding(self.FileName)
	switch encoding {
	case charset.GBK:
		utf8Data, _ := simplifiedchinese.GBK.NewDecoder().Bytes(self.FileName) //将gbk再转换为utf-8
		s.WriteString(fmt.Sprintln("name:", string(utf8Data)))
	default:
		s.WriteString(fmt.Sprintln("name:", string(self.FileName)))
	}
	return s.String()
}

type DirOpRes struct {
	FileHandle []byte `xdr:"limit=32"`
	Attr       Fattr
}
