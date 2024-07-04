MAKEFILEDIR := $(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))
PROJECT_ROOT := $(MAKEFILEDIR)
PROJECT_NAME ?= $(notdir $(PROJECT_ROOT))
WS_ROOT := $(realpath $(PROJECT_ROOT)/..)
PYDIR := $(PROJECT_ROOT)/py
PYPROJECT_TOML := $(PYDIR)/pyproject.toml
POETRY_LOCK := $(PYDIR)/poetry.lock
REQUIREMENTS := $(MAKEFILEDIR)/py/requirements.txt
LOCAL_REQUIREMENTS := $(MAKEFILEDIR)/py/local_requirements.txt
TEMP_COMBINED_REQUIREMENTS := $(MAKEFILEDIR)/py/temp_combined_requirements.txt
TEMP_REQ_CLEANUP ?= yes
PIP_NO_CACHE_DIR ?= no
ifeq ($(PIP_NO_CACHE_DIR),yes)
PIP_NO_CACHE_DIR_STR := --no-cache-dir
else
PIP_NO_CACHE_DIR_STR :=
endif

PY_ENV_MGR ?= poetry

PYSRCDIR := $(PYDIR)/src

ifdef PY_ENV_MGR
include $(MAKEFILEDIR)/py/$(PY_ENV_MGR).mk
endif

.PHONY: lint

pylint: py-install
	@$(ACTIVATE_VENV); \
  echo "running pylint at location: " `which pylint`; \
  $(VENV_BIN)/pylint --rcfile $(PYDIR)/.pylintrc $(PYSRCDIR)/apconf; \
	find $(PYSRCDIR)/scripts -type f -name "*.py" -exec $(VENV_BIN)/pylint --rcfile $(PYDIR)/.pylintrc {} + ; \
	find $(PYDIR)/test -type f -name "*.py" -exec $(VENV_BIN)/pylint --rcfile $(PYDIR)/.pylintrc {} +

test: py-install
	@$(ACTIVATE_VENV); \
  $(VENV_BIN)/python -m unittest discover -s $(PYDIR)/test

check: pylint test

all: check

clean: py-clean artifact-clean

artifact-clean:
	rm -rf $(PROJECT_ROOT)/artifacts
