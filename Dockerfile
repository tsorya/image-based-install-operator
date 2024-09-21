FROM registry.ci.openshift.org/ocp/builder:rhel-9-golang-1.22-openshift-4.17 as builder
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

#FROM registry.access.redhat.com/ubi9/ubi-minimal:9.4

FROM registry.ci.openshift.org/ocp/4.17:base-rhel9

ARG DATA_DIR=/data
RUN mkdir $DATA_DIR && chmod 775 $DATA_DIR

#RUN microdnf update
#RUN microdnf module enable postgresql:13
#RUN microdnf module enable nmstate:2.2
#RUN microdnf update -y && microdnf install --enablerepo=codeready-builder-for-rhel-9-x86_64-rpms -y nmstate-libs /usr/bin/nmstatectl && microdnf clean all

RUN dnf install -y nmstate-libs nmstate && dnf clean all && rm -rf /var/cache/dnf/*

WORKDIR /
COPY --from=builder /opt/app-root/src/build/manager /usr/local/bin/
COPY --from=builder /opt/app-root/src/build/server /usr/local/bin/
USER 65532:65532

ENTRYPOINT ["/usr/local/bin/manager"]
