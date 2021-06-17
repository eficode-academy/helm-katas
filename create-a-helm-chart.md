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

- Running the sentences application on Kubernetes
- Create a skeleton Helm chart
- Copy `sentences` kubernetes yaml into the chart
- Lint and deploy our new chart

### Step by step

<details>
      <summary>More details</summary>

**Deploy the sentences application to Kubernetes**

First, let's run the application in bare Kubernetes to see that our YAML is right.

- `kubectl apply -f sentences-app/deploy/kubernetes`

This will create three microservice deployments
with a single POD instance each.

**Test the deployed application**

- `kubectl get pods`

> :bulb: The front-end microservice for the
> sentences application is exposed with a
> Kubernetes service of type `NodePort`.

When all three PODs are in a running state, look
up the actual NodePort used by the frontend
microservice:

- `kubectl get svc sentence`

Output:

```shell
NAME        TYPE       CLUSTER-IP      EXTERNAL-IP   PORT(S)          AGE
sentences   NodePort   10.15.245.208   <none>        8080:30250/TCP   37s
```

In the example above, the relevant NodePort is
`30250`.

- look up an external accessible IP address that
  can be used to access the front-end
  microservice.

- `kubectl get nodes -o wide`

Any of the IP addresses from the
`EXTERNAL-IP`-column can be used.

To request a sentence from the sentences
application, use curl with the external IP address
and `NodePort` found above:

- `curl <EXTERNAL-IP>:<NodePort>`

Output:

```shell
John is 73 years
```

> :bulb: in the above example `NodePort` should be
> changed with your nodeport found above

- Clean up the application deployed with `kubectl delete -f sentences-app/deploy/kubernetes/`

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

Try to reach it again like we did with the raw kubernetes objects application to begin with.

- Note down the NodePort from the service `kubectl get svc`

Look up an external accessible IP address that can be used to access the front-end microservice.

- `kubectl get nodes -o wide`

Any of the IP addresses from the `EXTERNAL-IP`-column can be used.

- `curl <EXTERNAL-IP>:<NodePort>` and see that your application is running once again.

Output:

```shell
John is 47 years
```

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
