FROM docker.io/golang:1.22.6@sha256:d5e49f92b9566b0ddfc59a0d9d85cd8a848e88c8dc40d97e29f306f07c3f8338 as builder
ARG TARGETOS
ARG TARGETARCH

WORKDIR /opt/app-root/src

COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download


# Copy the go source
COPY cmd/ cmd/
COPY api/ api/
COPY controllers/ controllers/
COPY internal/ internal/
COPY vendor/ vendor/

RUN CGO_ENABLED=1 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH} go build -a -o build/manager cmd/manager/main.go
RUN CGO_ENABLED=1 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH} go build -a -o build/server cmd/server/main.go

FROM registry.access.redhat.com/ubi9/ubi-minimal:latest

ARG DATA_DIR=/data
RUN mkdir $DATA_DIR && chmod 775 $DATA_DIR

WORKDIR /
#COPY --from=builder /opt/app-root/src/build/manager /usr/local/bin/
#COPY --from=builder /opt/app-root/src/build/server /usr/local/bin/
USER 65532:65532

ENTRYPOINT ["/usr/local/bin/manager"]
