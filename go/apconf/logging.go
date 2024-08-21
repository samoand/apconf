package apconf

import (
	"os"
	"path/filepath"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func SetupLogging(loggingConfig map[string]any) (*zap.Logger, error) {
	// Set the log level
	level := zapcore.InfoLevel
	if lvl, ok := loggingConfig["level"].(string); ok {
		var err error
		level, err = zapcore.ParseLevel(lvl)
		if err != nil {
			return nil, err
		}
	}

	// Set the log format
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "time"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	var encoder zapcore.Encoder

	if formatter, ok := loggingConfig["formatter"].(string); ok && formatter == "json" {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}

	// Set the output destination (file or stdout)
	var output zapcore.WriteSyncer
	if outputFile, ok := loggingConfig["output_file"].(string); ok {
		file, err := os.OpenFile(filepath.Clean(outputFile), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return nil, err
		}
		output = zapcore.AddSync(file)
	} else {
		output = zapcore.AddSync(os.Stdout)
	}

	core := zapcore.NewCore(encoder, output, level)

	// Create the logger
	logger := zap.New(core)

	// Replace the global logger with the new logger if desired
	zap.ReplaceGlobals(logger)

	return logger, nil
}
