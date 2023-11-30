# Build the manager binary
FROM golang:1.19 as builder
ARG TARGETOS
ARG TARGETARCH

WORKDIR /workspace
# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum

RUN go mod download

COPY main.go main.go

COPY config/ config/
COPY database/ database/
COPY ginserver/ ginserver/
COPY monitoring/ monitoring/
COPY note/ note/

COPY view/ view/

RUN CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH} go build -a -o notepad main.go

FROM gcr.io/distroless/static:nonroot
WORKDIR /
COPY --from=builder /workspace/notepad .
COPY assets/ assets/
COPY templates/ templates/
USER 65532:65532

ENTRYPOINT ["/notepad"]