ARG VARIANT="1.18-bullseye"
FROM mcr.microsoft.com/devcontainers/go:0-${VARIANT}

SHELL ["/bin/bash", "-o", "pipefail", "-c"]

# Install Node.js
# hadolint ignore=DL3008
RUN curl -fsSL https://deb.nodesource.com/setup_lts.x | bash - \
  && export DEBIAN_FRONTEND=noninteractive \
  && apt-get install -y --no-install-recommends nodejs \
  && apt-get autoremove -y && apt-get clean -y && rm -rf /var/lib/apt/lists/*

# Install golangci-lint
ARG GOLANGCI_LINT_VERSION="1.49.0"
RUN curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh \
  | sh -s -- -b "$(go env GOPATH)/bin" "v${GOLANGCI_LINT_VERSION}"

# Install hadolint
ARG HADOLINT_VERSION="2.10.0"
RUN curl -SL "https://github.com/hadolint/hadolint/releases/download/v${HADOLINT_VERSION}/hadolint-Linux-$(uname -m)" \
  -o /usr/local/bin/hadolint \
  && chmod +x /usr/local/bin/hadolint

# Install commitlint
ARG COMMITLINT_CLI_VERSION="17.5.1"
ARG COMMITLINT_CONFIG_CONVENTIONAL_VERSION="17.4.4"
RUN npm install -g "@commitlint/cli@${COMMITLINT_CLI_VERSION}" "@commitlint/config-conventional@${COMMITLINT_CONFIG_CONVENTIONAL_VERSION}"

# Install go tools
RUN go install github.com/goreleaser/goreleaser@latest \
  && go install gotest.tools/gotestsum@latest

# Install act
RUN curl https://raw.githubusercontent.com/nektos/act/master/install.sh | bash
