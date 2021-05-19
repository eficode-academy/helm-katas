# Advanced usage of other charts

Install wordpress on your cluster

- install bitnami repo (if not already done from first)
- helm pull --untar bitnami/wordpress
- inspect the `chart.yaml` and the `charts/` folder
- install with a `values.yaml` file `helm install my-wordpress -f values.yaml wordpress`
- change the dependency version of mariaDB
- helm dependency update
- install the new one

https://helm.sh/docs/chart_best_practices/dependencies/#versions
