package apconf

import (
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	configRoot          string
	configBasenames     []string
	templateParams      map[string]any
	configPreprocessors []func(map[string]any)
	configDeployers     []func(map[string]any, map[string]any, ConfigDiffResult) error
	configValidators    []func(map[string]any, map[string]any, ConfigDiffResult) bool
	config              map[string]any
}

type ApconfException struct {
	message string
}

func (e *ApconfException) Error() string {
	return e.message
}

func NewConfig(
	configRoot string,
	configBasenames []string,
	templateParams map[string]any,
	configPreprocessors []func(map[string]any),
	configDeployers []func(map[string]any, map[string]any, ConfigDiffResult) error,
	configValidators []func(map[string]any, map[string]any, ConfigDiffResult) bool) *Config {

	c := &Config{
		configRoot:          configRoot,
		configBasenames:     configBasenames,
		templateParams:      templateParams,
		configPreprocessors: configPreprocessors,
		configDeployers:     configDeployers,
		configValidators:    configValidators,
		config:              make(map[string]any),
	}
	c.init()
	return c
}

func (c *Config) init() {
	configDirs := make([]string, len(c.configBasenames))
	for i, basename := range c.configBasenames {
		configDirs[i] = filepath.Join(c.configRoot, basename)
	}

	for _, dir := range configDirs {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			panic(&ApconfException{message: fmt.Sprintf("Config directory %s does not exist.", dir)})
		}
	}

	config := processYamlDirs(configDirs, c.templateParams)["Config"].(map[string]any)
	c.apply(config)
}

func (c *Config) preprocess(config map[string]any) {
	if c.configPreprocessors == nil {
		return
	}
	for _, preprocessor := range c.configPreprocessors {
		preprocessor(config)
	}
}

func (c *Config) validate(configNew map[string]any, configDiffResult ConfigDiffResult) bool {
	if c.configValidators == nil {
		return true
	}
	for _, validator := range c.configValidators {
		if !validator(configNew, c.config, configDiffResult) {
			return false
		}
	}
	return true
}

func (c *Config) deploy(configNew map[string]any, configDiffResult ConfigDiffResult) []error {
	if c.configDeployers == nil {
		return nil
	}
	var errors []error
	for _, deployer := range c.configDeployers {
		if err := deployer(configNew, c.config, configDiffResult); err != nil {
			errors = append(errors, err)
		}
	}
	return errors
}

func (c *Config) apply(config map[string]any) []error {
	configDiffResult := ConfigDiff(config, c.config)
	c.preprocess(config)
	valid := c.validate(config, configDiffResult)
	if !valid {
		return []error{fmt.Errorf("config %s failed validation", config)}
	}
	errs := c.deploy(config, configDiffResult)
	if errs != nil {
		return errs
	}
	c.config = config
	return nil
}
