package config

import (
	"bytes"
	_ "embed"
	"log"
	"strings"

	"github.com/spf13/viper"
)

type SettingStruct struct {
	MysqlConfig struct {
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		DataBase string `yaml:"database"`
		UserName string `yaml:"username"`
		PassWord string `yaml:"password"`
		Charset  string `yaml:"charset"`
	}
}

var Setting = new(SettingStruct)

var (
	V = viper.New()

	//go:embed config.yaml
	configFile []byte
)

func init() {
	V.SetConfigName("config2")
	V.SetConfigType("yaml")
	V.AddConfigPath("./config")
	V.AddConfigPath(".")

	V.SetEnvPrefix("VIPER")
	V.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	V.AutomaticEnv()

	var err error
	if err = V.ReadInConfig(); err != nil {
		if err = V.ReadConfig(bytes.NewBuffer(configFile)); err != nil {
			log.Fatalf("fatal error with reading embed config: %s", err)
		}
	}

	if err = V.Unmarshal(&Setting); err != nil {
		log.Fatalln(err)
	}
	log.Println("Init Config Complate")
}
