FROM golang:1.23.4-bullseye

# Set working directory
WORKDIR /app

# Install only essential system tools
RUN apt-get update && apt-get install -y \
    git \
    redis-tools \
    openssh-client \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

# Create a non-root user
RUN groupadd --gid 1000 vscode \
    && useradd --uid 1000 --gid 1000 -m vscode \
    && chown -R vscode:vscode /app

# Switch to the vscode user for the remaining operations
USER vscode

# Set up SSH and Git configuration
RUN mkdir -p /home/vscode/.ssh \
    && ssh-keyscan gitlab.com >> /home/vscode/.ssh/known_hosts \
    && chmod 700 /home/vscode/.ssh \
    && chmod 600 /home/vscode/.ssh/known_hosts \
    && git config --global url."ssh://git@gitlab.com/".insteadOf "https://gitlab.com/"

# Install Go tools
RUN go install golang.org/x/tools/gopls@latest \
    && go install github.com/go-delve/delve/cmd/dlv@latest \
    && go install github.com/air-verse/air@latest

# Set environment variables
ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    GO111MODULE=on

# Default command when container starts
CMD ["bash"]