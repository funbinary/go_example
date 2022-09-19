package config

import (
	"github.com/funbinary/go_example/pkg/bfile"
	"github.com/spf13/viper"
	"time"
)

type PathConfig struct {
	LogPath       string `mapstructure:"Log"`
	DumpPath      string `mapstructure:"Dump"`
	CacheNetaudit string `mapstructure:"CacheNetaudit"` // NFS缓存路径
	PcapDumpPath  string `mapstructure:"PcapDumpPath"`
}

type ZmqConfig struct {
	ZmqClient string `mapstructure:"Client"`
	ZmqId     string `mapstructure:"ModuleNetaudit"`
	ZmqUiId   string `mapstructure:"ModuleUi"`
}

type AuditConfig struct {
	TimeOut        int    `mapstructure:"TimeOut"`
	Snaplen        int    `mapstructure:"Snaplen"`
	Promisc        bool   `mapstructure:"Promisc"`
	BufferSize     int    `mapstructure:"BufferSize"`
	Interface      string `mapstructure:"Interface"`
	EnablePcapDump bool   `mapstructure:"EnablePcapDump"`
	BPF            string `mapstructure:"BPF"`
}

type DebugConfig struct {
	LogLevel int `mapstructure:"LogLevel"`
}

type IniFile struct {
	PathConfig  `mapstructure:"Path"`
	ZmqConfig   `mapstructure:"Zmq"`
	AuditConfig `mapstructure:"Audit"`
	DebugConfig `mapstructure:"Debug"`
}

var Config IniFile

func init() {
	//viper.SetDefault("PortConfig.Km1Com", "/dev/ttyUSB2")
	//viper.SetDefault("PortConfig.Km2Com", "/dev/ttyUSB3")
	//viper.SetDefault("PathConfig.LogPath", "/kds/log")
	//viper.SetDefault("ZmqConfig.ZmqClient", "tcp://localhost:10000")
	//viper.SetDefault("ZmqConfig.ModuleKm", "bysc_km")
	for {
		time.Sleep(100 * time.Millisecond)
		viper.SetConfigName("config")
		viper.SetConfigType("ini")
		viper.AddConfigPath(bfile.SelfDir())
		err := viper.ReadInConfig()
		if err != nil {
			continue
		}
		err = viper.Unmarshal(&Config)
		if err != nil {
			continue
		}
		break
	}

}

func Get(k string) interface{} {
	return viper.Get(k)
}

func GetString(k string) interface{} {
	return viper.GetString(k)
}
