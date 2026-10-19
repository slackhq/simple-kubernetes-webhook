# simple-kubernetes-webhook

This is a simple [Kubernetes admission webhook](https://kubernetes.io/docs/reference/access-authn-authz/extensible-admission-controllers/). It is meant to be used as a validating and mutating admission webhook only and does not support any controller logic. It has been developed as a simple Go web service without using any framework or boilerplate such as kubebuilder.

This project is aimed at illustrating how to build a fully functioning admission webhook in the simplest way possible. Most existing examples found on the web rely on heavy machinery using powerful frameworks, yet fail to illustrate how to implement a lightweight webhook that can do much needed actions such as rejecting a pod for compliance reasons, or inject helpful environment variables.

For readability, this project has been stripped of the usual production items such as: observability instrumentation, release scripts, redundant deployment configurations, etc. As such, it is not meant to use as-is in a production environment. This project is, in fact, a simplified fork of a system used across all Kubernetes production environments at Slack.

## Installation
This project runs locally using Docker and Kind. Run the commands below from the repository root.

### Requirements
* Docker with a running Linux container engine
* kubectl 1.37.x
* [Kind 0.33.0](https://kind.sigs.k8s.io/docs/user/quick-start/#installation)
* Bash, OpenSSL, and make
* Go >=1.26.9 (optional; needed for `make test` and `make build`)

The demo uses Kubernetes 1.37.0. Docker builds use Go 1.27.2 and support both amd64 and arm64. Local Go commands also select Go 1.27.2 automatically with the default `GOTOOLCHAIN=auto` setting.

## Usage
### Create Cluster
Create the local cluster and wait for its control plane to become ready:
```sh
make cluster
kubectl --context kind-kind get nodes
kubectl --context kind-kind -n kube-system get pods
```
The node should be `Ready` and system pods should be running. Port 8443 must be available on your host.

Make targets explicitly use the `kind-kind` context. To use another local Kind cluster name, pass `CLUSTER_NAME=example` to each make command; its context will be `kind-example`.

### Deploy Admission Webhook
```sh
make deploy
kubectl --context kind-kind -n default get pods
```
This builds and loads the image, generates fresh TLS certificates, recreates the webhook deployment, waits for HTTPS readiness, and registers the mutating and validating webhooks. The `apps` namespace is labelled to enable admission; the webhook itself runs in `default`.

Certificates and rendered manifests are generated in ignored `dev/generated/`. Run `make deploy` again to rebuild or renew the certificates; each deployment rotates the local CA and updates both webhook configurations.

Stream logs in another terminal (Ctrl-C to stop):
```sh
make logs
```

Check the HTTPS endpoint using the generated CA:
```sh
curl --cacert dev/generated/ca.crt https://localhost:8443/health
```
Expected response: `OK`. The local certificate also covers the service DNS name used by Kubernetes.

### Deploying pods
Create the valid example and wait for it to run:
```sh
make pod
kubectl --context kind-kind -n apps wait --for=condition=Ready pod/lifespan-seven --timeout=120s
kubectl --context kind-kind -n apps exec lifespan-seven -- printenv KUBE
kubectl --context kind-kind -n apps get pod lifespan-seven -o jsonpath='{.spec.tolerations}'
```
The environment variable should be `true`. The tolerations should include `acme.com/lifespan-remaining` values 7 through 14. The webhook logs show both mutation and validation.

Create the invalid example:
```sh
make bad-pod
```
This command is **expected to fail** with `pod name contains "offensive"`. It should report an admission rejection, not a TLS or connection error. Mutation runs before validation, so logs can show both for a rejected pod.

### Cleanup
```sh
make delete-all
make delete-cluster
```
Cleanup works even when the invalid pod was rejected or resources are already absent. It targets the local Kind cluster. Generated files remain in `dev/generated/` for inspection and are replaced on the next deployment.

## Testing
Run the unit and HTTP handler tests:
```sh
make test
```

Build local binaries for macOS and Linux on amd64 and arm64:
```sh
make build
```
The binaries are written to ignored `bin/`.

## Admission Logic
A set of validations and mutations are implemented in an extensible framework. Those happen on the fly when a pod is deployed and no further resources are tracked and updated (ie. no controller logic).

### Validating Webhooks
#### Implemented
- [name validation](pkg/validation/name_validator.go): validates that a pod name doesn't contain any offensive string

#### How to add a new pod validation
To add a new pod validation, create a file `pkg/validation/VALIDATION_NAME.go`, then create a new struct implementing the `validation.podValidator` interface. Register an instance in the `validations` list in `pkg/validation/validation.go`.

### Mutating Webhooks
#### Implemented
- [inject env](pkg/mutation/inject_env.go): inject environment variables into the pod such as `KUBE: true`
- [minimum pod lifespan](pkg/mutation/minimum_lifespan.go): inject a set of tolerations used to match pods to nodes of a certain age, the tolerations injected are controlled via the `acme.com/lifespan-requested` pod label.

#### How to add a new pod mutation
To add a new pod mutation, create a file `pkg/mutation/MUTATION_NAME.go`, then create a new struct implementing the `mutation.podMutator` interface. Register an instance in the `mutations` list in `pkg/mutation/mutation.go`.
