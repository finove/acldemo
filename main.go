package main

import (
	"io"
	"os"
	"strconv"

	"github.com/finove/acldemo/server"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"gopkg.in/natefinch/lumberjack.v2"
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
		SetupZerolog(cmdLevel, cmdLogFile, cmdVerbose)
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

func SetupZerolog(logLevel, logFile string, verbose bool) {
	var iowrites []io.Writer
	if verbose {
		var consoleWriter = zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: "2006/01/02 15:04:05.000"}
		iowrites = append(iowrites, consoleWriter)
	}
	if logFile != "" {
		var logfileWriter = &lumberjack.Logger{
			Filename:   logFile,
			MaxSize:    100,
			MaxBackups: 14,
			MaxAge:     14,
			LocalTime:  true,
			Compress:   false,
		}
		logfileWriter.Rotate()
		iowrites = append(iowrites, logfileWriter)
	}
	if len(iowrites) > 0 {
		var multi = zerolog.MultiLevelWriter(iowrites...)
		zerolog.SetGlobalLevel(getZerologLevel(logLevel))
		zerolog.TimeFieldFormat = "2006/01/02 15:04:05.000"
		zerolog.CallerMarshalFunc = func(file string, line int) string {
			short := file
			for i := len(file) - 1; i > 0; i-- {
				if file[i] == '/' {
					short = file[i+1:]
					break
				}
			}
			file = short
			return file + ":" + strconv.Itoa(line)
		}
		log.Logger = log.With().Caller().Logger().Output(multi)
	} else {
		log.Logger = zerolog.Nop()
	}
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
