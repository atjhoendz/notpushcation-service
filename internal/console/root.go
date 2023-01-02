package console

import (
	"fmt"
	"os"

	"github.com/atjhoendz/notpushcation-service/internal/config"
	runtime "github.com/banzaicloud/logrus-runtime-formatter"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// RootCmd :nodoc:
var RootCmd = &cobra.Command{
	Use: "root-cmd",
}

// Execute :nodoc:
func Execute() {
	err := RootCmd.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	config.GetConf()
	setupLogger()
}

func setupLogger() {
	formatter := runtime.Formatter{
		ChildFormatter: &log.JSONFormatter{},
		Line:           true,
		File:           true,
	}

	if config.Env() == "development" {
		formatter = runtime.Formatter{
			ChildFormatter: &log.TextFormatter{
				ForceColors:   true,
				FullTimestamp: true,
			},
			Line: true,
			File: true,
		}
	}

	log.SetFormatter(&formatter)
	log.SetOutput(os.Stdout)

	logLevel, err := log.ParseLevel(config.LogLevel())
	if err != nil {
		logLevel = log.DebugLevel
	}
	log.SetLevel(logLevel)
}
