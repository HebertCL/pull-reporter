# Check to see if we can use ash, in Alpine images, or default to BASH.
SHELL_PATH = /bin/ash
SHELL = $(if $(wildcard $(SHELL_PATH)),/bin/ash,/bin/bash)

# ==============================================================================
# Define dependencies

GOLANG          := golang:1.20
ALPINE          := alpine:3.18
KIND            := kindest/node:v1.27.1
POSTGRES        := postgres:15.3
VAULT           := hashicorp/vault:1.13
ZIPKIN          := openzipkin/zipkin:2.24
TELEPRESENCE    := datawire/tel2:2.13.1

KIND_CLUSTER    := hebert-starter-cluster
NAMESPACE       := pr-manager-system
APP             := prm
BASE_IMAGE_NAME := hebertcl/service
SERVICE_NAME    := pr-manager
VERSION         := 0.0.1
SERVICE_IMAGE   := $(BASE_IMAGE_NAME)/$(SERVICE_NAME):$(VERSION)

# VERSION       := "0.0.1-$(shell git rev-parse --short HEAD)"

# ==============================================================================
# Running from within k8s/kind

dev-bill:
	kind load docker-image $(POSTGRES) --name $(KIND_CLUSTER)

dev-up-local:
	kind create cluster \
		--image $(KIND) \
		--name $(KIND_CLUSTER) \
		--config zarf/k8s/dev/kind-config.yaml

	kubectl wait --timeout=120s --namespace=local-path-storage --for=condition=Available deployment/local-path-provisioner

	kind load docker-image $(TELEPRESENCE) --name $(KIND_CLUSTER)
	kind load docker-image $(POSTGRES) --name $(KIND_CLUSTER)

dev-up: dev-up-local

dev-down-local:
	kind delete cluster --name $(KIND_CLUSTER)

dev-down:
	kind delete cluster --name $(KIND_CLUSTER)

# ==============================================================================

dev-status:
	kubectl get nodes -o wide
	kubectl get svc -o wide
	kubectl get pods -o wide --watch --all-namespaces

# ==============================================================================

run-local:
	go run app/services/pull-manager/main.go

tidy:
	go mod tidy
	go mod vendor