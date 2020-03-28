package main

import (
	"errors"
	"github.com/pxecore/pxecore/pkg/ipxe"
	"github.com/pxecore/pxecore/pkg/ipxe/script"
	"github.com/pxecore/pxecore/pkg/tftp"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"os"
	"strings"
	"time"
)

type coreConfig struct {
	configFile     string
	debug          bool
	logFile        string
	singleMode     bool
	singleModeFile string
	ipxeBiosFile   string
	ipxeUEFIFile   string
}

var config = new(coreConfig)
var tftpServer = new(tftp.Server)

func main() {
	log.Info("Starting PXECORE Server...")
	loadCoreConfig(config)
	loadLogging(config)
	loadConfigFile(config)
	log.WithField("config", viper.AllSettings()).Debug("Config Loaded.")

	if config.ipxeBiosFile != "" {
		if err := ipxe.LoadIPXEBiosFile(config.ipxeBiosFile); err != nil {
			log.WithError(err).Fatal("Error loading ipxe-bios file.")
		}
	}
	if config.ipxeUEFIFile != "" {
		if err := ipxe.LoadIPXEUEFIFile(config.ipxeUEFIFile); err != nil {
			log.WithError(err).Fatal("Error loading ipxe-uefi file.")
		}
	}

	var ipxeScript tftp.IPXEScript
	if config.singleMode {
		i, err := script.NewSingleIPXEScriptFromFile(config.singleModeFile)
		if err != nil {
			log.WithError(err).Fatal("Error loading single file.")
		}
		ipxeScript = i
	}
	tftpServer = new(tftp.Server)
	tftpServer.Start(tftp.ServerConfig{
		Address:    ":69",
		Timeout:    5 * time.Second,
		IPXEScript: &ipxeScript,
	})
}

// loadCoreConfig defines the flags and environment used by the server.
func loadCoreConfig(c *coreConfig) {
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

	c.configFile = viper.GetString("config")
	c.debug = viper.GetBool("debug")
	c.logFile = viper.GetString("logfile")
	c.singleMode = viper.GetString("single") != ""
	c.singleModeFile = viper.GetString("single")
	c.ipxeBiosFile = viper.GetString("ipxe-bios")
	c.ipxeUEFIFile = viper.GetString("ipxe-uefi")
}

// loadLogging reads from flags and env variables the logging level and file.
func loadLogging(config *coreConfig) {
	if config.debug {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}

	if config.logFile != "" {
		file, err := os.OpenFile(config.logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err == nil {
			log.SetOutput(file)
			log.WithField("logfile", config.logFile).Debug("Logging into File")
		} else {
			log.WithField("logfile", config.logFile).Warn("Failed to open logfile, using stdout")
		}
	} else {
		log.Debug("Logging into STDERR")
	}
}

// loadConfigFile load de config files relative to the current path or on the config flag.
func loadConfigFile(config *coreConfig) {
	if config.configFile != "" {
		viper.SetConfigFile(config.configFile)
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
