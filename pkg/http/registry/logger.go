package registry

import (
	ILog "examples/pkg/http/entity/infra/log"
	"examples/pkg/http/infra/log"
	"os"
)

var Logger ILog.Logger

func init() {
	Logger = log.InitLogger(os.Stdout)
}
