
# Install/Update to the pinned CLI tool.
GF_CLI_VERSION ?= v2.10.0

.PHONY: cli
cli:
	@set -e; \
	wget -O gf https://github.com/gogf/gf/releases/download/$(GF_CLI_VERSION)/gf_$(shell go env GOOS)_$(shell go env GOARCH) && \
	chmod +x gf && \
	./gf install -y && \
	rm ./gf


# Check and install CLI tool.
.PHONY: cli.install
cli.install:
	@set -e; \
	gf -v > /dev/null 2>&1 || if [[ "$?" -ne "0" ]]; then \
  		echo "GoFame CLI is not installed, start proceeding auto installation..."; \
		make cli; \
	fi;