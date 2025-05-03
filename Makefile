
.PHONY: preload
preload:
	cp ./hack/preload-extras.txt ./hack/preload.txt
	oc get pod -A -o json \
      | jq -r '.items[].spec.containers[].image' \
      | grep -v registry.k8s.io \
      | grep -v docker.io/kindest \
      | sort \
      | uniq \
      >> ./hack/preload.txt

.PHONY: start
start:
	./hack/run-in-kind.sh

.PHONY: namespace
namespace:
	oc get ns app-cider || oc create ns app-cider
	oc project app-cider

app-deploy: namespace
	helm -n app-cider upgrade -i cider app/helm -f app/helm/values-deploy.yaml

app-rollout: namespace
	helm -n app-cider upgrade -i cider app/helm -f app/helm/values-rollout.yaml

app-rollout-metrics: namespace
	helm -n app-cider upgrade -i cider app/helm-metrics

app-rollout-traffic: namespace
	helm -n app-cider upgrade -i cider app/helm-traffic-routing

app-rollout-get:
	kubectl argo rollouts get rollout cider

app-rollout-status:
	kubectl argo rollouts status cider
