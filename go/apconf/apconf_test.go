package apconf

import (
	"bufio"
	"fmt"
	"go.uber.org/zap"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"sync/atomic"
	"testing"
)

func pathFromDesc(pathDesc string) string {
	re := regexp.MustCompile(`\s+`)
	return filepath.Clean(filepath.Join(re.Split(pathDesc, -1)...))
}

func containsKey(m map[string]any, key string) bool {
	_, exists := m[key]
	return exists
}

func applyLogBuilder(atomicValue *atomic.Value) func(
	newConfig map[string]any,
	oldConfig map[string]any,
	configDiffResult ConfigDiffResult) error {

	return func(
		newConfig map[string]any,
		_ map[string]any,
		configDiffResult ConfigDiffResult) error {

		// Check if logging config has changed
		if !configDiffResult.Contains([]string{"zap_logging_config"}) {
			return nil
		}

		// Extract new logging configuration
		logConfigMap := newConfig["zap_logging_config"].(map[string]any)["spec"].(map[string]any)

		// Ensure all necessary directories for output paths are created
		cores := logConfigMap["cores"].(map[string]any)
		for _, core := range cores {
			if coreMap, ok := core.(map[string]any); ok {
				if outputPath, ok := coreMap["outputPath"].(string); ok {
					if err := os.MkdirAll(filepath.Dir(outputPath), os.ModePerm); err != nil {
						return fmt.Errorf("failed to create output path directories: %w", err)
					}
				}
			}
		}

		// Convert the map to LogConfig structure
		logConfig, err := ConvertMapToLogConfig(logConfigMap)
		if err != nil {
			return fmt.Errorf("failed to convert configuration: %w", err)
		}

		// Initialize atomicValue using Zap based on logConfig
		newLogger, _, err := InitializeZapLogger(logConfig)
		if err != nil {
			return fmt.Errorf("failed to initialize zap atomicValue: %w", err)
		}

		// Atomically store the new atomicValue
		atomicValue.Store(newLogger)

		return nil
	}
}

func TestConfig(t *testing.T) {
	// Setup
	configPreprocessors := []func(map[string]any){
		Preprocessor(
			func(key string) bool { return key == "outputPathDesc" },
			func(oldKey string) string { return "outputPath" },
			func(oldValue any) any { return pathFromDesc(oldValue.(string)) },
			false,
		),
	}
	var atomicValue atomic.Value
	configDeployers := []func(map[string]any, map[string]any, ConfigDiffResult) error{
		applyLogBuilder(&atomicValue),
	}

	// Use runtime.Caller to get the current file's directory
	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatalf("Failed to get current file")
	}

	projectRoot, err := findGitRoot(filepath.Dir(currentFile))
	if err != nil {
		t.Fatalf("Failed to find project root: %v", err)
	}
	configRoot := filepath.Join(projectRoot, "test_config")
	procID := os.Getpid()
	cfg := NewConfig(
		configRoot,
		[]string{"crawl", "logging-zap-go"},
		map[string]any{
			"ProjectRoot": projectRoot,
			"ProcId":      procID,
		},
		configPreprocessors,
		configDeployers,
		nil,
	)

	configuredFilename := cfg.config["zap_logging_config"].(map[string]any)["spec"].(map[string]any)["cores"].(map[string]any)["rotating_file"].(map[string]any)["outputPath"].(string)

	t.Run("test_init_config", func(t *testing.T) {
		// See that "crawler" config is there.
		if numWorkers, ok := cfg.config["crawler_config"].(map[string]any)["spec"].(map[string]any)["num_workers"]; !ok || numWorkers.(int) != 30 {
			t.Errorf("Expected 'num_workers' to be 30, but got %v", numWorkers)
		}

		// Check that "logging" config is properly set and enforced
		logFileEnd := filepath.Join("artifacts", "test", "log", "myapp."+strconv.Itoa(procID)+".log") // Correct conversion using strconv.Itoa
		if !strings.HasSuffix(configuredFilename, logFileEnd) {
			t.Errorf("Expected log filename to end with %s, but got %s", logFileEnd, configuredFilename)
		}

		logger := atomicValue.Load().(*zap.Logger)
		warningMsg := "This is a warning"
		logger.Warn(warningMsg)
		infoMsg := "This is a first info message, it should not be logged"
		logger.Info(infoMsg)

		// Verify that the warning message is in the log file
		if !checkLogContains(configuredFilename, warningMsg) {
			t.Errorf("Expected warning message to be in log file, but it was not found")
		}

		// Verify that the info message is NOT in the log file
		if checkLogContains(configuredFilename, infoMsg) {
			t.Errorf("Expected info message not to be in log file, but it was found")
		}
	})

	t.Run("test_modified_log_config", func(t *testing.T) {
		newConfig := deepClone(cfg.config)
		newConfig["zap_logging_config"].(map[string]any)["spec"].(map[string]any)["cores"].(map[string]any)["rotating_file"].(map[string]any)["level"] = "debug"
		newConfig["zap_logging_config"].(map[string]any)["spec"].(map[string]any)["cores"].(map[string]any)["console"].(map[string]any)["level"] = "debug"

		errs := cfg.apply(newConfig)
		if errs != nil {
			t.Fatalf("Failed to apply new config: %v", errs)
		}

		infoMsg := "This is a second info message, it should be logged"
		logger := atomicValue.Load().(*zap.Logger)
		logger.Info(infoMsg)

		// Verify that the second info message is in the log file
		if !checkLogContains(configuredFilename, infoMsg) {
			t.Errorf("Expected second info message to be in log file, but it was not found")
		}
	})
}

// checkLogContains checks if a given log message exists in the log file
func checkLogContains(logFilePath, logMsg string) bool {
	file, err := os.Open(logFilePath)
	if err != nil {
		fmt.Printf("Failed to open log file: %v\n", err)
		return false
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), logMsg) {
			return true
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading log file: %v\n", err)
	}
	return false
}
