package config

import (
	"github.com/aioncore/ouroboros/pkg/service/config"
	"path/filepath"
	"reflect"
)

var Core *CoreConfig

type CoreConfig struct {
	NodeKeyPath string
}

func InitConfig(filePath string, serviceType string) {
	Core = config.InitConfig(
		filePath,
		serviceType,
		DefaultCoreConfig(filePath),
		reflect.TypeOf(Core),
	).(*CoreConfig)
}

func DefaultCoreConfig(filePath string) interface{} {
	return &CoreConfig{
		NodeKeyPath: filepath.Join(filePath, "keyt", "node_key.json"),
	}
}
