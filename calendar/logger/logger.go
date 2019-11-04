package logger

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var L *zap.SugaredLogger

func ConfigureLogger() {
	// The zap.Config struct includes an AtomicLevel. To use it, keep a
	// reference to the Config.
	rawJSON := []byte(`{
		"level": "info",
		"outputPaths": ["stdout"],
		"errorOutputPaths": ["stderr"],
		"encoding": "console",
		"encoderConfig": {
			"messageKey": "msg",
			"levelKey": "level",
			"levelEncoder": "lowercase",
			"timeEncoder": "ISO8601"
		}
	}`)

	var cfg zap.Config
	if err := json.Unmarshal(rawJSON, &cfg); err != nil {
		panic(err)
	}

	logLevel := viper.GetString("log_level")
	switch logLevel {
	case "debug":
		cfg.Level.SetLevel(zap.DebugLevel)
	case "info":
		cfg.Level.SetLevel(zap.InfoLevel)
	case "warn":
		cfg.Level.SetLevel(zap.WarnLevel)
	case "error":
		cfg.Level.SetLevel(zap.ErrorLevel)
	default:
		panic(fmt.Errorf("Unrecognized log_level: %s\n", logLevel))
	}

	// configure log_file
	logFile := viper.GetString("log_file")
	cfg.OutputPaths = append(cfg.OutputPaths, logFile)
	cfg.ErrorOutputPaths = append(cfg.ErrorOutputPaths, logFile)
	logger, err := cfg.Build()
	if err != nil {
		panic(err)
	}
	L = logger.Sugar()
	L.Info("logging configured")
}
