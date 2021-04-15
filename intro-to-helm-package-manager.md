# The Kubernetes package manager

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

## Using helm charts

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

When you install Helm, you are provided with a
default repository of Charts from the
[official Helm Chart Repository](https://artifacthub.io/).

This is a very dynamic list that always changes
due to updates and new additions. To keep Helm's
local list updated with all these changes, we need
to occasionally run the
[repository update](https://docs.helm.sh/helm/#helm-repo-update)
command.

To update Helm's local list of Charts, run:

```
helm repo update
```

Instead of figuring out which docker images to run
manually, we will let helm find them. Let helm
install MySql:

```
helm install stable/mysql
```

This will output information about your newly
deployed mysql setup similar to this:

```
NAME:   invinvible-serval
LAST DEPLOYED: Tue Nov 14 14:46:15 2017
NAMESPACE: default
STATUS: DEPLOYED

RESOURCES:
==> v1beta1/Deployment
NAME                     DESIRED  CURRENT  UP-TO-DATE  AVAILABLE  AGE
invinvible-serval-mysql  1        1        1           0          0s

==> v1/Secret
NAME                     TYPE    DATA  AGE
invinvible-serval-mysql  Opaque  2     0s

==> v1/PersistentVolumeClaim
NAME                     STATUS  VOLUME                                    CAPACITY  ACCESSMODES  STORAGECLASS  AGE
invinvible-serval-mysql  Bound   pvc-2f95ebb1-c942-11e7-9e2e-080027f4e367  8Gi       RWO          standard      0s

==> v1/Service
NAME                     CLUSTER-IP  EXTERNAL-IP  PORT(S)   AGE
invinvible-serval-mysql  10.0.0.25   <none>       3306/TCP  0s
```

Running `helm ls` will show all current
deployments and `helm delete <deployment name>`
(in above example helm delete invinvible-serval)
will remove the service again.

As said before Helm deals with the concept of
[charts](https://github.com/kubernetes/charts) for
its deployment logic. Stable/mysql was a chart,
[found here](https://github.com/kubernetes/charts/tree/master/stable/mysql)
that describes how helm should deploy it. It
interpolates values into the deployment, which for
mysql looks
[like this](https://github.com/kubernetes/charts/blob/master/stable/mysql/templates/deployment.yaml).

The charts describe which values can be given for
overwriting default behaviour, and there is an
active community around it.

Praqma also has
[Helmsman](https://github.com/Praqma/Helmsman)
which is another layer of abstraction on top of
Helm, allowing it to be automated as code
deployments.
