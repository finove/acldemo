package main

import (
	"os"

	"github.com/finove/acldemo/server"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var (
	cmdVerbose    bool
	cmdLevel      string
	cmdLogFile    string
	cmdConfigFile string
)

func main() {
	root.Execute()
}

var root cobra.Command = cobra.Command{
	Use:     "acldemo",
	Version: "0.0.1",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		zerolog.SetGlobalLevel(getZerologLevel(cmdLevel))
		zerolog.TimeFieldFormat = "2006/01/02 15:04:05.000"
		if cmdVerbose {
			log.Logger = log.With().Caller().Logger().Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: "2006/01/02 15:04:05.000"})
		} else {
			log.Logger = log.Output(os.Stdout)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		log.Error().Msg("error info")
		log.Warn().Msg("warning info")
		log.Info().Msg("info info")
		log.Debug().Msg("debug info")
		server.Run()
	},
}

func init() {
	root.PersistentFlags().BoolVarP(&cmdVerbose, "verbose", "v", false, "show log in console")
	root.PersistentFlags().StringVarP(&cmdLevel, "level", "l", "info", "log level (debug|info|warn|error|fatal)")
	root.PersistentFlags().StringVar(&cmdLogFile, "logpath", "", "log file path")
	root.PersistentFlags().StringVar(&cmdConfigFile, "config", "", "use config file")
}

func getZerologLevel(level string) (l zerolog.Level) {
	switch level {
	case "fatal":
		l = zerolog.FatalLevel
	case "error":
		l = zerolog.ErrorLevel
	case "warn":
		l = zerolog.WarnLevel
	case "info":
		l = zerolog.InfoLevel
	case "debug":
		l = zerolog.DebugLevel
	default:
		l = zerolog.Disabled
	}
	return
}
