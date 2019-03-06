all: binary

docker: binary
	docker build .

dockertag: binary
	./hack/container/docker-tag.sh

dockerpush: binary
	./hack/container/docker-push.sh

vendor:
	dep ensure

binary: vendor
	./hack/build/build.sh

clean:
	rm -f cmd/kubevirt-template-validator/kubevirt-template-validator

.PHONY: all docker dockertag dockerpush binary clean

