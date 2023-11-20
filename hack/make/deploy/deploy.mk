ENABLE_CSI ?= true
PLATFORM ?= "kubernetes"

## Deploy the operator without the csi-driver, with platform specified in % (kubernetes or openshift)
deploy/%/no-csi:
	@make ENABLE_CSI=false $(@D)

## Deploy the operator with csi-driver, with platform specified in % (kubernetes or openshift)
deploy/%:
	@make PLATFORM=$(@F) $(@D)

## Deploy the operator with csi-driver, on kubernetes
deploy-orig: manifests/crd/helm
	kubectl get namespace dynatrace || kubectl create namespace dynatrace
	helm template dynatrace-operator config/helm/chart/default \
			--namespace dynatrace \
			--set installCRD=true \
			--set platform=$(PLATFORM) \
			--set csidriver.enabled=$(ENABLE_CSI) \
			--set manifests=true \
			--set image="$(IMAGE_URI)" | kubectl apply -f -

deploy: manifests/crd/helm
	kubectl get namespace dynatrace || kubectl create namespace dynatrace
	helm template dynatrace-operator config/helm/chart/default \
			--namespace dynatrace \
			--set installCRD=true \
			--set platform=$(PLATFORM) \
			--set csidriver.enabled=$(ENABLE_CSI) \
			--set manifests=true \
			--set image="$(IMAGE_URI)" \
			--set operator.tolerations[0].key="injection_failure_policy",operator.tolerations[0].operator="Exists",operator.tolerations[0].effect="NoExecute" \
			--set webhook.tolerations[0].key="injection_failure_policy",webhook.tolerations[0].operator="Exists",webhook.tolerations[0].effect="NoExecute" \
			--set csidriver.tolerations[0].key="node-role.kubernetes.io/master",csidriver.tolerations[0].operator="Exists",csidriver.tolerations[0].effect="NoSchedule" \
			--set csidriver.tolerations[1].key="node-role.kubernetes.io/control-plane",csidriver.tolerations[1].operator="Exists",csidriver.tolerations[1].effect="NoSchedule" \
			--set csidriver.tolerations[2].key="injection_failure_policy",csidriver.tolerations[2].operator="Exists",csidriver.tolerations[2].effect="NoExecute" | kubectl apply -f -

