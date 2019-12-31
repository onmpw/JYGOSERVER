package main

import (
	"JYGO/config"
	"flag"
)

const (
	// CLIENT_MAX 表示客户端队列中最高允许的客户端连接数
	DateFormat  = "2006-01-02 15:04:05"
	CLIENT_MAX = 1024
	HOST = "localhost"
	PORT = "9002"
	NETWORK = "tcp"
)

var (
	configFile = "/etc/jygoserver.ini"
	Host string
	Port string
)

// loadConfig 解析配置项
func loadConfig() {
	file := flag.String("-c",configFile,"通过-c指定配置文件 如果没有指定则使用默认的配置文件")
	err := config.Init(*file)
	if err != nil {
		panic(err.Error())
	}
	Host = config.Conf.C("ListenAddr")
	Port = config.Conf.C("ListenPort")

	if Host == "" {
		Host = HOST
	}

	if Port == "" {
		Port = PORT
	}
}

