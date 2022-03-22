package database

import (
	"fmt"
	"os"

	"gopkg.in/ini.v1"
)

type Config struct {
	Port           int
	Host           string
	DoDefaultRoute int
}

func ReadConfiguration() Config {
	var c Config

	cfg, err := ini.Load("./db/config.ini")
	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
		os.Exit(1)
	}

	c.Port = cfg.Section("ads").Key("port").MustInt()
	c.Host = cfg.Section("ads").Key("host").String()
	c.DoDefaultRoute = cfg.Section("ads").Key("do_default_route").MustInt()
	return c
}
