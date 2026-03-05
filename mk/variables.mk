# Variables
DIST_DIR = dist
CMD_DIR = gopher
GO = go
GOFLAGS = -trimpath
LOCAL_BIN = .local/bin
LOCAL_LICENSES = .local/licenses

# Go environment
GOROOT = $(shell $(GO) env GOROOT)
GOPATH = $(shell $(GO) env GOPATH)
# Go 1.25+ uses lib/wasm, older versions use misc/wasm
WASM_EXEC_JS = $(shell if [ -f $(GOROOT)/lib/wasm/wasm_exec.js ]; then echo $(GOROOT)/lib/wasm/wasm_exec.js; else echo $(GOROOT)/misc/wasm/wasm_exec.js; fi)

# Source files
GO_SOURCES = $(shell find . -name '*.go' -not -path './.*' -not -path './vendor/*')

# File targets
DIST_WASM = $(DIST_DIR)/main.wasm
DIST_WASM_EXEC = $(DIST_DIR)/wasm_exec.js
SRC_INDEX = resources/index.html
DIST_INDEX = $(DIST_DIR)/index.html
GO_LICENSE = resources/credits/go/the-go-programming-language.md
LICENSES_MD = resources/credits/generated/licenses.md
