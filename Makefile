CURR_DIR:=$(shell pwd)
BIN_DIR=/usr/local/bin

# Go parameters
GOCMD = go
GOBUILD = $(GOCMD) build
GOCLEAN = $(GOCMD) clean

# Main package and executable name
PACKAGE = ./cmd/tomodoro-cli
EXECUTABLE = tomodoro
BUILD_PATH = ./build

# Run CLI target
run: clean build
	$(BUILD_PATH)/$(EXECUTABLE)

# Build CLI target
build: clean
	$(GOBUILD) -o $(BUILD_PATH)/$(EXECUTABLE) $(PACKAGE)

# Run CLI target with reflex to hot reload
watch:
	ulimit -n 1000
	reflex -s -r '\.go$$' make run

# Run API target
run-api: clean build-api
	$(BUILD_PATH)/$(EXECUTABLE_API)

# Clean target
clean:
	$(GOCLEAN)
	rm -rf $(BUILD_PATH)


# Install CLI target
install: build
	sudo cp $(BUILD_PATH)/$(EXECUTABLE) $(BIN_DIR)/tomodoro

# Uninstall CLI target
uninstall:
	sudo rm -rf $(BIN_DIR)/$(EXECUTABLE)

# Update CLI target
update:
	git pull origin main
	make uninstall-cli
	make install




