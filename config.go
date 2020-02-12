package main

import (
	"flag"
	"github.com/onmpw/JYGO/config"
)

const (
	DateFormat = "2006-01-02 15:04:05"
	HOST       = "localhost"
	PORT       = "9002"
	NETWORK    = "tcp"
	AUTH       = "off"
)

var (
	configFile = "/etc/jygoserver.ini"
	Host       string
	Port       string
	Auth       string
)

// loadConfig 解析配置项
func loadConfig() {
	file := flag.String("-c", configFile, "通过-c指定配置文件 如果没有指定则使用默认的配置文件")
	err := config.Init(*file)
	if err != nil {
		panic(err.Error())
	}
	Host = config.Conf.C("ListenAddr")
	Port = config.Conf.C("ListenPort")

	Auth = config.Conf.C("Auth")

	if Auth == "" {
		Auth = AUTH
	}

	if Host == "" {
		Host = HOST
	}

	if Port == "" {
		Port = PORT
	}
}

func checkAuth() bool {
	if Auth == "off" {
		return true
	}

	return false
}
