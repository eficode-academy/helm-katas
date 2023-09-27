# The Kubernetes package manager

## Learning goal

- Try the Helm cli to spin up a chart

## Introduction

[Enter Helm](https://github.com/helm/helm) - the
answer to how to package multi-container
applications, and how to easily install packages
on Kubernetes.

Helm helps you to:

- Achieve a simple (one command) and repeatable
  deployment
- Manage application dependency, using specific
  versions of other application and services
- Manage multiple deployment configurations: test,
  staging, production and others
- Execute post/pre deployment jobs during
  application deployment
- Update/rollback and test application deployments

## Using Helm charts

Helm uses a packaging format called charts. A
Chart is a collection of files that describe k8s
resources.

Charts can be simple, describing something like a
standalone web server but they can also be more
complex, for example, a chart that represents a
full web application stack included web servers,
databases, proxies, etc.

Instead of installing k8s resources manually via
kubectl, we can use Helm to install pre-defined
Charts faster, with less chance of typos or other
operator errors.

When you install Helm, it does not have a
connection to any default repositories. This is
because Helm wants to decouple the application to
the repository in use.

One of the largest Chart Repositories is the
[BitNami Chart Repository](https://charts.bitnami.com),
which we will be using in these exercises.

Helm chart repositories are very dynamic due to
updates and new additions. To keep Helm's local
list updated with all these changes, we need to
occasionally run the
[repository update](https://helm.sh/docs/helm/helm_repo_update/)
command.

## Exercise

### Overview

- Add a chart repository to your Helm cli
- Install Nginx chart
- Access the Nginx load balanced service
- Look at the status of the deployment with
  `helm ls`
- Clean up the chart deployment

### Step by step

<details>
      <summary>More details</summary>

**Add a chart repository to your Helm cli**

To install the Bitnami Helm Repo and update Helm's
local list of Charts, run:

```shell
helm repo add bitnami https://charts.bitnami.com/bitnami
```
```shell
helm repo update
```

**Install Nginx Chart**

We use the Nginx chart because it is fast and easy to install, and allows us to access the Nginx webserver from our browser to verify that it was deployed.

```shell
helm install my-release bitnami/nginx --set service.type=NodePort
```

This command creates a release called `my-release`
with the bitnami/nginx chart.

The command will output information about your
newly deployed nginx setup similar to this:

```shell
NAME: my-release
LAST DEPLOYED: Wed Sep 27 09:21:48 2023
NAMESPACE: student-3
STATUS: deployed
REVISION: 1
TEST SUITE: None
NOTES:
CHART NAME: nginx
CHART VERSION: 15.3.1
APP VERSION: 1.25.2

** Please be patient while the chart is being deployed **
NGINX can be accessed through the following DNS name from within your cluster:

    my-release-nginx.student-3.svc.cluster.local (port 80)

To access NGINX from outside the cluster, follow the steps below:

1. Get the NGINX URL by running these commands:

    export NODE_PORT=$(kubectl get --namespace student-3 -o jsonpath="{.spec.ports[0].nodePort}" services my-release-nginx)
    export NODE_IP=$(kubectl get nodes --namespace student-3 -o jsonpath="{.items[0].status.addresses[0].address}")
    echo "http://${NODE_IP}:${NODE_PORT}"

```

**Access the Nginx NodePort service**

Get the external port of Nginx with the following commands:

```shell 
kubectl get services
```
- Note down the external port of the Nginx service, it should be in the range of 30000-32767

- Note down one of the external IP addresses of the nodes in the cluster

```shell
kubectl get nodes -o wide
```

result:
  
  ```shell
  NAME                                            STATUS   ROLES    AGE     VERSION   INTERNAL-IP      EXTERNAL-IP     OS-IMAGE             KERNEL-VERSION    CONTAINER-RUNTIME
ip-192-168-83-125.eu-north-1.compute.internal   Ready    <none>   4h29m   v1.26.8   192.168.83.125   13.51.165.230   Ubuntu 20.04.6 LTS   5.15.0-1045-aws   cri-o://1.26.1
ip-192-168-85-161.eu-north-1.compute.internal   Ready    <none>   4h29m   v1.26.8   192.168.85.161   13.53.106.222   Ubuntu 20.04.6 LTS   5.15.0-1045-aws   cri-o://1.26.1
```

- Open a browser and enter the IP address and port of the Nginx service, e.g. `http://<node-ip>:<node-port>`

**Look at the status of the deployment with `helm`
and `kubectl`**

Running `helm ls` will show all current
deployments.

- Run `helm ls` and observe that you have a
  release named `my-release`
- Run `kubectl get pods,deployments,svc` and look
  at a few of the kubernetes objects the release
  created.

> :bulb: As said before Helm deals with the
> concept of
> [charts](https://github.com/kubernetes/charts)
> for its deployment logic. bitnami/nginx was a
> chart,
> [found here](https://github.com/bitnami/charts/tree/master/bitnami/nginx)
> that describes how helm should deploy it. It
> interpolates values into the deployment, which
> for nginx looks
> [like this](https://github.com/bitnami/charts/blob/master/bitnami/nginx/templates/deployment.yaml).
> The charts describe which values can be given
> for overwriting default behavior, and there is
> an active community around it.

**Clean up the chart deployment**

To remove the `my-release` release run:

```shell
helm uninstall my-release
```

</details>
