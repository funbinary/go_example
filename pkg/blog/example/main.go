package main

import (
	log "github.com/funbinary/go_example/pkg/blog"
)

func main() {
	// log.Debug("aaa")
	opts := &log.Options{
		Level:             "debug",
		Format:            "console",
		EnableColor:       true,
		DisableCaller:     false,
		DisableStacktrace: false,
		OutputPath:        "/kds/log",
		ModuleName:        "bysc_ftp",
		MaxSize:           10,
		TimeFormat:        "2006-01-02",
		// ErrorOutputPaths: []string{"error.log"},
	}
	log.Init(opts)
	log.Debug("asd")
	log.Info("asd")
	log.Warn("ffff")
	log.Error("fsafaf")

	// log.Panic("fsafaf")
	// lv := log.WithName("hello")
	// lv.Debug("afsf")
	// log.Debug("fff")

}
