FROM centos:7

MAINTAINER "Francesco Romani" <fromani@redhat.com>
ENV container docker

RUN mkdir -p /etc/webhook/certs
RUN yum install -y net-tools
COPY cmd/kubevirt-template-validator/kubevirt-template-validator /usr/sbin/kubevirt-template-validator

ENTRYPOINT [ "/usr/sbin/kubevirt-template-validator" ]
