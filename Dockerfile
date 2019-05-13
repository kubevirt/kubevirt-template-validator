FROM registry.access.redhat.com/ubi8/ubi-minimal

MAINTAINER "Francesco Romani" <fromani@redhat.com>
ENV container docker

RUN mkdir -p /etc/webhook/certs
COPY cmd/kubevirt-template-validator/kubevirt-template-validator /usr/sbin/kubevirt-template-validator

ENTRYPOINT [ "/usr/sbin/kubevirt-template-validator" ]
