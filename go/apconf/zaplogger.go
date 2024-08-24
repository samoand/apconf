package apconf

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// LogConfig represents the structure of the logging configuration.
type LogConfig struct {
	Level string              `json:"level"`
	Cores map[string]*LogCore `json:"cores"`
}

// LogCore represents the configuration for each logging core.
type LogCore struct {
	Level      string          `json:"level"`
	Encoding   string          `json:"encoding"`
	OutputPath string          `json:"outputPath,omitempty"`
	Rotation   *RotationConfig `json:"rotation,omitempty"`
	Encoder    map[string]any  `json:"encoder,omitempty"`
}

// RotationConfig represents the rotation settings for file logging.
type RotationConfig struct {
	MaxSize    int  `json:"maxSize"`
	MaxBackups int  `json:"maxBackups"`
	MaxAge     int  `json:"maxAge"`
	Compress   bool `json:"compress"`
}

// ConvertMapToLogConfig converts a map[string]any to a LogConfig struct.
func ConvertMapToLogConfig(config map[string]any) (*LogConfig, error) {
	logConfig := &LogConfig{}

	// Convert level
	if level, ok := config["level"].(string); ok {
		logConfig.Level = level
	} else {
		return nil, errors.New(
			"invalid or missing 'level' field in config")
	}

	// Convert cores
	if cores, ok := config["cores"].(map[string]any); ok {
		logConfig.Cores = make(map[string]*LogCore)
		for coreName, core := range cores {
			if coreMap, ok := core.(map[string]any); ok {
				logCore, err := convertMapToLogCore(coreMap)
				if err != nil {
					return nil, fmt.Errorf(
						"failed to convert core '%s': %w", coreName, err)
				}
				logConfig.Cores[coreName] = logCore
			} else {
				return nil, fmt.Errorf(
					"invalid core configuration format for '%s'", coreName)
			}
		}
	} else {
		return nil, errors.New(
			"invalid or missing 'cores' field in config")
	}

	return logConfig, nil
}

// convertMapToLogCore converts a map[string]any to a LogCore struct.
func convertMapToLogCore(core map[string]any) (*LogCore, error) {
	logCore := &LogCore{}

	// Convert level
	if level, ok := core["level"].(string); ok {
		logCore.Level = level
	} else {
		return nil, errors.New(
			"invalid or missing 'level' field in core config")
	}

	// Convert encoding
	if encoding, ok := core["encoding"].(string); ok {
		logCore.Encoding = encoding
	} else {
		return nil, errors.New(
			"invalid or missing 'encoding' field in core config")
	}

	// Convert output path (optional)
	if outputPath, ok := core["outputPath"].(string); ok {
		logCore.OutputPath = outputPath
	}

	// Convert rotation config (optional)
	if rotation, ok := core["rotation"].(map[string]any); ok {
		rotationConfig, err := convertMapToRotationConfig(rotation)
		if err != nil {
			return nil, fmt.Errorf(
				"failed to convert rotation config: %w", err)
		}
		logCore.Rotation = rotationConfig
	}

	// Convert encoder (optional)
	if encoder, ok := core["encoder"].(map[string]any); ok {
		logCore.Encoder = encoder
	}

	return logCore, nil
}

// convertMapToRotationConfig converts a map[string]any to a RotationConfig.
func convertMapToRotationConfig(rotation map[string]any) (*RotationConfig, error) {
	rotationConfig := &RotationConfig{}

	// Convert max size
	if maxSizeValue, exists := rotation["maxSize"]; exists {
		// Convert the value to int using ToInt utility function
		if maxSize, err := ToInt(maxSizeValue); err == nil {
			rotationConfig.MaxSize = maxSize
		} else {
			fmt.Printf("Error converting maxSize: %v\n", err)
		}
	} else {
		fmt.Println("maxSize key not found in rotation map.")
	}

	// Convert max backups
	if maxBackupsValue, exists := rotation["maxBackups"]; exists {
		// Convert the value to int using ToInt utility function
		if maxBackups, err := ToInt(maxBackupsValue); err == nil {
			rotationConfig.MaxBackups = maxBackups
		} else {
			fmt.Printf("Error converting maxBackups: %v\n", err)
		}
	} else {
		fmt.Println("maxBackups key not found in rotation map.")
	}

	// Convert max age
	if maxAgeValue, exists := rotation["maxAge"]; exists {
		// Convert the value to int using ToInt utility function
		if maxAge, err := ToInt(maxAgeValue); err == nil {
			rotationConfig.MaxAge = maxAge
		} else {
			fmt.Printf("Error converting maxAge: %v\n", err)
		}
	} else {
		fmt.Println("maxAge key not found in rotation map.")
	}

	// Convert compress
	if compress, ok := rotation["compress"].(bool); ok {
		rotationConfig.Compress = compress
	} else {
		return nil, errors.New(
			"invalid or missing 'compress' field in rotation config")
	}

	return rotationConfig, nil
}

// InitializeZapLogger initializes zap logger based on the logging configuration.
func InitializeZapLogger(config *LogConfig) (
	*zap.Logger, map[string]zap.AtomicLevel, error) {
	defaultEncoderConfig := zap.NewProductionEncoderConfig()
	defaultEncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	defaultEncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	var cores []zapcore.Core
	atomicLevels := make(map[string]zap.AtomicLevel)

	for coreName, coreConfig := range config.Cores {
		// Create a new atomic level for each core
		atomicLevel := zap.NewAtomicLevel()

		// Set the initial level for the atomic level
		var level zapcore.Level
		if err := level.UnmarshalText([]byte(coreConfig.Level)); err != nil {
			level = zapcore.InfoLevel
		}
		atomicLevel.SetLevel(level)
		atomicLevels[coreName] = atomicLevel

		// Determine the encoder type
		var encoder zapcore.Encoder
		if coreConfig.Encoding == "json" {
			encoder = zapcore.NewJSONEncoder(defaultEncoderConfig)
		} else {
			encoder = zapcore.NewConsoleEncoder(defaultEncoderConfig)
		}

		// Set up the core based on the type of output
		switch coreName {
		case "rotating_file":
			if coreConfig.Rotation == nil {
				return nil, nil, errors.New("missing rotation configuration for file core")
			}
			writer := zapcore.AddSync(&lumberjack.Logger{
				Filename:   coreConfig.OutputPath,
				MaxSize:    coreConfig.Rotation.MaxSize,
				MaxBackups: coreConfig.Rotation.MaxBackups,
				MaxAge:     coreConfig.Rotation.MaxAge,
				Compress:   coreConfig.Rotation.Compress,
			})
			core := zapcore.NewCore(encoder, writer, atomicLevel)
			cores = append(cores, core)

		case "console":
			core := zapcore.NewCore(encoder, zapcore.Lock(os.Stdout), atomicLevel)
			cores = append(cores, core)

		default:
			return nil, nil, fmt.Errorf("unknown core type: %v", coreName)
		}
	}

	logger := zap.New(zapcore.NewTee(cores...), zap.AddCaller())
	return logger, atomicLevels, nil
}

// Example function to log map[string]any as pretty JSON using zap.
func LogMapAsPrettyJSON(logger *zap.Logger, data map[string]any) {
	// Convert the map to pretty JSON
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		logger.Error("Failed to marshal map to pretty JSON", zap.Error(err))
		return
	}

	// Log the pretty JSON string
	logger.Info("Logging map as pretty JSON", zap.String("data", string(jsonData)))
}
