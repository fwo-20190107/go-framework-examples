package config

import (
	"examples/pkg/util"
	"fmt"
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml/v2"
)

var C config

func LoadConfig() {
	path := filepath.Join(util.RootDir(), "config", "config.toml")
	file, err := os.Open(path)
	if err != nil {
		panic("failed load config")
	}
	defer file.Close()

	decoder := toml.NewDecoder(file)
	if err := decoder.Decode(&C); err != nil {
		panic(fmt.Sprintf("failed decode: %v\n", err))
	}
}

type config struct {
	DB struct {
		User   string `toml:"user"`
		Passwd string `toml:"passwd"`
		Net    string `toml:"net"`
		Addr   string `toml:"addr"`
		DBName string `toml:"dbname"`
	}
}
