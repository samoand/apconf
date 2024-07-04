# Application Configurator

This Python library simplifies the declaration and management of configurations using an easy-to-read, Kubernetes-like format. It supports dynamic changes and offers APIs for deploying configurations.

The library addresses the challenge of managing multiple complex configurations linked to various deployment scenarios while maintaining the same application logic.

Users can break down configurations into different profiles, store these profiles in respective directories, and specify the relevant combination of these directories during initialization.

## Features

- Readable declaration of configurations (in yaml).
- Support for templates in the declarations.
- Dynamic modification of configurations.
- Plugin APIs for configuration deployment.
- Self-explanatory usage through provided examples.

## Usage

### Sample Usage

Refer to the `py/src/scripts/sample-main.py` for an example of how to use this library

This sample focuses on logging configuration as an example.

### Example Test

Refer to `py/test/*.py`

### Sample Config Snippets

```yaml
---
kind: Config
metadata:
  name: crawler_config
spec:
  num_workers: {{ num_workers }}
---
kind: Config
metadata:
  name: logging_config
spec:
  version: 1
  disable_existing_loggers: false
  formatters:
    standard:
      format: "%(threadName)s - %(name)s - %(levelname)s - %(funcName)s - %(message)s - %(asctime)s"
  handlers:
    console_handler:
      level: ERROR
      class: logging.StreamHandler
      formatter: standard
    file_handler:
      level: WARNING
      class: logging.handlers.RotatingFileHandler
      formatter: standard
      # log_path_desc will be normalized
      # and the key converted to "filename"
      # for compatibility with dictConfig
      log_path_desc: >-
        {{ project_root }}
        artifacts
        test
        log
        myapp.{{ proc_id }}.log
      maxBytes: 65536
      backupCount: 5000
      encoding: utf8
  loggers:
    "root":
      handlers:
        - file_handler
        - console_handler
      level: WARNING
```
