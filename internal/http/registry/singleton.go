package registry

import (
	"examples/internal/http/entity/infra"
	"examples/internal/http/infra/log"
	"os"
)

var Logger infra.Logger

func init() {
	Logger = log.InitLogger(os.Stdout)
}
