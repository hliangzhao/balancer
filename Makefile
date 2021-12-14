SHELL=/bin/bash -o pipefail

FIRST_GOPATH:=$(firstword $(subst :, ,$(shell go env GOPATH)))
GO_PKG=github.com/hliangzhao/balancer
K8S_GEN_BINARIES:=deepcopy-gen informer-gen lister-gen client-gen

TYPES_v1alpha1_TARGET:=pkg/apis/balancer/v1alpha1/balancer.go

K8S_GEN_DEPS:=.header
K8S_GEN_DEPS+=$(TYPES_v1alpha1_TARGET)
K8S_GEN_DEPS+=$(foreach bin,$(K8S_GEN_BINARIES),$(FIRST_GOPATH)/bin/$(bin))
K8S_GEN_DEPS+=$(OPENAPI_GEN_BINARY)

# TODO: 待修改
OPERATOR_E2E_IMAGE_TAG:=$(shell tar -cf - pkg | md5)
OPERATOR_E2E_IMAGE_NAME:=draveness/proxier-e2e:$(OPERATOR_E2E_IMAGE_TAG)

.PHONY: test
test:
	go test -count=1 ./pkg/...

e2e:
	./hack/docker-image-exists.sh || \
	(operator-sdk build $(OPERATOR_E2E_IMAGE_NAME) && docker push $(OPERATOR_E2E_IMAGE_NAME))
	go test -v ./test/e2e/ --kubeconfig "$(HOME)/.kube/k8s-playground-kubeconfig.yaml" --operator-image $(OPERATOR_E2E_IMAGE_NAME)

run:
	operator-sdk up local --namespace=default

LISTER_TARGET := pkg/client/listers/balancer/v1alpha1/balancer.go
$(LISTER_TARGET): $(K8S_GEN_DEPS)
	$(LISTER_GEN_BINARY) \
	$(K8S_GEN_ARGS) \
	--input-dirs     "$(GO_PKG)/pkg/apis/balancer/v1alpha1" \
	--output-package "$(GO_PKG)/pkg/client/listers"

# generate listers
CLIENT_TARGET := pkg/client/versioned/clientset.go
$(CLIENT_TARGET): $(K8S_GEN_DEPS)
	$(CLIENT_GEN_BINARY) \
	$(K8S_GEN_ARGS) \
	--input-base     "" \
	--clientset-name "versioned" \
	--input	         "$(GO_PKG)/pkg/apis/balancer/v1alpha1" \
	--output-package "$(GO_PKG)/pkg/client"

# generate informers
INFORMER_TARGET := pkg/client/informers/externalversions/balancer/v1alpha1/balancer.go
$(INFORMER_TARGET): $(K8S_GEN_DEPS) $(LISTER_TARGET) $(CLIENT_TARGET)
	$(INFORMER_GEN_BINARY) \
	$(K8S_GEN_ARGS) \
	--versioned-clientset-package "$(GO_PKG)/pkg/client/versioned" \
	--listers-package "$(GO_PKG)/pkg/client/listers" \
	--input-dirs      "$(GO_PKG)/pkg/apis/balancer/v1alpha1" \
	--output-package  "$(GO_PKG)/pkg/client/informers"

.PHONY: k8s-gen
k8s-gen: \
  $(CLIENT_TARGET) \
  $(LISTER_TARGET) \
  $(INFORMER_TARGET)

# TODO: 待修改
.PHONY: release
release:
	./hack/make-release.sh
	operator-sdk build draveness/proxier:$(RELEASE)
	docker push draveness/proxier:$(RELEASE)
	git tag $(RELEASE)

define _K8S_GEN_VAR_TARGET_
$(shell echo $(1) | tr '[:lower:]' '[:upper:]' | tr '-' '_')_BINARY:=$(FIRST_GOPATH)/bin/$(1)

$(FIRST_GOPATH)/bin/$(1):
	go get -u -d k8s.io/code-generator/cmd/$(1)
	cd $(FIRST_GOPATH)/src/k8s.io/code-generator; git checkout $(K8S_GEN_VERSION)
	go install k8s.io/code-generator/cmd/$(1)

endef

$(foreach binary,$(K8S_GEN_BINARIES),$(eval $(call _K8S_GEN_VAR_TARGET_,$(binary))))
