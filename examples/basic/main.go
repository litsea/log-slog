package main

import (
	"bytes"
	"embed"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/spf13/viper"

	log "github.com/litsea/log-slog"
)

//go:embed config.yaml
var config embed.FS

func main() {
	data, err := config.ReadFile("config.yaml")
	if err != nil {
		fmt.Println("main: read config, ", err)
		os.Exit(1)
	}

	viper.SetConfigType("yaml")
	err = viper.ReadConfig(bytes.NewReader(data))
	if err != nil {
		fmt.Println("main: load config, ", err)
		os.Exit(1)
	}

	logCfg := viper.Sub("log")

	err = log.Set(logCfg, log.WithVersion("0.1.0-dev"), log.WithAddSource(true))
	if err != nil {
		fmt.Println("main: set log, ", err)
		os.Exit(1)
	}

	log.With("time2", time.Now()).With("foo", "bar").Debug("debug")
	log.Info("info")
	log.Warn("warn", "err", fmt.Errorf("warn"))
	log.Error("error", "err", fmt.Errorf("test err: %w", errors.New("foo err")))

	log.WithGroup("group1").With("foo", "bar").Debug("debug with group and with")

	// Wait for send event
	time.Sleep(time.Second * 5)
}
