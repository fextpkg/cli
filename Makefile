ifndef $(GOPATH)
	export GOPATH=$(shell go env GOPATH)
endif
ifndef $(GOBIN)
	export GOBIN=$(shell go env GOBIN)
endif

build:
	go build -i fext.go

deps:
	go get

install:
ifeq ($(OS),Windows_NT)
	mv fext.exe "$(USERPROFILE)\\"
else
ifeq ($(shell uname),Linux)
	mv fext ~/.local/bin/
endif
ifeq ($(shell uname),Darwin)
	mv fext usr/local/bin/
endif
endif

all: deps build install
