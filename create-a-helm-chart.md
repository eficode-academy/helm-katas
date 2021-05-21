# Create an Helm Chart

## Learning goal

- Create a basic Helm chart for the `sentences`
  application

## Introduction

This exercise will create a simple Helm chart for
the sentence application. The chart will be
'simple' in the sense that it will not provide
support for customizing application parameters.

<details>
      <summary>More information</summary>

In the `sentences-app/deploy/kubernetes/` folder
we have Kubernetes YAML definitions for the three
microservices that make up the sentence
application (three Deployments and three
Services):

```shell
$ ls -1 sentences-app/deploy/kubernetes/
sentences-age-deployment.yaml
sentences-age-svc.yaml
sentences-deployment.yaml
sentences-name-deployment.yaml
sentences-name-svc.yaml
sentences-svc.yaml
```

</details>

## Exercise

### Overview

- Create a skeleton Helm chart
- Copy `sentences` kubernetes yaml into the chart
- Lint and deploy our new chart

### Step by step

<details>
      <summary>More details</summary>

**Create a skeleton Helm chart**

First we create a new directory for our Helm chart, and then use the `helm create` command to create the chart skeleton:

- `mkdir helm-chart`
- `cd helm-chart`
- `helm create sentence-app`

The `helm create` command we just issued created a lot of files that you might want to use when creating a new Helm chart.
We do not need all of those files for the chart we will be creating, therefore we will remove the files we do not need:

- `rm -rf sentence-app/templates/*`
- `echo "" > sentence-app/values.yaml`

This provides us with skeleton chart without any
template files.

**Copy `sentences` kubernetes yaml into the
chart**

Next, we copy the original Kubernetes YAML files
to the template folder:

- `cp -v ../sentences-app/deploy/kubernetes/*.yaml sentence-app/templates/`

That's it - now we have a Helm chart for our
sentences application.

> :bulb: It is a simple Helm chart in the sense
> that it has no configurable values, but it is a
> complete installable chart and it will use the
> correct sentence application Kubernetes YAML
> definitions.

**Lint and deploy our new chart**

Before deploying the chart, we run a static
validation of it:

- `helm lint sentence-app/`

Running this command produces the following output:

```shell
==> Linting sentence-app/
[INFO] Chart.yaml: icon is recommended

1 chart(s) linted, 0 chart(s) failed
```

> :bulb: Normally a chart is fetched from a chart
> registry (like a container registry), however, a
> chart stored locally can also be deployed with
> Helm.

To deploy the chart from the newly created chart
run the following:

- `helm install sentences sentence-app/`

Running this command produces the following output:

```shell
NAME: sentences
LAST DEPLOYED: Wed Apr 21 10:43:55 2021
NAMESPACE: user1
STATUS: deployed
REVISION: 1
TEST SUITE: None
```

To see all the different objects that Helm has
created, use:

```shell
$ kubectl get pods,services,deployments
NAME                                READY   STATUS    RESTARTS   AGE
pod/sentence-age-78fc854dd5-w9gdq   1/1     Running   0          64s
pod/sentence-name-ff4c584b9-txp5n   1/1     Running   0          64s
pod/sentences-746cc46db8-khp85      1/1     Running   0          64s

NAME               TYPE        CLUSTER-IP       EXTERNAL-IP   PORT(S)          AGE
service/age        ClusterIP   10.191.240.60    <none>        8080/TCP         66s
service/name       ClusterIP   10.191.251.238   <none>        8080/TCP         66s
service/sentence   NodePort    10.191.245.72    <none>        8080:32665/TCP   66s

NAME                            READY   UP-TO-DATE   AVAILABLE   AGE
deployment.apps/sentence-age    1/1     1            1           66s
deployment.apps/sentence-name   1/1     1            1           66s
deployment.apps/sentences       1/1     1            1           66s
```

To see the applications installed with Helm use
the `helm ls` operation:

- `helm ls`

```shell
NAME            NAMESPACE       REVISION        UPDATED                                 STATUS         CHART                    APP VERSION
sentences       user1           1               2021-04-21 10:43:55.789048706 +0000 UTC deployed       sentence-app-0.1.0       1.16.0
```

To see the Kubernetes YAML which Helm used to
install the application use the `helm get`
operation:

- `helm get all sentences`

In our case this will be identical to the YAML
files we copied previously since we haven't
provided any means of customizing the application
installation.

</details>

## Food for Thought

In this exercise we created a single Helm chart
for the complete application even though it is based
on three microservices. When would it make sense
to have a Helm chart for each microservice?

## Cleanup

Uninstall the application release with Helm:

```shell
$ helm uninstall sentences
```
