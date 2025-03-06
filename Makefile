
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
