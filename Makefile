REGISTRY ?= quay.io/myan
IMAGE_TAG ?= latest

build-lifecycle-image: 
	docker build -t ${REGISTRY}/lifecycle:${IMAGE_TAG} . -f lifecycle/Dockerfile