# External charts and chart dependencies

## Learning goal

- Use external charts in your deployment
- Configure the chart dependencies in the external chart

## Introduction

## Exercise

Install the Bitnami wordpress chart on your cluster, by first pulling it down to your machine, alternating it and use it to apply the chart on the cluster.

### Overview

- Install bitnami repo
- `helm pull --untar bitnami/wordpress`
- Inspect the `chart.yaml` and the `charts/` folder
- Install with a `values.yaml` file `helm install my-wordpress -f values.yaml wordpress`
- Change the dependency version of mariaDB
- `helm dependency update`
- Install the new version

### Step-by-step

<details>
      <summary>More information</summary>

**Install Bitnami Helm repo**

Install the Bitnami Helm Repo and update Helm's local list of Charts:

```
helm repo add bitnami https://charts.bitnami.com/bitnami
helm repo update
```

**Pull down the wordpress chart**

We are going to have a look at the wordpress chart before applying it. 

- Navigate to the `external-charts/start` folder in the root of this repository in your terminal.
- Pull down the chart from bitnami: `helm pull --untar bitnami/wordpress`

Your folder should now look something like this:

```sh
.
├── values.yaml
└── wordpress
    ├── Chart.lock
    ├── Chart.yaml
    ├── README.md
    ├── charts
    ├── ci
    ├── templates
    ├── values.schema.json
    └── values.yaml
```
- Open `external-charts/wordpress/values.yaml` to see all the possible values that are configurable.

**Inspect the `chart.yaml` and the `charts/` folder**

- Look at the `external-charts/wordpress/Chart.yaml` file to see the three dependencies that wordpress depends on; MariaDB, Memcached, and Common.
- Look in the `external-charts/wordpress/charts` folder to see the three dependencies also getting pulled down, each in their own folder.

**Update the values file**

- Set your own username and password in our pre-made values file in `external-charts/values.yaml`

**Install the chart**


- Install the chart in your cluster: `helm install my-wordpress wordpress -f values.yaml`
- Inspect that all pods comes online: `kubectl get pods,deployments`
- Try to access the wordpress site with the new external loadbalancer ip: `kubectl get svc`

> Note: You might be redirected to a HTTPS site with an untrusted certificate. It's perfectly fine in this exercise. It is because we do not use TLS certificates signed by a trusted certificate authority in the training environment. You should always use properly trusted certificates in production.

- Try to log into wordpress backend by accessing the admin site on: `https://<LoadBalancerIP>/admin`

**Change the dependency version of memcached**

When pulling a chart down with dependencies, the dependency charts are getting pulled down as well.
We will try alternating one of the dependencies before deploying again.

- Find the avaliable versions for memcached with `helm search repo memcached -l`

> :bulb: remember that a chart has two versions: Chart version called `version` and application version `appVersion`

- Try changing the version

**Helm dependency update**



**Install the new one**

</details>

### Clean up

- `helm uninstall my-wordpress`
- `kubectl delete pvc data-my-wordpress-mariadb-0`

#### Resources

https://helm.sh/docs/chart_best_practices/dependencies/#versions
https://github.com/bitnami/charts/tree/master/bitnami/wordpress/#installing-the-chart
https://bitnami.com/stack/wordpress/helm
https://helm.sh/docs/helm/helm_pull/