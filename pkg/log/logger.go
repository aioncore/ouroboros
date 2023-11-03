package log

import (
	"github.com/aioncore/ouroboros/pkg/service/log"
)

func InitLogger(filePath string, serviceType string) {
	log.InitLogger(filePath, serviceType)
}
