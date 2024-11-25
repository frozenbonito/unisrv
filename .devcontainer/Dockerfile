ARG VARIANT="1.21-bullseye"
FROM mcr.microsoft.com/devcontainers/go:1-${VARIANT}

SHELL ["/bin/bash", "-o", "pipefail", "-c"]

# Install Node.js
# hadolint ignore=DL3008
RUN curl -fsSL https://deb.nodesource.com/setup_lts.x | bash - \
  && export DEBIAN_FRONTEND=noninteractive \
  && apt-get install -y --no-install-recommends nodejs \
  && apt-get autoremove -y && apt-get clean -y && rm -rf /var/lib/apt/lists/*

# Install golangci-lint
ARG GOLANGCI_LINT_VERSION="1.62.2"
RUN curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh \
  | sh -s -- -b "$(go env GOPATH)/bin" "v${GOLANGCI_LINT_VERSION}"

# Install hadolint
ARG HADOLINT_VERSION="2.12.0"
RUN curl -SL "https://github.com/hadolint/hadolint/releases/download/v${HADOLINT_VERSION}/hadolint-Linux-$(uname -m)" \
  -o /usr/local/bin/hadolint \
  && chmod +x /usr/local/bin/hadolint

# Install commitlint
ARG COMMITLINT_CLI_VERSION="19.5.0"
ARG COMMITLINT_CONFIG_CONVENTIONAL_VERSION="19.5.0"
RUN npm install -g "@commitlint/cli@${COMMITLINT_CLI_VERSION}" "@commitlint/config-conventional@${COMMITLINT_CONFIG_CONVENTIONAL_VERSION}"
