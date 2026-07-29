GO_VERSION ?= 1.23
RUST_VERSION ?= 1.80
HAL_MANIFEST := core/hal/Cargo.toml

.PHONY: all build build-core build-runtime build-hal fmt lint typecheck
.PHONY: test-unit test-integration test-hal security-scan clean qemu protoc-gen

all: build

build: build-core build-runtime build-hal

build-core:
	@echo "[build-core] Not yet implemented -- requires Yocto setup, kernel config, boot chain"

build-runtime:
	@echo "[build-runtime] Building Go runtime services..."
	@for dir in runtime/*/; do \
		if [ -f "$$dir/go.mod" ]; then \
			(cd "$$dir" && go build ./...) || exit 1; \
		fi; \
	done
	@echo "[build-runtime] Done"

build-hal:
	@echo "[build-hal] Building Rust HAL crate..."
	@if [ -f "$(HAL_MANIFEST)" ]; then \
		cargo build --manifest-path $(HAL_MANIFEST); \
	else \
		echo "  Skipping -- no HAL crate found"; \
	fi
	@echo "[build-hal] Done"

fmt:
	@echo "[fmt] Formatting Go code..."
	@if ls runtime/*/go.mod 2>/dev/null >/dev/null; then \
		for dir in runtime/*/; do \
			if [ -f "$$dir/go.mod" ]; then \
				(cd "$$dir" && gofmt -w .) || true; \
			fi; \
		done; \
	fi
	@echo "[fmt] Formatting Rust code..."
	@if [ -f "$(HAL_MANIFEST)" ]; then \
		cargo fmt --manifest-path $(HAL_MANIFEST); \
	else \
		echo "  Skipping -- no HAL crate found"; \
	fi
	@echo "[fmt] Done"

lint:
	@echo "[lint] Linting Go code..."
	@if ls runtime/*/go.mod 2>/dev/null >/dev/null; then \
		golangci-lint run runtime/...; \
	else \
		echo "  Skipping -- no Go modules found"; \
	fi
	@echo "[lint] Linting Rust code..."
	@if [ -f "$(HAL_MANIFEST)" ]; then \
		cargo clippy --manifest-path $(HAL_MANIFEST) -- -D warnings; \
	else \
		echo "  Skipping -- no HAL crate found"; \
	fi
	@echo "[lint] Done"

typecheck:
	@echo "[typecheck] Type-checking Go code..."
	@if ls runtime/*/go.mod 2>/dev/null >/dev/null; then \
		for dir in runtime/*/; do \
			if [ -f "$$dir/go.mod" ]; then \
				(cd "$$dir" && go vet ./...) || exit 1; \
			fi; \
		done; \
	fi
	@echo "[typecheck] Type-checking Rust code..."
	@if [ -f "$(HAL_MANIFEST)" ]; then \
		cargo check --manifest-path $(HAL_MANIFEST); \
	else \
		echo "  Skipping -- no HAL crate found"; \
	fi
	@echo "[typecheck] Done"

test-unit:
	@echo "[test-unit] Running Go unit tests..."
	@if ls runtime/*/go.mod 2>/dev/null >/dev/null; then \
		for dir in runtime/*/; do \
			if [ -f "$$dir/go.mod" ]; then \
				(cd "$$dir" && go test ./...) || exit 1; \
			fi; \
		done; \
	fi
	@echo "[test-unit] Running Rust unit tests..."
	@if [ -f "$(HAL_MANIFEST)" ]; then \
		cargo test --manifest-path $(HAL_MANIFEST); \
	else \
		echo "  Skipping -- no HAL crate found"; \
	fi
	@echo "[test-unit] Done"

test-integration:
	@echo "[test-integration] Running Go integration tests..."
	@if ls runtime/*/go.mod 2>/dev/null >/dev/null; then \
		for dir in runtime/*/; do \
			if [ -f "$$dir/go.mod" ]; then \
				(cd "$$dir" && go test -tags=integration ./...) || exit 1; \
			fi; \
		done; \
	fi
	@echo "[test-integration] Running Rust integration tests..."
	@if [ -f "$(HAL_MANIFEST)" ]; then \
		cargo test --manifest-path $(HAL_MANIFEST); \
	else \
		echo "  Skipping -- no HAL crate found"; \
	fi
	@echo "[test-integration] Done"

test-hal:
	@echo "[test-hal] Running HAL tests (requires hardware emulation)..."
	@if [ -f "$(HAL_MANIFEST)" ]; then \
		cargo test --manifest-path $(HAL_MANIFEST) --features hwtest; \
	else \
		echo "  Skipping -- no HAL crate found"; \
	fi
	@echo "[test-hal] Done"

security-scan:
	@echo "[security-scan] Scanning Go vulnerabilities..."
	@if ls runtime/*/go.mod 2>/dev/null >/dev/null; then \
		govulncheck ./runtime/...; \
	else \
		echo "  Skipping -- no Go modules found"; \
	fi
	@echo "[security-scan] Scanning Rust vulnerabilities..."
	@if [ -f "$(HAL_MANIFEST)" ]; then \
		cargo audit --manifest-path $(HAL_MANIFEST); \
	else \
		echo "  Skipping -- no HAL crate found"; \
	fi
	@echo "[security-scan] Done"

protoc-gen:
	@echo "[protoc-gen] Regenerating Protobuf Go code..."
	@for proto in contracts/protobuf/*/v1/*.proto; do \
		if [ -f "$$proto" ]; then \
			protoc --go_out=. --go_opt=paths=source_relative \
				--go-grpc_out=. --go-grpc_opt=paths=source_relative \
				-I contracts/protobuf \
				"$$proto"; \
		fi; \
	done
	@echo "[protoc-gen] Done"

qemu:
	@echo "[qemu] Not yet implemented -- requires bootable core OS image"

clean:
	@echo "[clean] Cleaning Go build artifacts..."
	@if ls runtime/*/go.mod 2>/dev/null >/dev/null; then \
		for dir in runtime/*/; do \
			if [ -f "$$dir/go.mod" ]; then \
				(cd "$$dir" && go clean -cache) || true; \
			fi; \
		done; \
	fi
	@echo "[clean] Cleaning Rust build artifacts..."
	@if [ -f "$(HAL_MANIFEST)" ]; then \
		cargo clean --manifest-path $(HAL_MANIFEST); \
	else \
		echo "  Skipping -- no HAL crate found"; \
	fi
	@echo "[clean] Done"
