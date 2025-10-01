.PHONY: build

build:
	GOOS=darwin GOARCH=arm64 go build -o ~/.terraform.d/plugins/github.com/cipherbullet/hetzner/0.1.0/darwin_arm64/terraform-provider-hetzner