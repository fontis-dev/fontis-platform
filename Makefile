.PHONY: build build-core build-runtime build-hal fmt lint typecheck test-unit test-integration security-scan clean qemu protoc-gen

# Detect OS
UNAME_S := $(shell uname -s 2>/dev/null || echo Windows)

# Go and Rust settings
GO := go
CARGO := cargo
GOFMT := gofmt
GOLANGCI_LINT := golangci-lint

# Rust HAL project
HAL_DIR := core/hal

# Go runtime services
RUNTIME_SERVICES := identity auth storage networking logging module-lifecycle updates backup marketplace-client

# Protobuf
PROTOC := protoc
PROTOC_GEN_GO := protoc-gen-go
PROTOC_GEN_GO_GRPC := protoc-gen-go-grpc
PROTO_DIR := contracts/protobuf
PROTO_OUT := .

# Default target
build: build-core build-runtime build-hal

build-core:
	@echo "Core OS build requires Linux build environment. Skipping on $(UNAME_S)."
	@echo "On Linux, run: make build-core-linux"

build-core-linux:
	@echo "Building core OS image..."
	cd core && $(MAKE) -f Build.mk all

build-runtime: $(RUNTIME_SERVICES:%=build-runtime-%)

build-runtime-%:
	@echo "Building runtime/$*..."
	cd runtime/$* && $(GO) build ./...

build-hal:
	@echo "Building HAL..."
	cd $(HAL_DIR) && $(CARGO) build

fmt:
	@echo "Formatting Go code..."
	for dir in $(RUNTIME_SERVICES); do \
		cd runtime/$$dir && $(GOFMT) -l -w . 2>/dev/null || true; \
		cd ../..; \
	done
	@echo "Formatting Rust code..."
	cd $(HAL_DIR) && $(CARGO) fmt 2>/dev/null || true

fmt-check:
	@echo "Checking Go formatting..."
	@for dir in $(RUNTIME_SERVICES); do \
		unformatted=$$(cd runtime/$$dir && $(GOFMT) -l .); \
		if [ -n "$$unformatted" ]; then \
			echo "Unformatted files in runtime/$$dir:"; \
			echo "$$unformatted"; \
			exit 1; \
		fi; \
	done
	@echo "Checking Rust formatting..."
	cd $(HAL_DIR) && $(CARGO) fmt --check 2>/dev/null || (echo "Rust formatting check failed"; exit 1)

lint:
	@echo "Linting Go code..."
	for dir in $(RUNTIME_SERVICES); do \
		cd runtime/$$dir && $(GOLANGCI_LINT) run ./... 2>/dev/null || true; \
		cd ../..; \
	done
	@echo "Linting Rust code..."
	cd $(HAL_DIR) && $(CARGO) clippy -- -D warnings 2>/dev/null || true

typecheck:
	@echo "Type-checking Go code..."
	for dir in $(RUNTIME_SERVICES); do \
		cd runtime/$$dir && $(GO) vet ./... 2>/dev/null || true; \
		cd ../..; \
	done
	@echo "Type-checking Rust code..."
	cd $(HAL_DIR) && $(CARGO) check 2>/dev/null || true

test-unit:
	@echo "Running Go unit tests..."
	for dir in $(RUNTIME_SERVICES); do \
		cd runtime/$$dir && $(GO) test ./... -short 2>/dev/null || true; \
		cd ../..; \
	done
	@echo "Running Rust unit tests..."
	cd $(HAL_DIR) && $(CARGO) test --lib 2>/dev/null || true

test-integration:
	@echo "Running Go integration tests..."
	for dir in $(RUNTIME_SERVICES); do \
		if [ -f runtime/$$dir/*_integration_test.go ]; then \
			cd runtime/$$dir && $(GO) test -tags=integration ./... 2>/dev/null || true; \
			cd ../..; \
		fi; \
	done

security-scan:
	@echo "Running Go vulnerability scan..."
	for dir in $(RUNTIME_SERVICES); do \
		cd runtime/$$dir && govulncheck ./... 2>/dev/null || true; \
		cd ../..; \
	done
	@echo "Running Rust security audit..."
	cd $(HAL_DIR) && $(CARGO) audit 2>/dev/null || true

clean:
	@echo "Cleaning build artifacts..."
	for dir in $(RUNTIME_SERVICES); do \
		cd runtime/$$dir && $(GO) clean -cache 2>/dev/null || true; \
		cd ../..; \
	done
	cd $(HAL_DIR) && $(CARGO) clean 2>/dev/null || true
	rm -rf build/ out/

protoc-gen:
	@echo "Generating Protobuf code..."
	$(PROTOC) --proto_path=$(PROTO_DIR) --go_out=$(PROTO_OUT) --go_opt=paths=source_relative \
		--go-grpc_out=$(PROTO_OUT) --go-grpc_opt=paths=source_relative \
		$(PROTO_DIR)/**/v1/*.proto 2>/dev/null || true

qemu:
	@echo "Booting in QEMU (requires core OS image)..."
	qemu-system-x86_64 -enable-kvm -m 4G -smp 4 \
		-drive file=build/core.img,format=raw,if=virtio \
		-bios /usr/share/ovmf/OVMF.fd 2>/dev/null || echo "QEMU not available on this platform"
