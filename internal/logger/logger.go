package logger

import (
	"github.com/BroNaz/proxy/internal/config"
	"github.com/mattn/go-isatty"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"log/syslog"
	"os"
	"time"
)

func SetupLogging(cfg *config.LogConfig) {
	zerolog.DurationFieldUnit = time.Microsecond
	zerolog.TimeFieldFormat = "2006-01-02T15:04:05.000000-07:00"

	level := logLevel(cfg.Level)
	log.Logger = log.Logger.
		Level(level)

	switch cfg.Output {
	case "console", "":
		log.Logger = log.Output(zerolog.ConsoleWriter{
			Out:        os.Stdout,
			NoColor:    !isatty.IsTerminal(os.Stdout.Fd()),
			TimeFormat: "2006-01-02 15:04:05.000",
		})
	case "syslog":
		syswr, err := syslog.New(syslog.LOG_DEBUG|syslog.LOG_LOCAL0, "proxy")
		if err != nil {
			log.Fatal().Err(err).Msg("failed to connect to syslog")
		}
		log.Logger = log.Output(zerolog.SyslogLevelWriter(syswr))
	default:
		log.Fatal().Str("name", cfg.Output).Msg("unknown Log.Output in config")
	}
}

func logLevel(level string) zerolog.Level {
	switch level {
	case "debug":
		return zerolog.DebugLevel
	case "info":
		return zerolog.InfoLevel
	case "warn":
		return zerolog.WarnLevel
	case "error":
		return zerolog.ErrorLevel
	case "fatal":
		return zerolog.FatalLevel
	case "panic":
		return zerolog.PanicLevel
	default:
		return zerolog.InfoLevel
	}
}
