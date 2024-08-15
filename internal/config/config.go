package config

import (
	"os"
	"strings"
	"time"

	"github.com/spf13/viper"
)

var conFile string

type Config struct {
	Keenetic      Keenetic      `mapstructure:"keenetic"`
	Mqtt          Mqtt          `mapstructure:"mqtt"`
	Homeassistant HomeAssistant `mapstructure:"homeassistant"`
}

type Keenetic struct {
	Host     string `mapstructure:"host"`
	Login    string `mapstructure:"login"`
	Password string `mapstructure:"password"`
}

type Mqtt struct {
	Host      string `mapstructure:"host"`
	Login     string `mapstructure:"login"`
	Password  string `mapstructure:"password"`
	ClientID  string `mapstructure:"clientId"`
	BaseTopic string `mapstructure:"baseTopic"`
}

type HomeAssistant struct {
	UpdateInterval time.Duration `mapstructure:"updateInterval"`
	WhiteList      []string      `mapstructure:"whitelist"`
	DeviceID       string        `mapstructure:"deviceid"`
}

func SetConfigFile(path string) {
	conFile = path
}

func NewDefaultConfig() (*Config, error) {
	if err := InitializeConfig(); err != nil {
		return nil, err
	}
	config := Config{}
	err := viper.Unmarshal(&config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func InitializeConfig() error {
	path, err := os.Getwd()
	if err != nil {
		return err
	}

	if conFile != "" {
		viper.SetConfigFile(strings.TrimSpace(conFile))
	} else {
		viper.SetConfigFile(strings.TrimSpace(path + "/configs/config.yml"))
	}
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	if err = viper.ReadInConfig(); err != nil {
		return err
	}
	return nil
}
