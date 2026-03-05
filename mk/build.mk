# Build-related file targets

# File target: WebAssembly binary
$(DIST_WASM): $(GO_SOURCES)
	@mkdir -p $(DIST_DIR)
	@echo "Building WebAssembly..."
	GOOS=js GOARCH=wasm $(GO) build $(GOFLAGS) -ldflags="-s -w" -o $@ ./$(CMD_DIR)
	@echo "✓ Built $@"

# File target: wasm_exec.js
$(DIST_WASM_EXEC): $(WASM_EXEC_JS)
	@mkdir -p $(DIST_DIR)
	@echo "Copying wasm_exec.js..."
	@cp $< $@
	@echo "✓ Copied $@"

# File target: index.html
$(DIST_INDEX): $(SRC_INDEX)
	@mkdir -p $(DIST_DIR)
	@echo "Copying index.html..."
	@cp $< $@
	@echo "✓ Copied $@"
