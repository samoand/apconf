VENV_DIR := $(WS_ROOT)/external/py/venv
VENV_BIN := $(VENV_DIR)/bin

ARCH := $(shell uname -m)
ifeq ($(ARCH),arm64)
SYSTEM_PYTHON := ~/.pyenv/versions/3.11.2/bin/python
else ifeq ($(ARCH),x86_64)
SYSTEM_PYTHON := /usr/local/bin/python3
else
$(error Unsupported architecture)
endif

.PHONY: poetry-install, poetry-local-install

ACTIVATE_VENV = if [ "$$VIRTUAL_ENV" != "$(VENV_DIR)" ]; then \
	echo "Activating the correct virtual environment..."; \
		. $(VENV_DIR)/bin/activate; \
	else \
		echo "The correct virtual environment $(VENV_DIR) is already activated."; \
	fi

$(VENV_DIR):
	echo "Creating virtual environment..."
	$(SYSTEM_PYTHON) -m venv $(VENV_DIR)
	$(VENV_DIR)/bin/pip install --upgrade pip
	$(VENV_DIR)/bin/pip install $(PIP_NO_CACHE_DIR_STR) poetry

venv: $(VENV_DIR)

POETRY_LOCAL_INSTALL_TIMESTAMP := .poetry-local-install.timestamp
POETRY_INSTALL_TIMESTAMP := .poetry-install.timestamp

poetry-install: $(POETRY_INSTALL_TIMESTAMP)

$(POETRY_INSTALL_TIMESTAMP): $(PYPROJECT_TOML) $(VENV_DIR)
	@cd $(PYDIR) && . $(VENV_DIR)/bin/activate && pip install --upgrade pip &&  poetry install
	touch $(POETRY_INSTALL_TIMESTAMP)

poetry-local-install: $(POETRY_LOCAL_INSTALL_TIMESTAMP)

$(POETRY_LOCAL_INSTALL_TIMESTAMP): $(PYPROJECT_TOML) $(VENV_DIR)
	cd $(PYDIR) && . $(VENV_DIR)/bin/activate && pip install --upgrade pip &&  pip install -e $(PYDIR)
	touch $(POETRY_LOCAL_INSTALL_TIMESTAMP)

req-install: $(wildcard $(REQUIREMENTS)) $(wildcard $(LOCAL_REQUIREMENTS))
	$(MAKE) clean-temp-reqs
	@if [ -f ${REQUIREMENTS} ]; then \
		cp ${REQUIREMENTS} ${TEMP_COMBINED_REQUIREMENTS}; \
	fi
	@if  [ -f ${LOCAL_REQUIREMENTS} ]; then \
		while IFS= read -r line; do \
			export PYDIR="$(PYDIR)"; \
			export MAKEFILEDIR="$(MAKEFILEDIR)"; \
			resolved_line=$$(echo $$line | sed 's@\$$${PYDIR}@'"$$PYDIR"'@g' | sed 's@\$$${MAKEFILEDIR}@'"$$MAKEFILEDIR"'@g' | envsubst | xargs realpath); \
			echo "-e $$resolved_line" >> ${TEMP_COMBINED_REQUIREMENTS}; \
		done < ${LOCAL_REQUIREMENTS}; \
	fi
	@if [ -f $(TEMP_COMBINED_REQUIREMENTS) ]; then \
		$(VENV_DIR)/bin/pip install --upgrade pip; \
		$(VENV_DIR)/bin/pip install $(PIP_NO_CACHE_DIR_STR) -r $(TEMP_COMBINED_REQUIREMENTS); \
	fi
	if [ "$(TEMP_REQ_CLEANUP)" = "yes" ]; then \
		$(MAKE) clean-temp-reqs; \
	fi

clean-temp-reqs:
	@if [ -f $(TEMP_COMBINED_REQUIREMENTS) ]; then \
		rm -f ${TEMP_COMBINED_REQUIREMENTS}; \
	fi



py-install: poetry-local-install req-install

clean-poetry-lock:
	rm -rf $(POETRY_LOCK)

clean-venv:
	rm -rf $(VENV_DIR)

clean-timestamp-files:
	rm -rf $(POETRY_LOCAL_INSTALL_TIMESTAMP) $(POETRY_INSTALL_TIMESTAMP)

py-clean: clean-venv clean-poetry-lock clean-timestamp-files

