FROM centos:7

MAINTAINER "Francesco Romani" <fromani@redhat.com>
ENV container docker

RUN yum -y update
RUN yum install -y net-tools
RUN yum clean all

RUN mkdir -p /etc/webhook/certs
COPY cmd/kubevirt-template-validator/kubevirt-template-validator /usr/sbin/kubevirt-template-validator

ENTRYPOINT [ "/usr/sbin/kubevirt-template-validator" ]
