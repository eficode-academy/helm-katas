# Helm and Helmsman

This exercise will demonstrate a declarative approach to Helm chart installation
using [Helmsman](https://github.com/Praqma/helmsman).

## Prerequisites

This exercise assumes that the following repositories (or forks) have been
cloned into a shared folder:

```
https://github.com/eficode-academy/kubernetes-appdev-katas
https://github.com/eficode-academy/k8s-sentences-sentence
https://github.com/eficode-academy/k8s-sentences-age
https://github.com/eficode-academy/k8s-sentences-name
```

I.e. listing folders should produce output as follows:

```shell
$ ls -l
drwxr-xr-x  6 ubuntu ubuntu 4096 Jan  7 11:36 kubernetes-appdev-katas
drwxr-xr-x  7 ubuntu ubuntu 4096 Jan  7 11:36 k8s-sentences-sentence
drwxr-xr-x  7 ubuntu ubuntu 4096 Jan  7 11:36 k8s-sentences-age
drwxr-xr-x  7 ubuntu ubuntu 4096 Jan  7 11:36 k8s-sentences-name
```

Use the following commands to clone Praqma repositories or replace with your
own forks if possible:

```shell
$ cd ~
$ git clone https://github.com/eficode-academy/k8s-sentences-sentence.git
$ git clone https://github.com/eficode-academy/k8s-sentences-age.git
$ git clone https://github.com/eficode-academy/k8s-sentences-name.git
```

## Deploying with Helmsman

Go to the folder which holds our Helmsman specification for the three
microservices that make up the sentences application.

```shell
$ cd kubernetes-appdev-katas/sentences-app/deploy/helmsman/
```

The `helmsman.yaml` file in this folder is an Helmsman [desired state
file](https://github.com/Praqma/helmsman/blob/master/docs/desired_state_specification.md)
which basically is a declarative specification of the Helm charts we want
installed and the parameters we want changed from the defaults.

### Specify Target Namespace

Helmsman can create Kubernetes namespaces. We will not use this features,
however, Helmsman still need a list of namespaces which we will deploy
applications into.

In the `helmsman.yaml` file the application destination namespace is found in
four places. First in the overall list of namespaces (note the trailing ':'):

```
namespaces:
  default:
```

and a namespace destination for each microservce - here an example for the age
service:

```
apps:
  age:
    ...
    namespace: default
```

We want to deploy the sentence application to our own namespace. Investigate the
name of your namespace with:

```shell
$ kubectl config view | grep namespace
```

and then modify the `helmsman.yaml` file to use your namespace (four instances
to update).

### Deploy the Sentences Application

Deploy the sentences microservices (the three Helm charts) to your namespace with the following command:

```shell
$ helmsman -apply -f helmsman.yaml
```

Since Helmsman waits for resources to be deployed and ready, the command may
take a little while to complete.

When Helmsman have deployed the Helm charts, the main service is exposed with a
Kubernetes service of type `LoadBalancer`. Use the following command to get the
IP address of that service and test the deployed application.

```shell
$ kubectl get services
```

### Upgrade a Microservice of the Sentences Application

The `helmsman.yaml` file specifies both the Helm chart version for each
microservice and values that set both container image repository and tag.

Edit the `helmsman.yaml` file and change some of the image
repositories/tags. E.g. change the image repository to your own repository or
use different tags for Praqma images.  The following Praqma images have a
special observable behavior:

`releasepraqma/sentence:1.0-1fee`  (swaps age and name in the generated sentences)

After modifying the `helmsman.yaml`, re-run Helmsman to reconcile the
application specification with the application deployed to Kubernetes:

```shell
$ helmsman -apply -f helmsman.yaml
```

To see, that Helmsman actually installs Helm charts, use the following command
to show installed Helm charts:

```shell
$ helm ls
```

### Optional Exercise

The Sentences application was deployed with the frontend service exposed with a
Kubernetes service of type `LoadBalancer`. Update the `helmsman.yaml` to expose
the front-end service using a service of type `NodePort` instead, upgrade the
deployment and test access to the sentences application using the nodeport
service.


## Cleanup

Delete the application installed with Helmsman:

```shell
$ helmsman -destroy -f helmsman.yaml
```
