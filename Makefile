# Based on https://sohlich.github.io/post/go_makefile/

MKDIR_P = mkdir -p

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOINSTALL=$(GOCMD) install
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GODEP=dep

SRC_DIR=.
BIN_DIR=./bin
DATA_DIR=./data
VENDOR_DIR=./vendor 
CMD_DIR=$(SRC_DIR)/cmd/tcc

BINARY_NAME=tcc
BINARY_UNIX=$(BINARY_NAME)_unix

all: test build

build: $(BIN_DIR) $(VENDOR_DIR)
	$(GOBUILD) -o $(BIN_DIR)/$(BINARY_NAME) -v $(CMD_DIR)

install: $(VENDOR_DIR)
	$(GOINSTALL) -o $(BINARY_NAME) -v $(CMD_DIR)

test: $(VENDOR_DIR)
	$(GOTEST) -v $(SRC_DIR)/...

clean: 
	$(GOCLEAN)
	rm -rf $(BIN_DIR)
	rm -rf $(DATA_DIR)

run: build $(DATA_DIR)
	$(BIN_DIR)/$(BINARY_NAME)

$(VENDOR_DIR):
	$(GODEP) ensure

# run gen
$(DATA_DIR): build
	$(MKDIR_P) $(DATA_DIR)
	$(BIN_DIR)/$(BINARY_NAME) -d $(DATA_DIR) gen

$(BIN_DIR):
	$(MKDIR_P) $(BIN_DIR)
