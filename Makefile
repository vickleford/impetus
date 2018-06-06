container_name := impetus
# container_registry := quay.io/nordstrom
# container_release := 0.1

.PHONY: build/image tag/image push/image

build:
	docker build \
                -f Dockerfile \
		--build-arg HTTP_PROXY=${HTTP_PROXY} --build-arg HTTPS_PROXY=${HTTPS_PROXY} \
		--build-arg http_proxy=${HTTP_PROXY} --build-arg https_proxy=${HTTPS_PROXY} \
		-t $(container_name) .

run: build
	docker run \
        -v ${PWD}/artifacts:/artifacts \
        --rm \
        $(container_name)

