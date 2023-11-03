package config

import (
	"github.com/aioncore/ouroboros/pkg/service/utils"
	"gopkg.in/yaml.v2"
	"path/filepath"
	"reflect"
)

func InitConfig(filePath string, serviceType string, defaultConfig interface{}, configType reflect.Type) interface{} {
	config := reflect.New(configType)
	exist, err := utils.PathExist(filepath.Join(filePath, serviceType, "config", "core.yaml"))
	if err != nil {
		panic(err)
	}
	if exist {
		file, err := utils.ReadFile(filepath.Join(filePath, serviceType, "config", "core.yaml"))
		if err != nil {
			panic(err)
		}
		err = yaml.Unmarshal(file, config.Interface())
		if err != nil {
			panic(err)
		}
		return config.Elem().Interface()
	}
	jsonBytes, err := yaml.Marshal(defaultConfig)
	if err != nil {
		panic(err)
	}
	err = utils.WriteFile(filepath.Join(filePath, serviceType, "config", "core.yaml"), jsonBytes, 0666)
	if err != nil {
		panic(err)
	}
	return defaultConfig
}
