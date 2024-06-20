
DEFAULT: build-cur

ifeq ($(GOPATH),)
  GOPATH = $(HOME)/go
endif

build-cur:
	GOPATH=$(GOPATH) go install github.com/pefish/go-build-tool/cmd/...@latest
	$(GOPATH)/bin/go-build-tool

install: build-cur
	sudo install -C ./build/bin/linux/fuck-server /usr/local/bin/fuck-server

install-service: install
	sudo mkdir -p /etc/systemd/system
	sudo install -C -m 0644 ./script/fuck-server.service /etc/systemd/system/fuck-server.service
	sudo systemctl daemon-reload
	@echo
	@echo "fuck-server service installed."

