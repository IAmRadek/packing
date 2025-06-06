BINARY_NAME=webd

# Build directory
BUILD_DIR=build

# Main application directory
MAIN_DIR=cmd/webd

# Environment file
ENV_FILE=.development.env

# Check if env file exists
ifneq ("$(wildcard $(ENV_FILE))","")
	include $(ENV_FILE)
	export
endif

.PHONY: all build clean run env-check

all: build

# Add env-check to ensure .development.env exists
env-check:
	@if [ ! -f "$(ENV_FILE)" ]; then \
		echo "Error: $(ENV_FILE) does not exist"; \
		exit 1; \
	fi

build:
	@echo "Building..."
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) ./$(MAIN_DIR)

clean:
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)

run: env-check
	@echo "Running with $(ENV_FILE)..."
	@go run ./$(MAIN_DIR)

# Helper target to show current env variables
env:
	@echo "Environment variables from $(ENV_FILE):"
	@cat $(ENV_FILE)
