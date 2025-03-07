
.PHONY: preload
preload:
	oc get pod -A -o json \
      | jq -r '.items[].spec.containers[].image' \
      | grep -v registry.k8s.io \
      | grep -v docker.io/kindest \
      | sort \
      | uniq \
      > ./hack/preload.txt

.PHONY: start
start:
	./hack/run-in-kind.sh

app-deploy:
	helm -n app-cider upgrade -i cider app/helm -f app/helm/values-deploy.yaml

app-rollout:
	helm -n app-cider upgrade -i cider app/helm -f app/helm/values-rollout.yaml

app-rollout-metrics:
	helm -n app-cider upgrade -i cider app/helm -f app/helm/values-rollout.yaml --set rollout.metrics.use=true
