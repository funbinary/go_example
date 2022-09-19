package proc_netaudit

import (
	"github.com/funbinary/go_example/pkg/blog"
)

//  当前程序名
const PROC_NAME = "bysc_netaudit"

func InitLog(path string, level string) {
	opts := &blog.Options{
		Level:             level,
		Format:            "console",
		EnableColor:       false,
		DisableCaller:     false,
		DisableStacktrace: false,
		OutputPath:        path,
		ModuleName:        PROC_NAME,
		MaxSize:           100,
		TimeFormat:        "2006-01-02",
		// ErrorOutputPaths: []string{"error.log"},
	}
	blog.Init(opts)

}
