.PHONY: test
test:
	@echo "🛠️  Running unit tests..."
	go test ./...

.PHONY: build
build:
	@echo "🔧  Building Go binaries..."
	go build -o bin/admission-webhook-linux-amd64 .

.PHONY: podman-build
podman-build:
	@echo "📦 Building simple-kubernetes-webhook podman image..."
	podman build -t simple-kubernetes-webhook:1.0 .

.PHONY: cluster
cluster:
	@echo "🔧 Creating Kubernetes cluster..."
	kind create cluster --config dev/manifests/kind/kind.cluster.yaml

.PHONY: delete-cluster
delete-cluster:
	@echo "♻️  Deleting Kubernetes cluster..."
	kind delete cluster

.PHONY: push
push: podman-build
	@echo "📦 Pushing admission-webhook image into Kind's podman daemon..."
	rm -f image.tar
	podman save simple-kubernetes-webhook:1.0 -o image.tar
	kind load image-archive image.tar

.PHONY: deploy-config
deploy-config:
	@echo "⚙️  Applying cluster config..."
	kubectl apply -f dev/manifests/cluster-config/

.PHONY: delete-config
delete-config:
	@echo "♻️  Deleting Kubernetes cluster config..."
	kubectl delete -f dev/manifests/cluster-config/

.PHONY: deploy
deploy: push delete deploy-config
	@echo "🚀 Deploying simple-kubernetes-webhook..."
	kubectl apply -f dev/manifests/webhook/

.PHONY: delete
delete:
	@echo "♻️  Deleting simple-kubernetes-webhook deployment if existing..."
	kubectl delete -f dev/manifests/webhook/ || true

.PHONY: pod
pod:
	@echo "🚀 Deploying test pod..."
	kubectl apply -f dev/manifests/pods/lifespan-seven.pod.yaml

.PHONY: delete-pod
delete-pod:
	@echo "♻️ Deleting test pod..."
	kubectl delete -f dev/manifests/pods/lifespan-seven.pod.yaml

.PHONY: bad-pod
bad-pod:
	@echo "🚀 Deploying \"bad\" pod..."
	kubectl apply -f dev/manifests/pods/bad-name.pod.yaml

.PHONY: delete-bad-pod
delete-bad-pod:
	@echo "🚀 Deleting \"bad\" pod..."
	kubectl delete -f dev/manifests/pods/bad-name.pod.yaml

.PHONY: taint
taint:
	@echo "🎨 Taining Kubernetes node.."
	kubectl taint nodes kind-control-plane "acme.com/lifespan-remaining"=4:NoSchedule

.PHONY: logs
logs:
	@echo "🔍 Streaming simple-kubernetes-webhook logs..."
	kubectl logs -l app=simple-kubernetes-webhook -f

.PHONY: delete-all
delete-all: delete delete-config delete-pod delete-bad-pod
