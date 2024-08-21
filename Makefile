MAKEFILEDIR := $(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))
PROJECT_ROOT := $(MAKEFILEDIR)
PROJECT_NAME ?= $(notdir $(PROJECT_ROOT))
WS_ROOT := $(realpath $(PROJECT_ROOT)/..)
PROJECT_PY_DIR := $(PROJECT_ROOT)/py
PROJECT_GO_DIR := $(PROJECT_ROOT)/go
OS := $(shell uname)
ARCH := $(shell uname -m)

PY_ENV_MGR ?= poetry

ifdef PY_ENV_MGR
include $(PROJECT_PY_DIR)/$(PY_ENV_MGR).mk
endif
include $(PROJECT_GO_DIR)/go.mk

lint: py-lint go-lint
test: py-test go-test
check: py-check go-check

all: check

artifact-clean:
	rm -rf $(PROJECT_ROOT)/artifacts

clean: py-clean go-clean artifact-clean

.PHONY: lint test check
