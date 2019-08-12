# It's necessary to set this because some environments don't link sh -> bash.
SHELL := /bin/bash

#-----------------------------------------------------------------------------
# VERBOSE target
#-----------------------------------------------------------------------------

# When you run make VERBOSE=1 (the default), executed commands will be printed
# before executed. If you run make VERBOSE=2 verbose flags are turned on and
# quiet flags are turned off for various commands. Use V_FLAG in places where
# you can toggle on/off verbosity using -v. Use Q_FLAG in places where you can
# toggle on/off quiet mode using -q. Use S_FLAG where you want to toggle on/off
# silence mode using -s...
VERBOSE ?= 1
Q = @
Q_FLAG = -q
QUIET_FLAG = --quiet
V_FLAG =
VERBOSE_FLAG = 
S_FLAG = -s
X_FLAG =
ifeq ($(VERBOSE),1)
	Q =
endif
ifeq ($(VERBOSE),2)
	Q =
	Q_FLAG =
	QUIET_FLAG =
	S_FLAG =
	V_FLAG = -v
	VERBOSE_FLAG = --verbose
	X_FLAG = -x
endif

PROJECT_NAME ?= aws-rds
OPERATOR_NAME ?= $(PROJECT_NAME)-operator
OPERATOR_VERSION ?= $(shell cat ./VERSION)

IMAGE_NAME ?= quay.io/$(QUAY_USERNAME)/$(PROJECT_NAME)
ifeq ($(RELEASE_OPERATOR), true)
	TAG ?= $(OPERATOR_VERSION)
else
	TAG ?= $(OPERATOR_VERSION)-dev
endif
IMAGE ?= $(IMAGE_NAME):$(TAG)

OPERATOR_NS ?= openshift-operators

TMP_DIR := tmp
TEMPLATES_DIR := ./templates
MANIFESTS_DIR := ./$(TMP_DIR)/manifests
BUILD_TIMESTAMP = build-timestamp

ifeq ($(RELEASE_OPERATOR), true)
	APPR_NAMESPACE ?= $(QUAY_USERNAME)
else
	APPR_NAMESPACE ?= $(QUAY_USERNAME)-testing
endif

.DEFAULT_GOAL := help

## -- Utility targets --

## Print help message for all Makefile targets
## Run `make` or `make help` to see the help
.PHONY: help
help: ## Credit: https://gist.github.com/prwhite/8168133#gistcomment-2749866

	@printf "Usage:\n  make <target>";

	@awk '{ \
			if ($$0 ~ /^.PHONY: [a-zA-Z\-\_0-9]+$$/) { \
				helpCommand = substr($$0, index($$0, ":") + 2); \
				if (helpMessage) { \
					printf "\033[36m%-20s\033[0m %s\n", \
						helpCommand, helpMessage; \
					helpMessage = ""; \
				} \
			} else if ($$0 ~ /^[a-zA-Z\-\_0-9.]+:/) { \
				helpCommand = substr($$0, 0, index($$0, ":")); \
				if (helpMessage) { \
					printf "\033[36m%-20s\033[0m %s\n", \
						helpCommand, helpMessage; \
					helpMessage = ""; \
				} \
			} else if ($$0 ~ /^##/) { \
				if (helpMessage) { \
					helpMessage = helpMessage"\n                     "substr($$0, 3); \
				} else { \
					helpMessage = substr($$0, 3); \
				} \
			} else { \
				if (helpMessage) { \
					print "\n                     "helpMessage"\n" \
				} \
				helpMessage = ""; \
			} \
		}' \
		$(MAKEFILE_LIST)

.PHONY: dep
## Runs 'dep ensure -v'
dep:
	$(Q)dep ensure $(V_FLAG)

## -- Build targets --


.PHONY: refresh-timestamp
refresh-timestamp:
	$(eval NEW_TIMESTAMP := $(shell date +%s))
	@echo -n "$(NEW_TIMESTAMP)"  > $(BUILD_TIMESTAMP)
	@echo "Refreshing build timestamp to '$(NEW_TIMESTAMP)' ($(shell date --date="@$(NEW_TIMESTAMP)" '+%Y-%m-%d %H:%M:%S'))"

.PHONY: get-timestamp
ifeq "$(shell test -s $(BUILD_TIMESTAMP) && echo yes)" "yes"
get-timestamp:
else
get-timestamp: refresh-timestamp
endif
	$(eval export TIMESTAMP = $(shell cat $(BUILD_TIMESTAMP)))
	$(eval export READABLE_TIMESTAMP := $(shell date --date="@$(TIMESTAMP)" '+%Y%m%d%H%M%S'))
	$(eval export NICE_READABLE_TIMESTAMP := $(shell date --date="@$(TIMESTAMP)" '+%Y-%m-%d %H:%M:%S'))

.PHONY: build
## Compile the operator for Linux/AMD64
build: refresh-timestamp
	$(Q)GO111MODULE=off go build $(V_FLAG)

.PHONY: build-image
## Build the operator image
build-image: build
	$(Q)podman build -t $(IMAGE) .

.PHONY: push-image
## Push the operator image to quay.io
push-image:
	@podman login -u "$(QUAY_USERNAME)" -p "$(QUAY_PASSWORD)" quay.io
	$(Q)podman push $(IMAGE)

.PHONY: build-operator-csv
## Build the operator ClusterServiceVersion
build-operator-csv: get-timestamp
	$(Q)mkdir -p $(MANIFESTS_DIR)
	$(Q)sed -e 's,REPLACE_OPERATOR_NAME,$(OPERATOR_NAME),g' $(TEMPLATES_DIR)/clusterserviceversion.yaml | \
		sed -e 's,REPLACE_VERSION,$(TAG),g' | \
		sed -e 's,REPLACE_IMAGE,$(IMAGE),g' | \
		sed -e 's,REPLACE_CREATED_AT,$(NICE_READABLE_TIMESTAMP),g' | \
		sed -e 's,REPLACE_CSV_NAMESPACE,placeholder,g' | \
		sed -e 's,REPLACE_PACKAGE,$(OPERATOR_NAME),g' > $(MANIFESTS_DIR)/$(OPERATOR_NAME)-v$(TAG).clusterserviceversion.yaml

.PHONY: build-operator-olm-package
## Build the operator package for OpenShift Marketplace
build-operator-olm-package: get-timestamp build-operator-csv
	$(Q)mkdir -p $(MANIFESTS_DIR)
	$(Q)sed -e 's,REPLACE_OPERATOR_NAME,$(OPERATOR_NAME),g' $(TEMPLATES_DIR)/package.yaml | \
		sed -e 's,REPLACE_VERSION,$(TAG),g' | \
		sed -e 's,REPLACE_PACKAGE,$(OPERATOR_NAME),g' > $(MANIFESTS_DIR)/$(OPERATOR_NAME).package.yaml
	$(Q)cp -f $(TEMPLATES_DIR)/crd.yaml $(MANIFESTS_DIR)/aws-v1alpha1-rdsdatabases.crd.yaml
	$(Q)operator-courier verify --ui_validate_io $(MANIFESTS_DIR)

.PHONY: push-operator-olm-package
## Pushes the operator package to Quay.io app registry
push-operator-olm-package: get-timestamp
	@$(eval QUAY_API_TOKEN := $(shell curl -sH "Content-Type: application/json" -XPOST https://quay.io/cnr/api/v1/users/login -d '{"user":{"username":"'${QUAY_USERNAME}'","password":"'${QUAY_PASSWORD}'"}}' | jq -r '.token'))
	@echo "Pushing operator package $(APPR_NAMESPACE)/$(OPERATOR_NAME):$(TAG)-$(READABLE_TIMESTAMP)"
	@operator-courier $(VERBOSE_FLAG) push $(MANIFESTS_DIR) $(APPR_NAMESPACE) $(OPERATOR_NAME) $(TAG)-$(READABLE_TIMESTAMP) "$(QUAY_API_TOKEN)"

.PHONY: clean
## Clean up 
clean:
	@rm -rvf $(PROJECT_NAME) tmp/

## -- Install/Delete targets --

.PHONY: install-operator-secrets
## Create secret for operator
install-operator-secrets:
	$(Q)sed -e 's,REPLACE_OPERATOR_NAME,$(OPERATOR_NAME),g' $(TEMPLATES_DIR)/aws.secret.yaml | oc apply -f -

.PHONY: uninstall-operator-secrets
## Delete secret for operator
uninstall-operator-secrets:
	$(Q)-sed -e 's,REPLACE_OPERATOR_NAME,$(OPERATOR_NAME),g' $(TEMPLATES_DIR)/aws.secret.yaml | oc delete -f -

.PHONY: install-operator
## Create secret, role, account and crd for operator
install-operator: install-operator-secrets
	$(Q)sed -e 's,REPLACE_OPERATOR_NAME,$(OPERATOR_NAME),g' $(TEMPLATES_DIR)/operator-cluster-role.yaml | oc apply -f -
	$(Q)sed -e 's,REPLACE_OPERATOR_NAME,$(OPERATOR_NAME),g' $(TEMPLATES_DIR)/operator-service-account.yaml | oc apply -f -
	$(Q)sed -e 's,REPLACE_OPERATOR_NAME,$(OPERATOR_NAME),g' $(TEMPLATES_DIR)/operator-cluster-role-binding.yaml | oc apply -f -
	$(Q)sed -e 's,REPLACE_OPERATOR_NAME,$(OPERATOR_NAME),g' $(TEMPLATES_DIR)/crd.yaml | oc apply -f -

.PHONY: uninstall-operator
## Delete secret, role, account and crd for operator
uninstall-operator: 
	$(Q)-sed -e 's,REPLACE_OPERATOR_NAME,$(OPERATOR_NAME),g' $(TEMPLATES_DIR)/operator-cluster-role.yaml | oc delete -f -
	$(Q)-sed -e 's,REPLACE_OPERATOR_NAME,$(OPERATOR_NAME),g' $(TEMPLATES_DIR)/operator-service-account.yaml | oc delete -f -
	$(Q)-sed -e 's,REPLACE_OPERATOR_NAME,$(OPERATOR_NAME),g' $(TEMPLATES_DIR)/operator-cluster-role-binding.yaml | oc delete -f -
	$(Q)-sed -e 's,REPLACE_OPERATOR_NAME,$(OPERATOR_NAME),g' $(TEMPLATES_DIR)/crd.yaml | oc delete -f -

.PHONY: install-olm-operator-source
## Create OperatorSource for operator
install-olm-operator-source:
	$(Q)sed -e 's,REPLACE_OPERATOR_NAME,$(OPERATOR_NAME),g' $(TEMPLATES_DIR)/operator-source.yaml | \
		sed -e 's,REPLACE_NAMESPACE,$(APPR_NAMESPACE),g' | oc apply -f -

.PHONY: uninstall-olm-operator-source
## Delete OperatorSource for operator
uninstall-olm-operator-source:
	$(Q)-sed -e 's,REPLACE_OPERATOR_NAME,$(OPERATOR_NAME),g' $(TEMPLATES_DIR)/operator-source.yaml | \
		sed -e 's,REPLACE_NAMESPACE,$(APPR_NAMESPACE),g' | oc delete -f -

.PHONY: deploy-operator
## Create deployment for operator
deploy-operator: install-operator build-operator-csv
	#$(Q)sed -e 's,REPLACE_OPERATOR_NAME,$(OPERATOR_NAME),g' $(TEMPLATES_DIR)/deployment.yaml | \
	#	sed -e 's,REPLACE_IMAGE,$(IMAGE),g' | oc apply -f -
	$(Q) sed -e 's,namespace: placeholder,namespace: openshift-operators,g' $(MANIFESTS_DIR)/$(OPERATOR_NAME)-v$(TAG).clusterserviceversion.yaml | oc apply -f -

.PHONY: redeploy-operator
## Scale operator's deployment to 0 and back to 1 to re-deploy it
redeploy-operator:
	$(Q)oc scale --replicas=0 deploy $(OPERATOR_NAME) -n $(OPERATOR_NS) \
	&& oc scale --replicas=1 deploy $(OPERATOR_NAME) -n $(OPERATOR_NS)

.PHONY: undeploy-operator
## Delete deployment for operator
undeploy-operator: uninstall-operator
	#$(Q)-sed -e 's,REPLACE_OPERATOR_NAME,$(OPERATOR_NAME),g' $(TEMPLATES_DIR)/deployment.yaml | \
	#	sed -e 's,REPLACE_IMAGE,$(IMAGE),g' | oc delete -f -
	$(Q)-sed -e 's,namespace: placeholder,namespace: openshift-operators,g' $(MANIFESTS_DIR)/$(OPERATOR_NAME)-v$(TAG).clusterserviceversion.yaml | oc delete -f -

.PHONY: deploy-db
## Create database secret and deployment
deploy-db:
	$(Q)sed -e 's,REPLACE_OPERATOR_NAME,$(OPERATOR_NAME),g' $(TEMPLATES_DIR)/db.secret.yaml | oc apply -f -
	$(Q)sed -e 's,REPLACE_OPERATOR_NAME,$(OPERATOR_NAME),g' $(TEMPLATES_DIR)/db.yaml | oc apply -f -

.PHONY: undeploy-db
## Delete database secret, deployment and service
undeploy-db:
	$(Q)-sed -e 's,REPLACE_OPERATOR_NAME,$(OPERATOR_NAME),g' $(TEMPLATES_DIR)/db.yaml | oc delete -f -
	$(Q)-sed -e 's,REPLACE_OPERATOR_NAME,$(OPERATOR_NAME),g' $(TEMPLATES_DIR)/db.secret.yaml | oc delete -f -

.PHONY: undeploy-all
## Undeploy operator and related assets
undeploy-all: undeploy-db undeploy-operator

## -- Test targets --

.PHONY: test-unit
## Run the operator's unit tests
test-unit:
	$(Q)go test github.com/operator-backing-service-samples/aws-rds/pkg/{rds,crd} $(V_FLAG)

## -- Run targets --

.PHONY: run-locally
## Run the operator locally
run-locally:
	$(Q)./$(PROJECT_NAME)
