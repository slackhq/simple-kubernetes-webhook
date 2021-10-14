# simple-kubernetes-webhook

This is a simple [Kubernetes admission webhook](https://kubernetes.io/docs/reference/access-authn-authz/extensible-admission-controllers/). It is meant to be used as a validating and mutating admission webhook only and does not support any controller logic. It has been developed as a simple Go web service without using any framework or boilerplate such as kubebuilder.

This project is aimed at illustrating how to build a fully functioning admission webhook in the simplest way possible. Most existing examples found on the web rely on heavy machinery using powerful frameworks, yet fail to illustrate how to implement a lightweight webhook that can do much needed actions such as rejecting a pod for compliance reasons, or inject helpful environment variables.

For readability, this project has been stripped of the usual production items such as: observability instrumentation, release scripts, redundant deployment configurations, etc. As such, it is not meant to use as-is in a production environment. This project is, in fact, a simplified fork of a system used accross all Kubernetes production environments at Slack.

## Installation
This project can fully run locally and includes automation to deploy a local Kubernetes cluster (using Kind).

### Requirements
* Docker
* kubectl
* [Kind](https://kind.sigs.k8s.io/docs/user/quick-start/#installation)
* Go >=1.16 (optional)

## Usage
### Create Cluster
First, we need to create a Kubernetes cluster:
```
❯ make cluster

🔧 Creating Kubernetes cluster...
kind create cluster --config dev/manifests/kind/kind.cluster.yaml
Creating cluster "kind" ...
 ✓ Ensuring node image (kindest/node:v1.21.1) 🖼
 ✓ Preparing nodes 📦
 ✓ Writing configuration 📜
 ✓ Starting control-plane 🕹️
 ✓ Installing CNI 🔌
 ✓ Installing StorageClass 💾
Set kubectl context to "kind-kind"
You can now use your cluster with:

kubectl cluster-info --context kind-kind

Have a nice day! 👋
```

Make sure that the Kubernetes node is ready:
```
❯ kubectl get nodes
NAME                 STATUS   ROLES                  AGE     VERSION
kind-control-plane   Ready    control-plane,master   3m25s   v1.21.1
```

And that system pods are running happily:
```
❯ kubectl -n kube-system get pods
NAME                                         READY   STATUS    RESTARTS   AGE
coredns-558bd4d5db-thwvj                     1/1     Running   0          3m39s
coredns-558bd4d5db-w85ks                     1/1     Running   0          3m39s
etcd-kind-control-plane                      1/1     Running   0          3m56s
kindnet-84slq                                1/1     Running   0          3m40s
kube-apiserver-kind-control-plane            1/1     Running   0          3m54s
kube-controller-manager-kind-control-plane   1/1     Running   0          3m56s
kube-proxy-4h6sj                             1/1     Running   0          3m40s
kube-scheduler-kind-control-plane            1/1     Running   0          3m54s
```

### Deploy Admission Webhook
To configure the cluster to use the admission webhook and to deploy said webhook, simply run:
```
❯ make deploy

📦 Building simple-kubernetes-webhook Docker image...
docker build -t simple-kubernetes-webhook:latest .
[+] Building 14.3s (13/13) FINISHED
...

📦 Pushing admission-webhook image into Kind's Docker daemon...
kind load docker-image simple-kubernetes-webhook:latest
Image: "simple-kubernetes-webhook:latest" with ID "sha256:46b8603bcc11a8fa1825190d3ed99c099096395b22a709e13ec6e7ae2f54014d" not yet present on node "kind-control-plane", loading...

⚙️  Applying cluster config...
kubectl apply -f dev/manifests/cluster-config/
namespace/apps created
mutatingwebhookconfiguration.admissionregistration.k8s.io/simple-kubernetes-webhook.acme.com created
validatingwebhookconfiguration.admissionregistration.k8s.io/simple-kubernetes-webhook.acme.com created

🚀 Deploying simple-kubernetes-webhook...
kubectl apply -f dev/manifests/webhook/
deployment.apps/simple-kubernetes-webhook created
service/simple-kubernetes-webhook created
secret/simple-kubernetes-webhook-tls created
```

Then, make sure the admission webhook pod is running (in the `default` namespace):
```
❯ kubectl get pods
NAME                                        READY   STATUS    RESTARTS   AGE
simple-kubernetes-webhook-77444566b7-wzwmx   1/1     Running   0          2m21s
```

You can stream logs from it:
```
❯ make logs

🔍 Streaming simple-kubernetes-webhook logs...
kubectl logs -l app=simple-kubernetes-webhook -f
time="2021-09-03T04:59:10Z" level=info msg="Listening on port 443..."
time="2021-09-03T05:02:21Z" level=debug msg=healthy uri=/health
```

And hit it's health endpoint from your local machine:
```
❯ curl -k https://localhost:8443/health
OK
```

### Deploying pods
Deploy a valid test pod that gets succesfully created:
```
❯ make pod

🚀 Deploying test pod...
kubectl apply -f dev/manifests/pods/lifespan-seven.pod.yaml
pod/lifespan-seven created
```
You should see in the admission webhook logs that the pod got mutated and validated.

Deploy a non valid pod that gets rejected:
```
❯ make bad-pod

🚀 Deploying "bad" pod...
kubectl apply -f dev/manifests/pods/bad-name.pod.yaml
Error from server: error when creating "dev/manifests/pods/bad-name.pod.yaml": admission webhook "simple-kubernetes-webhook.acme.com" denied the request: pod name contains "offensive"
```
You should see in the admission webhook logs that the pod validation failed. It's possible you will also see that the pod was mutated, as webhook configurations are not ordered.

## Testing
Unit tests can be run with the following command:
```
$ make test
go test ./...
?   	github.com/slackhq/simple-kubernetes-webhook	[no test files]
ok  	github.com/slackhq/simple-kubernetes-webhook/pkg/admission	0.611s
ok  	github.com/slackhq/simple-kubernetes-webhook/pkg/mutation	1.064s
ok  	github.com/slackhq/simple-kubernetes-webhook/pkg/validation	0.749s
```

## Admission Logic
A set of validations and mutations are implemented in an extensible framework. Those happen on the fly when a pod is deployed and no further resources are tracked and updated (ie. no controller logic).

### Validating Webhooks
#### Implemented
- [name validation](pkg/validation/name_validator.go): validates that a pod name doesn't contain any offensive string

#### How to add a new pod validation
To add a new pod mutation, create a file `pkg/validation/MUTATION_NAME.go`, then create a new struct implementing the `validation.podValidator` interface.

### Mutating Webhooks
#### Implemented
- [inject env](pkg/mutation/inject_env.go): inject environment variables into the pod such as `KUBE: true`
- [minimum pod lifespan](pkg/mutation/minimum_lifespan.go): inject a set of tolerations used to match pods to nodes of a certain age, the tolerations injected are controlled via the `acme.com/lifespan-requested` pod label.

#### How to add a new pod mutation
To add a new pod mutation, create a file `pkg/mutation/MUTATION_NAME.go`, then create a new struct implementing the `mutation.podMutator` interface.



