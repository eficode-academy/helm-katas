# Say Hello to the Sentences Application

## Learning goal

- Try out the sentences application

## Introduction

This exercise introduces the example application,
aka. the `sentences application` and show how to
run it both with docker-compose and on Kubernetes.

## Running the Sentences Application with Docker Compose

<details>
      <summary>More information</summary>

The sentences application consists of three
microservices packaged in three different
container images. The sentences application can be
run on a Docker host using e.g. docker-compose.

</details>

## Exercise

Start docker-compose with the three microservices
that make up the sentences application.

### Step by step

<details>
      <summary>More details</summary>

Open a terminal in the root of the git repository (helm-katas) and use `docker-compose up` to deploy the stack:

- `docker-compose -f sentences-app/deploy/docker-compose.yaml up -d`

Use the following command to request sentences
from the sentences application:

- `curl 127.0.0.1:8080`
- Observe that you get something like
  `Eric is 66 years` back

</details>

## Running the Sentences Application on Kubernetes

Now we would like to deploy the application into
our cluster.

## Exercise

### Overview

- Deploy the sentences application to Kubernetes
- Test the deployed application

### Step by step

<details>
      <summary>More details</summary>

**Deploy the sentences application to Kubernetes**

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

- `curl <EXTERNAL-IP>:30250`

Output:

```shell
John is 73 years
```

> :bulb: in the above example `30250` should be
> changed with your nodeport found above

</details>

## Cleanup

- `docker-compose -f sentences-app/deploy/docker-compose.yaml down`
- `kubectl delete -f sentences-app/deploy/kubernetes/`
