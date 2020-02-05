package config

import (
	"flag"
	"fmt"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Config struct {
	Env         string      `yaml:"env"`
	Mysql       Mysql       `yaml:"mysql"`
	Redis       Redis       `yaml:"redis"`
	AmapKey     string      `yaml:"amapKey"`
	AliYun      AliYun      `yaml:"aliYun"`
}

type Mysql struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
}

type Redis struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Password string `yaml:"password"`
	Db       int    `yaml:"db"`
}

type AliYun struct {
	AccessKeyID     string `yaml:"accessKeyID"`
	AccessKeySecret string `yaml:"accessKeySecret"`
}

var Conf *Config

func init() {
	var confPath = "release/config/production.yaml"
	flag.StringVar(&confPath, "c", confPath, "runtime config yaml file")
	flag.Parse()
	logrus.Warn(fmt.Sprintf("config file: %s", confPath))

	bytes, err := ioutil.ReadFile(confPath)
	if err != nil {
		logrus.Error(err)
		return
	}
	err = yaml.Unmarshal([]byte(bytes), &Conf)
	if err != nil {
		logrus.Error(err)
	}
}
