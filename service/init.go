package service

import (
	"github.com/sirupsen/logrus"
	"github.com/killsari/api/config"
)

func init() {
	if config.Conf.Env == "prod" {
		logrus.Info("Service start .....")
		//todo
	}
}
