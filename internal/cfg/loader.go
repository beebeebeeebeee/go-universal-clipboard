package cfg

import (
	"github.com/spf13/viper"
	"log"
)

func LoadEnv(path string) {
	v := viper.New()
	v.AddConfigPath(path)
	v.SetConfigType("yaml")
	v.SetConfigName("")

	if err := v.ReadInConfig(); err != nil {
		log.Fatalf("error on reading config file, %s", err)
	}

	err := v.Unmarshal(&Cfg)
	if err != nil {
		log.Fatalf("unable to decode into config, %s", err)
	}

}
