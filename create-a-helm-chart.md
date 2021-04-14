# Create an Helm Chart

This exercise will create a simple Helm chart for the sentence application.  The
chart will be 'simple' in the sense that it will not provide support for
customizing application parameters.

In the `sentences-app/deploy/kubernetes/` folder we have Kubernetes YAML
definitions for the three microservices that make up the sentence application
(three Deployments and three Services):

```shell
$ ls -1 sentences-app/deploy/kubernetes/
sentences-age-deployment.yaml
sentences-age-svc.yaml
sentences-deployment.yaml
sentences-name-deployment.yaml
sentences-name-svc.yaml
sentences-svc.yaml
```

In the `kubernetes-appdev-katas`-folder. Use `helm create` to create a simple Helm chart:

```shell
$ mkdir helm-chart
$ cd helm-chart
$ helm create sentence-app
```

Since we will use the sentence application YAML as templates for the chart we
delete the ones created by `helm create`:

```shell
$ rm -rf sentence-app/templates/*
$ echo "" > sentence-app/values.yaml
```

This provides us with skeleton chart without any template files. Next, we copy
the original Kubernetes YAML files to the template folder:

```shell
$ cp -v ../sentences-app/deploy/kubernetes/*.yaml sentence-app/templates/
```

Thats it - now we have a Helm chart for our sentences application.

It is a simple Helm chart in the sense that it has no configurable values, but
it is a complete installable chart and it will use the correct sentence
application Kubernetes YAML definitions.

Before deploying the chart, we run a static validation of it:

```shell
$ helm lint sentence-app/
==> Linting sentence-app/
[INFO] Chart.yaml: icon is recommended

1 chart(s) linted, no failures
```

Normally a chart is fetched from a chart registry (like a container registry),
however, a chart stored locally can also be deployed with Helm. To deploy the
chart from the newly created chart run the following:

```shell
$ helm install sentences sentence-app/
```

Running this command produce the following output:

```
NAME:   sentences
LAST DEPLOYED: Thu Jul 11 14:26:22 2019
NAMESPACE: default
STATUS: deployed
REVISION: 1
TEST SUITE: None
```

To see all the different objects that helm has created, use:

```shell
$ kubectl get pods,services,deployments
NAME                                READY   STATUS    RESTARTS   AGE
pod/sentence-age-6969dc55b6-t6jfk   1/1     Running   0          116s
pod/sentence-name-bb79ff496-pkdsx   1/1     Running   0          116s
pod/sentences-b8c85b468-dr5fv       1/1     Running   0          116s

NAME                    TYPE        CLUSTER-IP      EXTERNAL-IP   PORT(S)          AGE
service/sentence-age    ClusterIP   10.63.255.23    <none>        5000/TCP         117s
service/sentence-name   ClusterIP   10.63.255.250   <none>        5000/TCP         117s
service/sentences       NodePort    10.63.250.194   <none>        5000:31379/TCP   117s

NAME                                  READY   UP-TO-DATE   AVAILABLE   AGE
deployment.extensions/sentence-age    1/1     1            1           117s
deployment.extensions/sentence-name   1/1     1            1           117s
deployment.extensions/sentences       1/1     1            1           117s

```

To see the applications installed with Helm use the `helm ls` operation:

```shell
$ helm ls
NAME         NAMESPACE  REVISION   UPDATED                    STATUS     CHART                APP VERSION
sentences    default    1          Wed Aug 14 08:44:55 2019   DEPLOYED   sentence-app-0.1.0   1.16.0
```

To see the Kubernetes YAML which Helm used to install the application use the `helm get` operation:

```shell
$ helm get all sentences
```

In our case this will be identical to the YAML files we copied previously since
we haven't provided any means of customizing the application installation.

# Food for Thought

In this exercise we created a single Helm chart for the complete application
even though its based on three microservices. When would it make sense to have a
Helm chart for each microservice?

# Cleanup

Delete the application installed with Helm:

```shell
$ helm delete sentences
```
