
VERSION                                := $(shell cat VERSION)
REGISTRY                               := eu.gcr.io/gardener-project/gardener
PREFIX                                 := ext-authz-server
EXTERNAL_AUTHZ_SERVER_IMAGE_REPOSITORY := $(REGISTRY)/$(PREFIX)
EXTERNAL_AUTHZ_SERVER_IMAGE_TAG        := $(VERSION)

.PHONY: ext-authz-server-docker-image
ext-authz-server-docker-image:
	@docker build -t $(EXTERNAL_AUTHZ_SERVER_IMAGE_REPOSITORY):$(EXTERNAL_AUTHZ_SERVER_IMAGE_TAG) -f Dockerfile --rm .

.PHONY: docker-images
docker-images: ext-authz-server-docker-image

.PHONY: release
release: docker-images docker-login docker-push

.PHONY: docker-login
docker-login:
	@gcloud auth activate-service-account --key-file .kube-secrets/gcr/gcr-readwrite.json

.PHONY: docker-push
docker-push:
	@if ! docker images $(EXTERNAL_AUTHZ_SERVER_IMAGE_REPOSITORY) | awk '{ print $$2 }' | grep -q -F $(EXTERNAL_AUTHZ_SERVER_IMAGE_TAG); then echo "$(EXTERNAL_AUTHZ_SERVER_IMAGE_REPOSITORY) version $(EXTERNAL_AUTHZ_SERVER_IMAGE_TAG) is not yet built. Please run 'ext-authz-server-docker-image'"; false; fi
	@docker -- push $(EXTERNAL_AUTHZ_SERVER_IMAGE_REPOSITORY):$(EXTERNAL_AUTHZ_SERVER_IMAGE_TAG)
