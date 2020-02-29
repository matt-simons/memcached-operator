# Memcached Go Operator

## Overview

This Memcached operator is a simple example operator for the [Operator SDK][operator_sdk] and includes some basic end-to-end tests.

## Prerequisites

- [go][go_tool] version v1.13+.
- [docker][docker_tool] version 17.03+
- [kubectl][kubectl_tool] v1.14.1+
- [operator-sdk][operator_install]
- Access to a Kubernetes v1.14.1+ cluster

## Getting Started

1. Create a namespace `kubectl create namespace my-name`
2. Edit Makefile and replace `default` with your namespace name.
3. Run the unit tests `make test-unit`. This will run using a mocked Kubernetes API.
4. Run the system tests `make test-e2e`. This will run using your real Kubernetes API.
5. Add the necessary functionality to `pkg/controllers/memcached/memcached_controller.go` to pass the tests.

### Pulling the dependencies

Run the following command

```
$ go mod tidy
```

### Running the operator locally

Instead of building and redeploying the operator on every code change you can run it locally by running:
```
$ make run
```
This enables quicker development iterations.

### Building the operator

Build the Memcached operator image and push it to a registry:

```
$ oc project <my-namespace>
$ make imagestream-init
$ make build
$ kubectl apply -f ./deploy
```


[dep_tool]: https://golang.github.io/dep/docs/installation.html
[go_tool]: https://golang.org/dl/
[kubectl_tool]: https://kubernetes.io/docs/tasks/tools/install-kubectl/
[docker_tool]: https://docs.docker.com/install/
[operator_sdk]: https://github.com/operator-framework/operator-sdk
[operator_install]: https://github.com/operator-framework/operator-sdk/blob/master/doc/user/install-operator-sdk.md
