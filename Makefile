ifndef $(GOPATH)
	export GOPATH=$(shell go env GOPATH)
endif
ifndef $(GOBIN)
	export GOBIN=$(shell go env GOBIN)
endif

build:
	go build -i upip.go

deps:
	go get

install:
ifeq ($(OS),Windows_NT)
	mv upip.exe "$(USERPROFILE)\\"
else  # linux and other
	mv upip ~/.local/bin/
endif

all: deps build install
