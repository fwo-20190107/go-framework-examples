package registry

import (
	ILog "examples/internal/http/entity/infra/log"
	"examples/internal/http/infra/log"
	"os"
)

var Logger ILog.Logger

func init() {
	Logger = log.InitLogger(os.Stdout)
}
