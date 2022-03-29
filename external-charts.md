# External charts and chart dependencies

## Learning Goals

- Use external charts as dependencies in your own chart
- Configure values of chart dependencies

## Introduction

Helm lets us write all of the custom templates we need, but it also allows us to include the work of others as dependencies for our own charts.

Just like when we are working with containers, we would pull down the official `mongodb` image whenever we need a mongodb instance, instead of building our own.
Similarly with helm, we can use off-the-shelf charts for things like databases as part of our own charts.

This is done by declaring external charts as dependencies in our `Chart.yaml` or by pulling them down manually and placing them in the `charts/` directory of our own chart.

We will be using the [Bitnami mysql](https://artifacthub.io/packages/helm/bitnami/mysql) helm chart in the exercise.

We want to use an external chart to install the database, because it means we do not have to worry about the implementation of the mysql templates.
We can instead focus on developing our own application, and use off-the-shelf charts for any external dependencies we might have.

You can read more about chart depenencies in the [documentation](https://helm.sh/docs/topics/charts/#chart-dependencies).

## Exercise

In this exercise we will be using an example todo application (originally from docker [getting started](https://github.com/eficode-academy/getting-started)) to show how you can use external charts as dependencies for your own Helm charts.

The example application is simple todo list.
It consists of a node application that serves the frontend and stores the todo items in a mysql database.
The example chart we have created thus contains two deployments and two services:

The two deployments are one for the `todo application` and one for a `mysql database`.

As well as two services to handle networking, a `ClusterIP` for the mysql database and a `NodePort` for the todo application.

We will start by deploying the example todo application with all custom Kubernetes yaml, and then modify the chart to use the `bitnami/mysql` chart as a dependency, instead of our custom mysql deployment.

If you get stuck, or you want to see how the finished chart looks, a completed version can be found in `helm-katas/external-charts/done`.

### Overview

- Deploy the initial version of the todo application
- Install Bitnami chart repository
- Pull the bitnami/mysql chart and inspect its contents
- Remove the custom mysql templates
- Declare the bitnami/mysql chart as a dependency
- Pull the Bitnami chart as a dependency
- Configure the mysql chart
- Upgrade your installed release
- Verify that the application is working correctly with the new database

### Step-by-step

<details>
<summary>Steps</summary>

**Deploy the initial version of the todo application**

- Let's start by navigating to the exercise `cd helm-katas/external-charts/start`

- Thereafter deploying the basic version of the todo application chart.

- Deploy the application:
```sh
$ helm install my-todo todo
NAME: my-todo
LAST DEPLOYED: Sun May 30 16:03:34 2021
NAMESPACE: user1
STATUS: deployed
REVISION: 1
TEST SUITE: None
```

- Get the node port of the `todo-app` service, with `$ kubectl get svc`

- Access the service from your browser using the node port.

<details>
<summary>:bulb: how do I get the nodeport?</summary>

You need the external IP of the one of the cluster nodes.
You can get that by issuing `kubectl get nodes -o wide`

</details>

- Play around with the application for a moment, add some items, delete some items, and refresh the page to verify that the state of the todo application is persisted to the mysql database.

Once you are confident that the application works, we then proceed to using a chart dependency instead of our custom deployment to install the mysql database.

**Install Bitnami chart repository**

In order to install a chart from the Bitnami chart repository, we have to `install` it in our local Helm client.

This is done with `repo add` and `repo update` commands.

- Use the following commands to install the Bitnami repository:

```sh
$ helm repo add bitnami https://charts.bitnami.com/bitnami
$ helm repo update
```

- Verify that the repository was installed:

```sh
$ helm repo list
NAME   	URL
bitnami	https://charts.bitnami.com/bitnami
```

**Pull the bitnami/mysql chart and inspect its contents**

Let's try to pull down the `bitnami/mysql` chart so that we can inspect it.

The chart can be found either on [artifacthub.io](https://artifacthub.io/packages/helm/bitnami/mysql) or on [Bitnami's own website](https://bitnami.com/stack/mysql), where we can read about the chart, and all of the different values that we might want to customize.

- Pull the repository down:

```sh
$ helm pull bitnami/mysql --untar
```

> :bulb: Helm charts are stored as compressed tar archives, therefore when you pull a chart, you will get the tar file.
> We can use the `--untar` option to automatically unpack the archive such that we can inspect its contents.

You should now have a directory named `mysql`, which contains the helm chart.

- Take a moment to inspect the different files in the chart, especially the `values.yaml` and how the values are propagated to the different `templates`.

- Once you are satisfied that the helm chart looks good, go ahead and delete the `mysql` directory that was created.

```sh
$ rm -rf mysql
```

> Note: Pulling down charts and unpacking them is useful for inspecting them before use.
> Charts pulled down and placed in the `charts/` directory are used as subcharts.
> If you want to do a lot of customization to a subchart, unpacking it can be the easiest way to do so, as you will have all of the files available.

In this case we do not want to customize any of the templates of the mysql chart, just pass some custom values, in order to make our own chart simpler.

**Remove the custom mysql templates**

Before we add the bitnami chart as a dependency, let's remove the existing custom templates for mysql:

- Delete the files:
    - `rm todo/templates/mysql-deployment.yaml`
    - `rm todo/templates/mysql-service.yaml`

**Declare the bitnami/mysql chart as a dependency**

We can declare external charts as dependencies for our own chart by adding an entry to the `dependencies` map in our `Chart.yaml`.

- Open your `Chart.yaml` - `todo/Chart.yaml` and add the following:

```yaml
...
dependencies:
  - repository: https://charts.bitnami.com/bitnami
    name: mysql
    version: 8.5.10
```

> :bulb: the `dependencies` key takes a list of chart dependencies.
> A chart can have an arbitrary number of dependencies.

Which declares that our todo chart is dependent upon the bitnami chart `mysql` of the specified version.

**Pull the Bitnami chart as a dependency**

Now that we have declared the bitnami mysql chart as a dependency, we can pull the chart down so that we can install it as part of our release:

- Use the `helm dependency update` command to install the `bitnami/mysql` chart locally.

```sh
$ helm dependency update todo
Hang tight while we grab the latest from your chart repositories...
...Successfully got an update from the "bitnami" chart repository
Update Complete. ⎈Happy Helming!⎈
Saving 1 charts
Downloading mysql from repo https://charts.bitnami.com/bitnami
Deleting outdated charts
```

This will place the `bitnami/mysql` chart as tar file in the in the `charts/` directory.

```sh
$ ls todo/charts
mysql-8.5.10.tgz
```

**Configure the mysql chart**

The bitnami mysql chart comes with a very long list of default values, as you saw when we inspected the chart earlier.

In this case we only care about configuring a few different values, which are the name of the created mysql database as well as the credentials for connecting to it.

Luckily the `bitnami/mysql` chart has values we can set for these.

We do this by adding the values to our own `values.yaml` and place them under a top key that has the same name as the chart, in this case `mysql`.

- Edit your `todo/values.yaml` and change:

```yaml
...
mysql:
  dbName: todos
  dbPassword: todos
  dbUser: todos
  dbRootPassword: todos
  service:
    name: todo-mysql
    type: ClusterIP
```

To:

```yaml
...
mysql:
  image:
    tag: 5.7
  auth:
    rootPassword: todos
    database: todos
    username: todos
    password: todos
  primary:
    resources:
      requests:
        cpu: 0.25
        memory: "250Mi"
      limits:
        cpu: 1.0
        memory: "1000Mi"
```

First we customize the container image tag to use, with the `mysql.image.tag` key.
We know that our todo application uses mysql version `5.7`, so we choose the container image tag with the appropriate version.

Next, by defining the `mysql.auth` values, we can have the chart automatically setup a database with the provided name, and create a user with the specified username and password.

Finally we set some resource requests and limits.

> :bulb: we have to specify the resource limits, as the mysql chart requires more resources than our cluster by default allows per pod.

We must also change the value that specifies the host of the mysql database in the todo application.

The service created from the bitnami chart will get the name `<release-name>-mysql`, so if we call our release `my-todo` then the service name will be `my-todo-mysql`.

- We change the value in our values.yaml:

From:

```yaml
todoApp:
  mysqlHost: todo-mysql
  ...
```

To:

```yaml
todoApp:
  ## the hostname of the service will be `<release name>-mysql`
  mysqlHost: my-todo-mysql
  ...
```

**Upgrade your installed release**

- See the differences that the change will make in kubernetes YAML with the helm diff plugin: `helm diff upgrade my-todo todo/`

- Upgrade your installation with the changes we have made:

```sh
$ helm upgrade my-todo todo
Release "my-todo" has been upgraded. Happy Helming!
NAME: my-todo
LAST DEPLOYED: Sun May 30 17:21:00 2021
NAMESPACE: user1
STATUS: deployed
REVISION: 2
TEST SUITE: None
```

Now we wait a moment for the new resources to be deployed to the cluster.

> :bulb: you can use a command like `watch kubectl get all` to continually monitor the rollout.

**Verify that the application is working correctly with the new database**

Now go back to your browser and navigate to the todo application endpoint.

The application should still be running, but all of the todo entries should be gone, as we have now connected to the new database.

> Note: the reason why all data is gone is because the initial version of MySQL we made, did not have any persistent volume. The Bitnami Chart however uses that, so new database application upgrades is doable while still keeping the application data.

Play around with the application again, and verify that the application is working correctly with the new database.

</details>

### Clean up

- `helm uninstall my-todo`
- `kubectl delete pvc data-my-todo-mysql-0` (your pvc name might be different)

### Resources

https://helm.sh/docs/chart_best_practices/dependencies/#versions
https://helm.sh/docs/helm/helm_pull/
