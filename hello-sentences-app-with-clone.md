# Say Hello to the Sentences Application

This exercise will introduce the example application, aka. the `sentences
application and show how to run it both with docker-compose and on Kubernetes.

## Clone Repositories

This exercise assumes that the following repositories (or forks) have been
cloned into a shared folder.  If you are planning to do the CI/CD exercises, the
three latter repositories should be your own forks.

```
https://github.com/praqma-training/kubernetes-appdev-katas
https://github.com/praqma-training/k8s-sentences-sentence
https://github.com/praqma-training/k8s-sentences-age
https://github.com/praqma-training/k8s-sentences-name
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
$ git clone https://github.com/praqma-training/k8s-sentences-sentence.git
$ git clone https://github.com/praqma-training/k8s-sentences-age.git
$ git clone https://github.com/praqma-training/k8s-sentences-name.git
```

## Running the Sentences Application with Docker Compose

The sentences application consists of three microservices packaged in three
different container images. The sentences application can be run on a Docker
host using e.g. docker-compose.

Use the following command to start containers with the three microservices that
make up the sentences application:

```shell
$ docker-compose -f kubernetes-appdev-katas/sentences-app/deploy/docker-compose.yaml up
```

In another shell, use the following command to request sentences from the
sentences application:

```shell
$ curl 127.0.0.1:8080
Eric is 66 years
```

Stop the docker-compose deployment with Ctrl-C.

## Running the Sentences Application on Kubernetes

Use the following command to deploy the sentences application to Kubernetes:

```shell
$ kubectl apply -f kubernetes-appdev-katas/sentences-app/deploy/kubernetes
```

This will create three microservice deployments with a single POD instance
each. Use the following command to see the status of each POD instance:

```shell
$ kubectl get pods
```

The front-end microservice for the sentences application is exposed with a
Kubernetes service of type `NodePort`.

When all three PODs are in a running state, use the following commands to look
up the actual NodePort used by the frontend microservice:

```shell
$ kubectl get svc sentence
NAME        TYPE       CLUSTER-IP      EXTERNAL-IP   PORT(S)          AGE
sentences   NodePort   10.15.245.208   <none>        8080:30250/TCP   37s
```

In the example above, the relevant NodePort is 30250. Use the following command
to look up an external accessible IP address that can be used to access the
front-end microservice.

```shell
$ kubectl get nodes -o wide
```
Any of the IP addresses from the 'EXTERNAL-IP' column can be used.

To request a sentence from the sentences application, use curl with the external
IP address and NodePort found above:

```shell
$ curl <EXTERNAL-IP>:30250
John is 73 years
```

## Cleanup

```shell
$ kubectl delete -f kubernetes-appdev-katas/sentences-app/deploy/kubernetes/
```