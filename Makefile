REGISTRY = default-route-openshift-image-registry.apps.eu-west-1-a5g8f.eng.msp.worldpay.io
VERSION = 0.0.1
NAMESPACE = default
NAME = memcached

.PHONY: imagestream-init
imagestream-init:
	docker login --username=`oc whoami` --password=`oc whoami -t` https://$(REGISTRY)
	oc create imagestream $(NAME) --lookup-local=true

.PHONY: run
run:
	operator-sdk run --local --namespace default

.PHONY: build
build:
	operator-sdk build $(REGISTRY)/$(NAMESPACE)/$(NAME):v$(VERSION)
	docker push $(REGISTRY)/$(NAMESPACE)/$(NAME):v$(VERSION)

.PHONY: code-gen
code-gen: ## Run the operator-sdk commands to generated code (k8s and openapi)
	@echo Updating the deep copy files with the changes in the API
	operator-sdk generate k8s
	@echo Updating the CRD files with the OpenAPI validations
	operator-sdk generate crds


.PHONY: test-unit
test-unit:
	go test -cover ./pkg/controller/...

.PHONY: test-e2e
test-e2e:
	operator-sdk test local --namespace=$(NAMESPACE) --up-local --verbose ./test/e2e/
