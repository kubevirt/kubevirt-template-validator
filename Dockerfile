FROM centos:7

MAINTAINER "Francesco Romani" <fromani@redhat.com>
ENV container docker

COPY cmd/kubevirt-template-validator/kubevirt-template-validator /usr/sbin/kubevirt-template-validator

ENTRYPOINT [ "/usr/sbin/kubevirt-template-validator" ]
