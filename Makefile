SHELL := /bin/bash

# --- Config (override on the command line or in env) ---
NAME  ?= Hetzner Terraform Provider
EMAIL ?= info@cipherbullet.com

# Isolated keyring (doesn't touch ~/.gnupg)
GNUPGHOME ?= $(CURDIR)/.gnupg-ci

# Output directory for exports
GPG_OUT ?= $(CURDIR)/.gpg

PUB_ASC  := $(GPG_OUT)/public.key
PRIV_ASC := $(GPG_OUT)/private.key
PUB_B64  := $(GPG_OUT)/public.key.b64
PRIV_B64 := $(GPG_OUT)/private.key.b64
FPR_FILE := $(GPG_OUT)/fingerprint.txt
KEYID_FILE := $(GPG_OUT)/keyid.txt

.ONESHELL:
.SHELLFLAGS := -eu -o pipefail -c

.PHONY: gpg-init gpg-print gpg-export gpg-export-b64 gpg-clean build


gpg-init:
	@mkdir -p "$(GNUPGHOME)" "$(GPG_OUT)"; chmod 700 "$(GNUPGHOME)"
	@{ \
	  echo "%no-protection"; \
	  echo "Key-Type: RSA"; \
	  echo "Key-Length: 4096"; \
	  echo "Subkey-Type: RSA"; \
	  echo "Subkey-Length: 4096"; \
	  echo "Name-Real: $(NAME)"; \
	  echo "Name-Email: $(EMAIL)"; \
	  echo "Expire-Date: 0"; \
	  echo "%commit"; \
	} > $(GPG_OUT)/gen-key.conf
	@gpg --batch --homedir "$(GNUPGHOME)" --generate-key $(GPG_OUT)/gen-key.conf
	@echo "✓ GPG key generated for: $(NAME) <$(EMAIL)>"
	$(MAKE) gpg-print gpg-export gpg-export-b64

gpg-print:
	@mkdir -p "$(GPG_OUT)"
	@FPR=$$(gpg --homedir "$(GNUPGHOME)" --list-keys --with-colons | awk -F: '/^fpr:/ {print $$10; exit}'); \
	LONGID=$$(gpg --homedir "$(GNUPGHOME)" --list-keys --keyid-format=long --with-colons | awk -F: '/^pub:/ {print $$5; exit}'); \
	echo "$$FPR" > $(FPR_FILE); \
	echo "$$LONGID" > $(KEYID_FILE); \
	echo "Long Key ID: $$LONGID"; \
	echo "Fingerprint: $$FPR" # secret: GPG_FINGERPRINT
	@gpg --homedir "$(GNUPGHOME)" --list-keys --keyid-format=long > $(GPG_OUT)/list.txt
	@gpg --homedir "$(GNUPGHOME)" --fingerprint > $(GPG_OUT)/fingerprint_details.txt || true
	@echo "✓ Info written to $(GPG_OUT)"

gpg-export: $(FPR_FILE)
	@mkdir -p "$(GPG_OUT)"
	@FPR=$$(cat $(FPR_FILE)); \
	gpg --homedir "$(GNUPGHOME)" --armor --export "$$FPR"            > "$(PUB_ASC)"; \
	gpg --homedir "$(GNUPGHOME)" --armor --export-secret-keys "$$FPR" > "$(PRIV_ASC)"; \
	echo "✓ Exported keys into $(GPG_OUT)"

gpg-export-b64: $(PUB_ASC) $(PRIV_ASC)
	@base64 < "$(PUB_ASC)"  > "$(PUB_B64)"
	@base64 < "$(PRIV_ASC)" > "$(PRIV_B64)" # secret: GPG_PRIVATE_KEY
	@echo "✓ Exported base64 versions into $(GPG_OUT)"

gpg-clean:
	@rm -rf "$(GNUPGHOME)" "$(GPG_OUT)"
	@echo "✓ Cleaned GPG dirs: $(GNUPGHOME), $(GPG_OUT)"

build:
	@GOOS=darwin GOARCH=arm64 go build -o ~/.terraform.d/plugins/github.com/cipherbullet/hetzner/0.1.0/darwin_arm64/terraform-provider-hetzner
	@echo "✓ Built provider binary"
