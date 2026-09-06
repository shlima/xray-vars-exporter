PLATFORMS := linux/amd64 linux/arm64 darwin/amd64 darwin/arm64
EXE_EXT := .exe
BUILD_DIR := build
APP_NAME := xray-vars-exporter

.PHONY: build-all
build-all:
	@echo "Building for multiple platforms..."
	@rm -rf $(BUILD_DIR)
	@mkdir -p $(BUILD_DIR)
	@for platform in $(PLATFORMS); do \
		os=$$(echo $$platform | cut -d'/' -f1); \
		arch=$$(echo $$platform | cut -d'/' -f2); \
		output=$(BUILD_DIR)/$(APP_NAME)-$$os-$$arch; \
		if [ "$$os" = "windows" ]; then \
			output="$$output$(EXE_EXT)"; \
		fi; \
		echo "Building $$output..."; \
		GOOS=$$os GOARCH=$$arch CGO_ENABLED=0 \
		go build -ldflags="-s -w" -trimpath \
		-o $$output .; \
	done
	@echo "All builds complete!"
	@ls -lh $(BUILD_DIR)
