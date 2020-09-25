package utils

import (
	"github.com/jinzhu/configor"
	"go.uber.org/zap"
)

var Config = struct {
	CARoot string `required:"true"`
}{}

func LoadConfig() {
	err := configor.Load(&Config, "config.yml")
	if err != nil {
		Log().Fatal("load config failed", zap.Error(err))
	}
	Log().Info("config loaded", zap.Any("config", Config))
}
