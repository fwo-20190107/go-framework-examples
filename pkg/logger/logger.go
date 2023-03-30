package logger

import (
	ILog "examples/pkg/entity/infra/log"
	"examples/pkg/infra/log"
	"os"
)

var L ILog.Logger

func init() {
	L = log.InitLogger(os.Stdout)
}
