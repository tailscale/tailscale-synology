current_dir := $(realpath -s $(dir $(abspath $(lastword $(MAKEFILE_LIST)))))

.PHONY: build build-image

build-image:
	docker build . -t synobuild

build: build-image
	docker run --rm --privileged \
		--env PACKAGE_ARCH=$(SYNO_PLATFORM) \
		--env DSM_VER=$(SYNO_DSM_VERSION) \
		-v $(current_dir)/toolkit:/toolkit_tarballs \
		-v $(current_dir)/artifacts:/result_spk \
		synobuild

shell: build-image
	docker run --rm -it --privileged \
		--entrypoint /bin/bash \
		--env PACKAGE_ARCH=$(SYNO_PLATFORM) \
		--env DSM_VER=$(SYNO_DSM_VERSION) \
		-v $(current_dir)/toolkit:/toolkit_tarballs \
		-v $(current_dir)/artifacts:/result_spk \
		synobuild
