# Maintainers Guide

This document describes tools, tasks and workflow that one needs to be familiar with in order to effectively maintain
this project. If you use this package within your own software as is but don't plan on modifying it, this guide is
**not** for you.

## Tools (optional)

Tools, dependencies, or other programs someone maintaining this project needs to be familiar with:
* Kubernetes
* Go
* Docker
* Kind

## Tasks

### Testing

Unit can be run like so:
```
❯ go test ./...
?   	github.com/slackhq/simple-kubernetes-webhook	[no test files]
ok  	github.com/slackhq/simple-kubernetes-webhook/pkg/admission	0.743s
ok  	github.com/slackhq/simple-kubernetes-webhook/pkg/mutation	1.065s
ok  	github.com/slackhq/simple-kubernetes-webhook/pkg/validation	0.413s
```

### TLS certificate
Kubernetes only allows admission webhooks running with `https`. To generate a TLS secret, run [`./dev/gen-certs.sh`](/dev/gen-certs.sh). The base64 caBundle needs to be manually copied and pasted in the `MutatingWebhookConfiguration` and `ValidattingWebhookConfiguration` at [`./dev/manifests/cluster-config/`](./dev/manifests/cluster-config/)

### Logs
The logs level defaults to `debug` and can be set with the env var:
```
LOG_LEVEL=info
```
The logs format defaults to `text` and can be set to `json` with the env var:
```
LOG_JSON=true
```

### Releasing

N/A: this demo project is not released

## Workflow

### Versioning and Tags

N/A: this demo project is not released

### Branches

The `main` branch is where active development occurs, feel free to name your feature / bug fix branch what your heart desires.

### Issue Management

Labels are used to run issues through an organized workflow. Here are the basic definitions:

*  `bug`: A confirmed bug report. A bug is considered confirmed when reproduction steps have been
   documented and the issue has been reproduced.
*  `enhancement`: A feature request for something this package might not already do.
*  `docs`: An issue that is purely about documentation work.
*  `tests`: An issue that is purely about testing work.
*  `needs feedback`: An issue that may have claimed to be a bug but was not reproducible, or was otherwise missing some information.
*  `discussion`: An issue that is purely meant to hold a discussion. Typically the maintainers are looking for feedback in this issues.
*  `question`: An issue that is like a support request because the user's usage was not correct.
*  `semver:major|minor|patch`: Metadata about how resolving this issue would affect the version number.
*  `security`: An issue that has special consideration for security reasons.
*  `good first contribution`: An issue that has a well-defined relatively-small scope, with clear expectations. It helps when the testing approach is also known.
*  `duplicate`: An issue that is functionally the same as another issue. Apply this only if you've linked the other issue by number.

> You may want to add more labels for subsystems of your project, depending on how complex it is.

**Triage** is the process of taking new issues that aren't yet "seen" and marking them with a basic
level of information with labels. An issue should have **one** of the following labels applied:
`bug`, `enhancement`, `question`, `needs feedback`, `docs`, `tests`, or `discussion`.

Issues are closed when a resolution has been reached. If for any reason a closed issue seems
relevant once again, reopening is great and better than creating a duplicate issue.

## Everything else

When in doubt, find the other maintainers and ask.
