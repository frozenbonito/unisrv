ARG VARIANT="1.18-bullseye"
FROM mcr.microsoft.com/devcontainers/go:0-${VARIANT}

SHELL ["/bin/bash", "-o", "pipefail", "-c"]

# Install golangci-lint
ARG GOLANGCI_LINT_VERSION="1.49.0"
RUN curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh \
  | sh -s -- -b "$(go env GOPATH)/bin" "v${GOLANGCI_LINT_VERSION}"

# Install hadolint
ARG HADOLINT_VERSION="2.10.0"
RUN curl -SL "https://github.com/hadolint/hadolint/releases/download/v${HADOLINT_VERSION}/hadolint-Linux-$(uname -m)" \
  -o /usr/local/bin/hadolint \
  && chmod +x /usr/local/bin/hadolint

# Install go tools
RUN go install github.com/goreleaser/goreleaser@latest \
  && go install gotest.tools/gotestsum@latest

# Install act
RUN curl https://raw.githubusercontent.com/nektos/act/master/install.sh | bash
