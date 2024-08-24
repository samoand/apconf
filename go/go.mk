GO_VERSION=1.23.0
GO_INSTALL_DIR=$(WS_ROOT)/external/go$(GO_VERSION)
GOPATH=$(GO_INSTALL_DIR)/gopath
GOROOT=$(GO_INSTALL_DIR)/go
GO_BIN=$(GOROOT)/bin
GO_PATH_BIN=$(GO_INSTALL_DIR)/gopath/bin
GO_MOD_MODE=readonly
LINTER=$(GO_PATH_BIN)/golangci-lint

$(GO_INSTALL_DIR):
	@mkdir -p $@

GO_FILES := $(wildcard $(PROJECT_GO_DIR)/**/*.go)
GO_MOD := $(PROJECT_GO_DIR)/go.mod

$(GOROOT): $(GO_INSTALL_DIR)
	$(PROJECT_GO_DIR)/go_install.sh $(GO_VERSION) $(GO_INSTALL_DIR) $(OS) $(ARCH)

go-install: $(GOROOT)

$(GO_PATH_BIN):
	@mkdir -p $@

$(GOPATH): $(GOROOT) $(GO_FILES) $(GO_MOD)
	@mkdir -p $@
	cd $(PROJECT_GO_DIR) && $(GO_BIN)/go mod tidy

go-dep-install: $(GOPATH)

$(LINTER): $(GO_PATH_BIN)
	@if [ ! -f "$(LINTER)" ]; then \
		echo "Installing golangci-lint..." ; \
		$(GO_BIN)/go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest ; \
		chmod +x $(GO_PATH_BIN)/golangci-lint ; \
	fi

go-lint: | $(LINTER)
	cd $(PROJECT_GO_DIR) && $(LINTER) run -c .golangci.yml --modules-download-mode=$(GO_MOD_MODE) ./...

go-test:
	cd $(PROJECT_GO_DIR) && $(GO_BIN)/go test -v ./...

go-check: go-install go-dep-install go-lint go-test

clean-goroot:
	rm -rf $(GOROOT)

clean-gopath:
	rm -rf $(GOPATH)

go-clean:
	chmod -R 777 $(GO_INSTALL_DIR)
	rm -rf $(GO_INSTALL_DIR)

.PHONY: go-clean clean-gopath clean-goroot go-check go-test go-lint go-install
