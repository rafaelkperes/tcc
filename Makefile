# Based on https://sohlich.github.io/post/go_makefile/

MKDIR_P = mkdir -p

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOINSTALL=$(GOCMD) install
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

SRC_DIR=./src
BIN_DIR=./bin
DATA_DIR=./data
CMD_DIR=$(SRC_DIR)/cmd

BINARY_NAME=tcc
BINARY_UNIX=$(BINARY_NAME)_unix

all: test build

build: $(BIN_DIR)
	$(GOBUILD) -o $(BIN_DIR)/$(BINARY_NAME) -v $(CMD_DIR)

install:
	$(GOINSTALL) -o $(BINARY_NAME) -v $(CMD_DIR)

test: 
	$(GOTEST) -v $(SRC_DIR)/...

clean: 
	$(GOCLEAN)
	rm -rf $(BIN_DIR)
	rm -rf $(DATA_DIR)

run: build $(DATA_DIR)
	$(BIN_DIR)/$(BINARY_NAME)

# run gen
$(DATA_DIR): build
	$(MKDIR_P) $(DATA_DIR)
	$(BIN_DIR)/$(BINARY_NAME) -d $(DATA_DIR) gen

$(BIN_DIR):
	$(MKDIR_P) $(BIN_DIR)
