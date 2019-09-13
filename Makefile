# Based on https://sohlich.github.io/post/go_makefile/

MKDIR_P = mkdir -p

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOINSTALL=$(GOCMD) install
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

SRC_DIR=.
BIN_DIR=./bin
CMD_DIR=$(SRC_DIR)/cmd/tcc

all: test build

consumer: $(BIN_DIR)
	$(GOBUILD) -o $(BIN_DIR)/consumer -v $(CMD_DIR)

producer: $(BIN_DIR)
	$(GOBUILD) -o $(BIN_DIR)/producer -v $(CMD_DIR)

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

$(BIN_DIR):
	$(MKDIR_P) $(BIN_DIR)
