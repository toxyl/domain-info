package main

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	ConcurrentRequests int `mapstructure:"concurrent_requests"`
	DNSServers         []struct {
		Host    string `mapstructure:"host"`
		Port    uint   `mapstructure:"port"`
		Timeout uint   `mapstructure:"timeout"`
	} `mapstructure:"dns_servers"`
}

var cfgFile string
var Conf Config

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath("/etc/domain-info/")
		viper.AddConfigPath("$HOME/.domain-info")
		viper.AddConfigPath(".")
	}

	err := viper.ReadInConfig()
	if err != nil {
		log.Panic(fmt.Errorf("[Config] Fatal error config file: %w", err))
	}

	err = viper.Unmarshal(&Conf)
	if err != nil {
		log.Printf("[Config] Unable to decode into Config struct, %v", err)
	}
}
