CLUSTER_NAME ?= kind
KUBECTL = kubectl --context kind-$(CLUSTER_NAME)

.PHONY: test
test:
	@echo "\n🛠️  Running unit tests..."
	go test ./...

.PHONY: build
build:
	@echo "\n🔧  Building Go binaries..."
	GOOS=darwin GOARCH=amd64 go build -o bin/admission-webhook-darwin-amd64 .
	GOOS=linux GOARCH=amd64 go build -o bin/admission-webhook-linux-amd64 .
	GOOS=darwin GOARCH=arm64 go build -o bin/admission-webhook-darwin-arm64 .
	GOOS=linux GOARCH=arm64 go build -o bin/admission-webhook-linux-arm64 .

.PHONY: docker-build
docker-build:
	@echo "\n📦 Building simple-kubernetes-webhook Docker image..."
	docker build -t simple-kubernetes-webhook:latest .

# From this point `kind` is required
.PHONY: cluster
cluster:
	@echo "\n🔧 Creating Kubernetes cluster..."
	kind create cluster --name $(CLUSTER_NAME) --config dev/manifests/kind/kind.cluster.yaml --wait 120s

.PHONY: delete-cluster
delete-cluster:
	@echo "\n♻️  Deleting Kubernetes cluster..."
	kind delete cluster --name $(CLUSTER_NAME)

.PHONY: push
push: docker-build
	@echo "\n📦 Pushing admission-webhook image into Kind's Docker daemon..."
	kind load docker-image simple-kubernetes-webhook:latest --name $(CLUSTER_NAME)

.PHONY: certs
certs:
	@echo "\n🔐 Generating local TLS certificates and manifests..."
	bash dev/gen-certs.sh

.PHONY: deploy-config
deploy-config:
	@echo "\n⚙️  Applying cluster config..."
	$(KUBECTL) apply -f dev/manifests/cluster-config/apps.ns.yaml
	$(KUBECTL) apply -f dev/generated/mutating.config.yaml -f dev/generated/validating.config.yaml

.PHONY: delete-config
delete-config:
	@echo "\n♻️  Deleting Kubernetes cluster config..."
	$(KUBECTL) delete -f dev/manifests/cluster-config/ --ignore-not-found

.PHONY: deploy
deploy: push certs
	$(MAKE) delete
	@echo "\n🚀 Deploying simple-kubernetes-webhook..."
	$(KUBECTL) apply -f dev/generated/webhook.tls.secret.yaml
	$(KUBECTL) apply -f dev/manifests/webhook/
	$(KUBECTL) -n default rollout status deployment/simple-kubernetes-webhook --timeout=120s
	$(MAKE) deploy-config

.PHONY: delete
delete:
	@echo "\n♻️  Deleting simple-kubernetes-webhook deployment if existing..."
	$(KUBECTL) delete -f dev/manifests/webhook/ --ignore-not-found
	$(KUBECTL) -n default delete secret simple-kubernetes-webhook-tls --ignore-not-found

.PHONY: pod
pod:
	@echo "\n🚀 Deploying test pod..."
	$(KUBECTL) apply -f dev/manifests/pods/lifespan-seven.pod.yaml

.PHONY: delete-pod
delete-pod:
	@echo "\n♻️ Deleting test pod..."
	$(KUBECTL) delete -f dev/manifests/pods/lifespan-seven.pod.yaml --ignore-not-found

.PHONY: bad-pod
bad-pod:
	@echo "\n🚀 Deploying \"bad\" pod..."
	$(KUBECTL) apply -f dev/manifests/pods/bad-name.pod.yaml

.PHONY: delete-bad-pod
delete-bad-pod:
	@echo "\n🚀 Deleting \"bad\" pod..."
	$(KUBECTL) delete -f dev/manifests/pods/bad-name.pod.yaml --ignore-not-found

.PHONY: taint
taint:
	@echo "\n🎨 Tainting Kubernetes node..."
	$(KUBECTL) taint nodes $(CLUSTER_NAME)-control-plane "acme.com/lifespan-remaining"=4:NoSchedule --overwrite

.PHONY: logs
logs:
	@echo "\n🔍 Streaming simple-kubernetes-webhook logs..."
	$(KUBECTL) -n default logs -l app=simple-kubernetes-webhook -f

.PHONY: delete-all
delete-all:
	$(MAKE) delete-pod delete-bad-pod
	$(MAKE) delete
	$(MAKE) delete-config
