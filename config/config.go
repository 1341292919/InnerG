package config

import (
	"errors"
	"github.com/bytedance/gopkg/util/logger"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var (
	Mysql        *mySQL
	Redis        *redis
	Smtp         *smtp
	Service      *service
	Api          *api
	Oss          *oss
	MongoDb      *mongodb
	runtimeViper = viper.New()
)

const (
	File     = "./config/config.yaml"
	FileType = "yaml"
)

func Init() {
	runtimeViper.SetConfigFile(File)
	runtimeViper.SetConfigType(FileType)
	if err := runtimeViper.ReadInConfig(); err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if errors.As(err, &configFileNotFoundError) {
			logger.Fatal("config.Init: could not find config files")
		}
		logger.Fatalf("config.Init: read config error: %v", err)
	}
	configMapping()
	// 设置持续监听
	runtimeViper.OnConfigChange(func(e fsnotify.Event) {
		logger.Infof("config: notice config changed: %v\n", e.String())
		configMapping()
	})
	runtimeViper.WatchConfig()
}

func configMapping() {
	c := new(config)
	if err := runtimeViper.Unmarshal(&c); err != nil {
		// 由于这个函数会在配置重载时被再次触发，所以需要判断日志记录方式
		logger.Fatalf("config.configMapping: config: unmarshal error: %v", err)
	}
	Mysql = &c.MySQL
	Redis = &c.Redis
	Smtp = &c.Smtp
	Oss = &c.OSS
	Service = &c.Service
	Api = &c.Api
	MongoDb = &c.MongoDb
}
