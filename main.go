package main

import (
	"errors"
	"github.com/pxecore/pxecore/pkg/controller"
	"github.com/pxecore/pxecore/pkg/http"
	repo "github.com/pxecore/pxecore/pkg/repository"
	"github.com/pxecore/pxecore/pkg/tftp"
	"github.com/pxecore/pxecore/pkg/tftp/locator"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"os"
	"strings"
	"time"
)

var tftpServer *tftp.Server
var repository repo.Repository

func main() {
	log.Info("Starting PXECORE Server...")

	loadDefaultConfig()
	loadCoreConfig()
	loadConfigFile()
	loadLogging()
	log.WithField("config", viper.AllSettings()).Debug("Config Loaded.")

	r, err := repo.NewRepository(viper.GetStringMap("db"))
	if err != nil {
		log.WithError(err).Fatal("Error loading repository server.")
	}
	repository = r

	tftpServer = new(tftp.Server)
	if err := tftpServer.StartInBackground(tftp.ServerConfig{
		Address: viper.GetString("tftp.address"),
		Timeout: viper.GetDuration("tftp.timeout"),
		FileLocators: []tftp.FileLocator{
			locator.NewIPXEFirmware(),
			locator.NewRepositoryIPXEScript(repository),
		},
	}); err != nil {
		log.Fatal(err)
	}

	s := http.Server{Controllers: []http.Controller{
		controller.Template{Repository: repository},
		controller.Host{Repository: repository},
		controller.Group{Repository: repository},
	}}
	c, err := http.NewConfig(viper.GetStringMap("http"))
	if err != nil {
		log.Fatal(err)
	}
	log.Fatal(s.Start(c))
}

// loadDefaultConfig loads default config.
func loadDefaultConfig() {
	viper.SetDefault("tftp", map[string]interface{}{
		"address": ":69",
		"timeout": 5 * time.Second,
	})
	viper.SetDefault("http", map[string]interface{}{
		"address":       ":80",
		"read-timeout":  10,
		"write-timeout": 10,
	})
	viper.SetDefault("db", map[string]interface{}{
		"driver": "memory",
	})
}

// loadCoreConfig defines the flags and environment used by the server.
func loadCoreConfig() {
	viper.SetEnvPrefix("pxecore")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	pflag.ErrHelp = errors.New("pxecore-server: help requested")
	pflag.StringP("config", "c", "", "Config file path.")
	viper.BindEnv("config")
	pflag.Bool("debug", false, "Verbose Output.")
	viper.BindEnv("debug")
	pflag.StringP("logfile", "l", "", "Log file Path.")
	viper.BindEnv("logfile")
	pflag.StringP("basedir", "b", "", "Static file directory.")
	viper.BindEnv("basedir")
	pflag.Parse()
	if err := viper.BindPFlags(pflag.CommandLine); err != nil {
		log.Warn("Error reading flags: ", err)
	}
}

// loadLogging reads from flags and env variables the logging level and file.
func loadLogging() {
	if viper.GetBool("debug") {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}

	if lf := viper.GetString("logfile"); lf != "" {
		file, err := os.OpenFile(lf, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err == nil {
			log.SetOutput(file)
			log.WithField("logfile", lf).Debug("Logging into File")
		} else {
			log.WithField("logfile", lf).Warn("Failed to open logfile, using stdout")
		}
	} else {
		log.Debug("Logging into STDERR")
	}
}

// loadConfigFile load de config files relative to the current path or on the config flag.
func loadConfigFile() {
	if cf := viper.GetString("config"); cf != "" {
		viper.SetConfigFile(cf)
	} else {
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(".")
	}

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			log.Fatal("Error loading viper config", err)
		}
	}
}
