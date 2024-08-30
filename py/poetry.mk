PY_PROJECT_TOML := $(PROJECT_PY_DIR)/pyproject.toml
PY_POETRY_LOCK := $(PROJECT_PY_DIR)/poetry.lock
PY_REQUIREMENTS := $(PROJECT_PY_DIR)/requirements.txt
PY_SRC_DIR := $(PROJECT_PY_DIR)/src

PY_LOCAL_REQUIREMENTS := $(MAKEFILEDIR)/py/local_requirements.txt
PY_TEMP_COMBINED_REQUIREMENTS := $(MAKEFILEDIR)/py/temp_combined_requirements.txt
TEMP_REQ_CLEANUP ?= yes
PIP_NO_CACHE_DIR ?= no
ifeq ($(PIP_NO_CACHE_DIR),yes)
PIP_NO_CACHE_DIR_STR := --no-cache-dir
else
PIP_NO_CACHE_DIR_STR :=
endif

PY_VENV_DIR := $(WS_ROOT)/external/py/venv
PY_VENV_BIN := $(PY_VENV_DIR)/bin

ifeq ($(ARCH),arm64)
SYSTEM_PYTHON := ~/.pyenv/versions/3.11.2/bin/python
else ifeq ($(ARCH),x86_64)
SYSTEM_PYTHON := /usr/local/bin/python3
else
$(error Unsupported architecture)
endif

.PHONY: poetry-install, poetry-local-install

ACTIVATE_VENV = if [ "$$VIRTUAL_ENV" != "$(PY_VENV_DIR)" ]; then \
	echo "Activating the correct virtual environment..."; \
		. $(PY_VENV_DIR)/bin/activate; \
	else \
		echo "The correct virtual environment $(PY_VENV_DIR) is already activated."; \
	fi

$(PY_VENV_DIR):
	echo "Creating virtual environment..."
	$(SYSTEM_PYTHON) -m venv $(PY_VENV_DIR)
	$(PY_VENV_DIR)/bin/pip install --upgrade pip
	$(PY_VENV_DIR)/bin/pip install $(PIP_NO_CACHE_DIR_STR) poetry

venv: $(PY_VENV_DIR)

PY_POETRY_LOCAL_INSTALL_TIMESTAMP := .poetry-local-install.timestamp
PY_POETRY_INSTALL_TIMESTAMP := .poetry-install.timestamp

poetry-install: $(PY_POETRY_INSTALL_TIMESTAMP)

$(PY_POETRY_INSTALL_TIMESTAMP): $(PY_PROJECT_TOML) $(PY_VENV_DIR)
	@cd $(PROJECT_PY_DIR) && . $(PY_VENV_DIR)/bin/activate && pip install --upgrade pip &&  poetry install
	touch $(PY_POETRY_INSTALL_TIMESTAMP)

poetry-local-install: $(PY_POETRY_LOCAL_INSTALL_TIMESTAMP)

$(PY_POETRY_LOCAL_INSTALL_TIMESTAMP): $(PY_PROJECT_TOML) $(PY_VENV_DIR)
	cd $(PROJECT_PY_DIR) && . $(PY_VENV_DIR)/bin/activate && pip install --upgrade pip &&  pip install -e $(PROJECT_PY_DIR)
	touch $(PY_POETRY_LOCAL_INSTALL_TIMESTAMP)

py-req-install: $(wildcard $(PY_REQUIREMENTS)) $(wildcard $(PY_LOCAL_REQUIREMENTS))
	$(MAKE) clean-temp-reqs
	@if [ -f $(PY_REQUIREMENTS) ]; then \
		cp $(PY_REQUIREMENTS) $(PY_TEMP_COMBINED_REQUIREMENTS); \
	fi
	@if  [ -f $(PY_LOCAL_REQUIREMENTS) ]; then \
		while IFS= read -r line; do \
			export PROJECT_PY_DIR="$(PROJECT_PY_DIR)"; \
			export MAKEFILEDIR="$(MAKEFILEDIR)"; \
			resolved_line=$$(echo $$line | sed 's@\$$$(PROJECT_PY_DIR)@'"$$PROJECT_PY_DIR"'@g' | sed 's@\$$$(MAKEFILEDIR)@'"$$MAKEFILEDIR"'@g' | envsubst | xargs realpath); \
			echo "-e $$resolved_line" >> $(PY_TEMP_COMBINED_REQUIREMENTS); \
		done < $(PY_LOCAL_REQUIREMENTS); \
	fi
	@if [ -f $(PY_TEMP_COMBINED_REQUIREMENTS) ]; then \
		$(PY_VENV_DIR)/bin/pip install --upgrade pip; \
		$(PY_VENV_DIR)/bin/pip install $(PIP_NO_CACHE_DIR_STR) -r $(PY_TEMP_COMBINED_REQUIREMENTS); \
	fi
	if [ "$(TEMP_REQ_CLEANUP)" = "yes" ]; then \
		$(MAKE) clean-temp-reqs; \
	fi

clean-temp-reqs:
	@if [ -f $(PY_TEMP_COMBINED_REQUIREMENTS) ]; then \
		rm -f $(PY_TEMP_COMBINED_REQUIREMENTS); \
	fi

py-install: poetry-local-install py-req-install

py-lint: py-install
	@$(ACTIVATE_VENV); \
  echo "running py-lint at location: " `which pylint`; \
  $(PY_VENV_BIN)/pylint --rcfile $(PROJECT_PY_DIR)/.pylintrc $(PY_SRC_DIR)/apconf; \
	find $(PY_SRC_DIR)/scripts -type f -name "*.py" -exec $(PY_VENV_BIN)/pylint --rcfile $(PROJECT_PY_DIR)/.pylintrc {} + ; \
	find $(PROJECT_PY_DIR)/test -type f -name "*.py" -exec $(PY_VENV_BIN)/pylint --rcfile $(PROJECT_PY_DIR)/.pylintrc {} +

py-test: py-install
	@$(ACTIVATE_VENV); \
  $(PY_VENV_BIN)/python -m unittest discover -s $(PROJECT_PY_DIR)/test

clean-poetry-lock:
	rm -rf $(PY_POETRY_LOCK)

clean-py-venv:
	rm -rf $(PY_VENV_DIR) $(PY_POETRY_LOCAL_INSTALL_TIMESTAMP) $(PY_POETRY_INSTALL_TIMESTAMP)

clean-poetry-timestamp-files:
	rm -rf $(PY_POETRY_LOCAL_INSTALL_TIMESTAMP) $(PY_POETRY_INSTALL_TIMESTAMP)

py-clean: clean-py-venv clean-poetry-lock clean-poetry-timestamp-files

py-check: py-lint py-test
