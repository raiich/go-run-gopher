# License-related internal tasks

# File target: go-licenses in .local/bin
$(LOCAL_BIN)/go-licenses: $(LOCAL_BIN)
	@echo "Installing go-licenses to $(LOCAL_BIN)..."
	@GOBIN=$(PWD)/$(LOCAL_BIN) $(GO) install github.com/google/go-licenses@latest
	@echo "✓ Installed $@"

# File target: .local/licenses directory
$(LOCAL_LICENSES):
	@mkdir -p $@
	@echo "✓ Created $@"

# File target: Go programming language license
$(GO_LICENSE):
	echo '## The Go Programming Language' > $(GO_LICENSE)
	echo '' >> $(GO_LICENSE)
	echo '```' >> $(GO_LICENSE)
	cat $(dirname $(which go))/../LICENSE >> $(GO_LICENSE)
	echo '```' >> $(GO_LICENSE)
	echo '' >> $(GO_LICENSE)

# File target: Generated licenses markdown
$(LICENSES_MD): $(LOCAL_BIN)/go-licenses $(GO_SOURCES) go.mod go.sum
	@echo "Cleaning $(LOCAL_LICENSES)..."
	@rm -rf $(LOCAL_LICENSES)
	@echo "Saving license files to $(LOCAL_LICENSES)..."
	@$(LOCAL_BIN)/go-licenses save ./$(CMD_DIR) --save_path $(LOCAL_LICENSES)
	@echo "Generating licenses.md..."
	@mkdir -p resources/credits/generated
	@$(LOCAL_BIN)/go-licenses report ./$(CMD_DIR) 2>/dev/null | \
		xargs -I {} ./scripts/generate-license-md.sh {} $(LOCAL_LICENSES) \
		> $@ 2>&1 || (cat $@ && false)
	@echo "✓ Generated $@"

.PHONY: licenses-save
licenses-save: $(LOCAL_BIN)/go-licenses $(LOCAL_LICENSES) ## Save license files to .local/licenses
	@echo "Saving license files to $(LOCAL_LICENSES)..."
	@$(LOCAL_BIN)/go-licenses save ./$(CMD_DIR) --save_path $(LOCAL_LICENSES)
	@echo "✓ Saved license files to $(LOCAL_LICENSES)"

.PHONY: licenses-check
licenses-check: $(LOCAL_BIN)/go-licenses ## Check third-party licenses
	@echo "Checking licenses..."
	@$(LOCAL_BIN)/go-licenses report ./$(CMD_DIR)
