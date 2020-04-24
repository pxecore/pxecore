package main

import (
	"errors"
	"github.com/pxecore/pxecore/pkg/controller"
	"github.com/pxecore/pxecore/pkg/http"
	"github.com/pxecore/pxecore/pkg/ipxe"
	"github.com/pxecore/pxecore/pkg/ipxe/script"
	repo "github.com/pxecore/pxecore/pkg/repository"
	"github.com/pxecore/pxecore/pkg/tftp"
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
	loadLogging()
	loadConfigFile()
	log.WithField("config", viper.AllSettings()).Debug("Config Loaded.")

	if err := overrideIPXEFiles(); err != nil {
		log.WithError(err).Fatal("Error loading ipxe file.")
	}

	if err := loadTFTPServer(); err != nil {
		log.WithError(err).Fatal("Error loading tftp server.")
	}

	r, err := repo.NewRepository(viper.GetStringMap("db"))
	if err != nil {
		log.WithError(err).Fatal("Error loading repository server.")
	}
	repository = r

	s := http.Server{Controllers: []http.Controller{
		controller.Template{Repository: repository},
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
	pflag.StringP("logfile", "l", "", "Logfile Path.")
	viper.BindEnv("logfile")
	pflag.StringP("single", "s", "", "Single Mode File Path.")
	viper.BindEnv("single")
	pflag.String("ipxe-bios", "", "Single Mode File Path.")
	viper.BindEnv("ipxe-bios")
	pflag.String("ipxe-uefi", "", "Single Mode File Path.")
	viper.BindEnv("ipxe-uefi")
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

// overrideIPXEFiles replaces the memory loaded IPXE File with the one provided.
func overrideIPXEFiles() error {
	if path := viper.GetString("ipxe-bios"); path != "" {
		if err := ipxe.LoadIPXEBiosFile(path); err != nil {
			return err
		}
	}
	if path := viper.GetString("ipxe-uefi"); path != "" {
		if err := ipxe.LoadIPXEUEFIFile(path); err != nil {
			return err
		}
	}
	return nil
}

// loadTFTPServer starts the TFTP Server.
func loadTFTPServer() error {
	tftpServer = new(tftp.Server)
	return tftpServer.StartInBackground(tftp.ServerConfig{
		Address:    viper.GetString("tftp.address"),
		Timeout:    viper.GetDuration("tftp.timeout"),
		IPXEScript: loadIPXEScript(),
	})
}

// loadIPXEScript checks the configuration and load the wanted IPXE resolver.
func loadIPXEScript() *tftp.IPXEScript {
	var i tftp.IPXEScript
	var err error
	if smf := viper.GetString("single"); smf != "" {
		i, err = script.NewSingleIPXEScriptFromFile(smf)
		if err != nil {
			log.WithError(err).Fatal("Error loading single file.")
		}
		return &i
	}
	return nil
}
