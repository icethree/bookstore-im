package config

import (
	"flag"
	"time"

	"github.com/go-redis/redis"
)

var (
	confFile = flag.String("c", "config.toml", "-c conf file path")
	level    = flag.Int("l", -1, "logs level, -1:debug, 0:info, 1:warn, 2:error")
)

type Conf struct {
	Listen string
	ENV    string `toml:"env"`
	MySql  struct {
		DSN         string
		Active      int
		Idle        int
		IdleTimeout time.Duration
	} `toml:"mysql"`
	Logger struct {
		Path         string
		NormalPrefix string
		SqlPrefix    string
		AccessPrefix string
	} `toml:"logger"`
	Redis *redis.UniversalOptions `toml:"redis"`
}

var Config = new(Conf)

//func (c *Conf) load(file string) {
//	if _, err := toml.DecodeFile(file, c); err != nil {
//		panic("service file read failed, err:" + err.Error())
//	}
//}
//
//func init() {
//	flag.Parse()
//
//	Config.load(*confFile)
//
//	log.Init(Config.Logger.Path, Config.Logger.NormalPrefix, zapcore.Level(*level))
//}
