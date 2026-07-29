# Core OS build system
# This file is included from the root Makefile

.PHONY: all kernel boot initramfs image

all: image

kernel:
	@echo "Building kernel..."
	@echo "  Config: core/kernel/config/fontis_defconfig"
	@echo "  Output: build/linux/arch/x86/boot/bzImage"
	@mkdir -p build/linux
	@echo "Run 'make kernel-linux' on a Linux build host."

boot:
	@echo "Building boot chain..."
	@echo "See core/boot/efi/ for UEFI Secure Boot configuration."

initramfs:
	@echo "Building initramfs..."
	@echo "See core/boot/initramfs/ for initramfs configuration."

image: kernel boot initramfs
	@echo "Assembling final OS image..."
	@echo "Output: build/core.img"
