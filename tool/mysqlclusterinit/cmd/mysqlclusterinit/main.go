package main

import (
	"fmt"
	"strings"

	"github.com/bingoohuang/gonet"

	"github.com/BurntSushi/toml"
	"github.com/bingoohuang/tool/mysqlclusterinit"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
)

func main() {
	checkmysql := pflag.BoolP("checkmysql", "m", false, "check mysql")
	ver := pflag.BoolP("version", "v", false, "show version")
	testaddr := pflag.StringP("testaddr", "t", "", "test addr is local or not")
	conf := pflag.StringP("config", "c", "./config.toml", "config file path")
	pflag.Parse()

	args := pflag.Args()
	if len(args) > 0 {
		fmt.Printf("Unknown args %s\n", strings.Join(args, " "))
		pflag.PrintDefaults()

		return
	}

	if *ver {
		fmt.Printf("Version: 1.2\n")
		return
	}

	if *testaddr != "" {
		yes, _ := gonet.IsLocalAddr(*testaddr)
		if yes {
			fmt.Printf("%s is a local address\n", *testaddr)
		} else {
			fmt.Printf("%s is a non-local address\n", *testaddr)
		}

		return
	}

	configFile, _ := homedir.Expand(*conf)
	settings := mustLoadConfig(configFile)

	if *checkmysql {
		settings.CheckMySQL()
		return
	}

	if r := settings.InitMySQLCluster(); r.Error != nil {
		logrus.Errorf("error %v", r.Error)
	}
}

func loadConfig(configFile string) (config mysqlclusterinit.Settings, err error) {
	if _, err = toml.DecodeFile(configFile, &config); err != nil {
		logrus.Errorf("DecodeFile error %v", err)
	}

	return
}

func mustLoadConfig(configFile string) (config mysqlclusterinit.Settings) {
	var err error
	if config, err = loadConfig(configFile); err != nil {
		logrus.Panic(err)
	}

	if config.Port <= 0 {
		config.Port = 3306
	}

	logrus.Debugf("config: %+v\n", config)

	return config
}