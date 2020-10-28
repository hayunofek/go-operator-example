# Mongo & MongoWeb Operator in go
OpenShift basic operator written in Go, which deploys a deployment, a service and a route.

# What does this operator do?
OpenShift basic operator written in Go, which deploys a deployment with [Size] replicas of a specified [Image], a service and a route. 
This operator has been tested on openshift v4.4.

# Requirements
- operator-sdk: https://sdk.operatorframework.io/docs/installation/install-operator-sdk/
- mercurial version 3.9+
- bazaar version 2.7.0+
- go version v1.15+.
- docker v17.03+ (or another tool compatible with multi-stage Dockerfiles).
- kubectl version v1.11.3+ (v1.16.0+ if using apiextensions.k8s.io/v1 CRDs).

# How was the project initialized
- `mkdir operator_name`
- `cd operator_name`
- `operator-sdk init --domain=ocp4.example --repo=gitlab.com/openshift4ee/go-operator-example` 
- `operator-sdk create api --group ocp4ee --version v1alpha1 --kind GoMongo`

# Defining your CRD
Edit the `api/v1alpha1/CRDNAME_types.go` file, and add the desired fields in _CRDNAME_ Spec. Afterwards, execute `make generate`.

# Updating your manifests
In order to update the project manifests, for example after changing your rbac permissions in gomongo_controller.go, execute `make manifests`.
(For more information consult the following tutorial - https://sdk.operatorframework.io/docs/building-operators/golang/tutorial/ )

# How to build and push the operator image
- `make docker-build IMG=<some-registry>/<project-name>:<tag>`
- `make docker-push IMG=<some-registry>/<project-name>:<tag>`

# How to deploy the operator (on a terminal with oc installed and logged in)
- `cd config/default/ && kustomize edit set namespace "namespace_to_deploy" && cd ../..`
- `make install`
- `make deploy IMG=<some-registry>/<project-name>:<tag>`

The first command sets the namespace in which the operator shall be deployed. The second registers the CRD in the cluster, and the third deploys the actual operator.

# How to deploy sample operator custom resource
- `oc apply -f config/samples/ocp4ee_v1_gomongo.yaml`

