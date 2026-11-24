FROM golang:1.27.4-bookworm@sha256:4fa7c968ba93f1635670522dbb01a2c78150443e1385746924566079f36570ea

RUN go install gotest.tools/gotestsum@latest

ENV PACT_LOG_LEVEL=debug
RUN go install github.com/pact-foundation/pact-go/v2@v2.7.1
RUN pact-go -l DEBUG install
