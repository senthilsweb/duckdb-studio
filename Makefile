# Variables
WEB_DIR=.
DIST_DIR=.dist
RELEASE_DIR=.release
OS=$(shell uname)
binary_name=duckdbstudio

# Ensure variables are exported to all uses in the Makefile
export

# Phony targets ensure that Make doesn't look for files named like the target
.PHONY: clean generate copy build package ensure_releases_dir

all: clean generate copy build package

clean:
	@echo "Cleaning previous web build artifacts..."
	@find $(DIST_DIR) -mindepth 1 -delete || mkdir -p $(DIST_DIR)
	@echo "Cleaned..."

ensure_releases_dir:
	@echo "Ensuring $(RELEASE_DIR) exists..."
	@mkdir -p $(RELEASE_DIR)

generate:
	@echo "Installing web dependencies..."
	@cd $(WEB_DIR) && npm i
	@echo "Dependencies installed..."
	@echo "Generating build..."
	@cd $(WEB_DIR) && npm run generate || (echo "Error generating build. Exiting..." && exit 1)
	@echo "Generated..."

copy:
	@echo "Copying to $(DIST_DIR) for embedding..."
	@mkdir -p $(DIST_DIR)
	@cp -r $(WEB_DIR)/.output/public/* $(DIST_DIR)/
	@echo "Copied..."

build: ensure_releases_dir
	@echo "Building Go binary for $(OS)..."
ifeq ($(OS),Linux)
	@go build -o $(RELEASE_DIR)/$(binary_name)-linux-aarch64 -v .
else ifeq ($(OS),Darwin)
	#@go build -o $(RELEASE_DIR)/$(binary_name)-macos -v .
	@go build -o $(RELEASE_DIR)/$(binary_name)-linux-aarch64 -v .
else
	@echo "Unsupported OS: $(OS). Exiting..."
	@exit 1
endif
	@echo "Binary built in $(RELEASE_DIR)..."

package: ensure_releases_dir
	@echo "Packaging files..."
	@ZIP_NAME="$(RELEASE_DIR)/$(binary_name)_bundle_`date +%s`.zip"; \
	zip $$ZIP_NAME $(RELEASE_DIR)/$(binary_name)-* config.yml
	@echo "Packaged into $$ZIP_NAME..."
