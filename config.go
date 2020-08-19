package test

import (
	"log"
	"os/user"
	"path/filepath"
	"sync"

	"github.com/spf13/viper"
)

var Conf *viper.Viper
var once sync.Once

// Config 加载测试私密文件
func Config() *viper.Viper {
	once.Do(func() {
		u, err := user.Current()
		if err != nil {
			panic(err)
		}
		Conf = viper.New()
		file := filepath.Join(u.HomeDir, ".testing.yml")
		Conf.SetConfigFile(file)
		log.Println("[Load Testing Config]", file)
		if err := Conf.ReadInConfig(); err != nil {
			panic(err)
		}
	})
	return Conf
}
