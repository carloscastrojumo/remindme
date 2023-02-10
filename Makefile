FILENAME = rmm

ifeq ($(OS),Windows_NT)
	GO_BIN_DIR = $(GOPATH)\bin\\
	FILENAME := $(FILENAME).exe
else
	UNAME := $(shell uname -s)
	ifeq ($(UNAME),Darwin)
		GO_BIN_DIR = /Users/$(USER)/go/bin/
	else ifeq ($(UNAME),Linux)
		GO_BIN_DIR = $(HOME)/go/bin/
	else
    	$(error OS not supported by this Makefile)
	endif
endif


install:
	@echo "Installing Remindme as $(FILENAME) in $(GO_BIN_DIR)"
	go build -o $(GO_BIN_DIR)$(FILENAME) main.go
