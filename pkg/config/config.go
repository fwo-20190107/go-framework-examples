package config

import (
	"examples/pkg/code"
	"examples/pkg/errors"
	"examples/pkg/util"
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml/v2"
)

var C config

func LoadConfig() error {
	path := filepath.Join(util.RootDir(), "config.toml")
	file, err := os.Open(path)
	if err != nil {
		return errors.Errorf(code.CodeInternal, "failed load config: %s\n", path)
	}
	defer file.Close()

	decoder := toml.NewDecoder(file)
	if err := decoder.Decode(&C); err != nil {
		return errors.Errorf(code.CodeInternal, "failed decode: %v\n", err)
	}
	return nil
}

type config struct {
	Server struct {
		Addr int `toml:"addr"`
	}
	DB struct {
		User   string `toml:"user"`
		Passwd string `toml:"passwd"`
		Net    string `toml:"net"`
		Addr   string `toml:"addr"`
		DBName string `toml:"dbname"`
	}
	// DB struct {
	// 	Schema   string `toml:"schema"`
	// 	Testdata string `toml:"testdata"`
	// }
}
