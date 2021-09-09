FROM registry.access.redhat.com/ubi8/ubi-minimal as builder

RUN microdnf install -y golang-1.15.* && microdnf clean all

ARG VERSION=latest
ARG COMPONENT="kubevirt-template-validator"
ARG BRANCH=master
ARG REVISION=master

WORKDIR /workspace

# Copy the Go Modules manifests and vendor directory
COPY go.mod go.mod
COPY go.sum go.sum
COPY vendor/ vendor/

# Copy the go source
COPY cmd/ cmd/
COPY pkg/ pkg/
COPY internal/ internal/

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -a -ldflags="-X 'kubevirt.io/ssp-operator/internal/template-validator/version.COMPONENT=$COMPONENT'\
-X 'kubevirt.io/ssp-operator/internal/template-validator/version.VERSION=$VERSION'\
-X 'kubevirt.io/ssp-operator/internal/template-validator/version.BRANCH=$BRANCH'\
-X 'kubevirt.io/ssp-operator/internal/template-validator/version.REVISION=$REVISION'" -o cmd/kubevirt-template-validator/kubevirt-template-validator cmd/kubevirt-template-validator/main.go

FROM registry.access.redhat.com/ubi8/ubi-minimal

MAINTAINER "Francesco Romani" <fromani@redhat.com>
ENV container docker

RUN mkdir -p /etc/webhook/certs
WORKDIR /
COPY --from=builder /workspace/cmd/kubevirt-template-validator/kubevirt-template-validator /usr/sbin/kubevirt-template-validator

ENTRYPOINT [ "/usr/sbin/kubevirt-template-validator" ]
