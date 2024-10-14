FROM mcr.microsoft.com/vscode/devcontainers/go:1-1.23-bookworm

# Install GUI dependencies
RUN sudo apt-get update && sudo apt-get install -y gcc pkg-config libwayland-dev libx11-dev libx11-xcb-dev libxkbcommon-x11-dev libgles2-mesa-dev libegl1-mesa-dev libffi-dev libxcursor-dev libvulkan-dev

# install go tools
RUN go install gioui.org/cmd/gogio@latest
RUN go install github.com/abferm/hostdir@latest
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# protobuf compliler
ENV PROTOC_URL https://github.com/protocolbuffers/protobuf/releases/download/v3.11.4/protoc-3.11.4-linux-x86_64.zip
RUN wget ${WGET_OPTS} -O protoc.zip ${PROTOC_URL} && \
    unzip protoc.zip -d /usr/local/ && \
    rm protoc.zip

# Fix wierd permissions issue
RUN sudo chmod -R g+w /go/pkg