package config

import(
	"github.com/spf13/viper"
	"os"
)

var conf *viper.Viper

func Get() *viper.Viper {
	if conf == nil {
		path, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		conf = viper.New()
		conf.AddConfigPath(path)
		conf.SetConfigType("yaml")
		conf.SetConfigName("config")
		err = conf.ReadInConfig()
		if err != nil {
			panic(err)
		}
	}
	return conf
}
