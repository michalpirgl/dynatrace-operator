
## Deploy the operator in a cluster configured in ~/.kube/config where platform and version are autodetected
deploy/helm-orig: manifests/crd/helm
	helm upgrade dynatrace-operator config/helm/chart/default \
			--install \
			--namespace dynatrace \
			--create-namespace \
			--atomic \
			--set installCRD=true \
			--set csidriver.enabled=$(ENABLE_CSI) \
			--set manifests=true \
			--set image="$(IMAGE_URI)"

deploy/helm: manifests/crd/helm
	helm upgrade dynatrace-operator config/helm/chart/default \
			--install \
			--namespace dynatrace \
			--create-namespace \
			--atomic \
			--set installCRD=true \
			--set csidriver.enabled=$(ENABLE_CSI) \
			--set manifests=true \
			--set image="$(IMAGE_URI)" \
			--set operator.tolerations[0].key="injection_failure_policy",operator.tolerations[0].operator="Exists",operator.tolerations[0].effect="NoExecute" \
			--set webhook.tolerations[0].key="injection_failure_policy",webhook.tolerations[0].operator="Exists",webhook.tolerations[0].effect="NoExecute" \
			--set csidriver.tolerations[0].key="node-role.kubernetes.io/master",csidriver.tolerations[0].operator="Exists",csidriver.tolerations[0].effect="NoSchedule" \
			--set csidriver.tolerations[1].key="node-role.kubernetes.io/control-plane",csidriver.tolerations[1].operator="Exists",csidriver.tolerations[1].effect="NoSchedule" \
			--set csidriver.tolerations[2].key="injection_failure_policy",csidriver.tolerations[2].operator="Exists",csidriver.tolerations[2].effect="NoExecute"

## Undeploy the operator in a cluster configured in ~/.kube/config where platform and k8s version are autodetected
undeploy/helm:
	helm uninstall dynatrace-operator \
			--namespace dynatrace
