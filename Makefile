all: binary

#vendor:
#	glide up -v

binary: #vendor
	./hack/build/build.sh ${VERSION}

release: binary
	[ -d _out ] && rm -rf _out
	mkdir -p _out
	cp cmd/kubevirt-template-validator/kubevirt-template-validator _out/kubevirt-template-validator-${VERSION}-linux-amd64

clean:
	rm -f cmd/kubevirt-template-validator/kubevirt-template-validator
	rm -rf _out

.PHONY: all vendor binary release clean

