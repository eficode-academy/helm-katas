# External charts and chart dependencies

## Learning goal

- Use external charts in your deployment
- Configure the chart dependencies in the external chart

## Introduction

### Helm Chart Dependencies

## Exercise

Add Bitnami PostgreSQL to your todo-mvc application

### Overview

- Install bitnami repo
- Pull down postgreSQL for inspection
- Add postgreSQL to chart dependency
- Configure backend to use the right db service endpoint
- Add loadbalancer type for frontend

### Step-by-step

TODO uncomment
<!-- <details> -->
<!-- <summary>More information</summary> -->

**Install Bitnami Helm repo**

Install the Bitnami Helm Repo and update Helm's local list of Charts:

```
helm repo add bitnami https://charts.bitnami.com/bitnami
helm repo update
```

**Pull down the PostgreSQL chart**

- Pull down the chart from bitnami: `helm pull --untar bitnami/postgresql`

- Open `external-charts/postgresql/values.yaml` to see all the possible values that are configurable.

**Inspect the `chart.yaml` and the `charts/` folder**

**Update the values file**

- Set your own username and password in our pre-made values file in `external-charts/values.yaml`

**Install the chart**


**Install the new one**

TODO uncomment
<!-- </details> -->

### Clean up

- `helm uninstall my-todo`
- `kubectl delete pvc `

#### Resources

https://helm.sh/docs/chart_best_practices/dependencies/#versions
https://helm.sh/docs/helm/helm_pull/
