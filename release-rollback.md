# Helm release and rollback

## Learning Goals

- Learn the commands `install`,`upgrade` and `rollback`.
- See the differences from `install --dry-run` and `helm template`
- Learn how helm stores a release in kubernetes

## Introduction

With the release of Helm3, the way Helm interacts with the kubernetes cluster has changed dramatically!


## Install and upgrade

The `helm upgrade --install` command will create a release if it does not exist in the namespace, or will upgrade a release if a release by that name is found.

<details>
<summary>:bulb: If an explanaition becomes too long, the more detailed parts can be encapsulated in a drop down section</summary>
</details>

## Exercise

### Overview

This exercise assumes you are working in the `release-rollback/` folder.

- Use the different non-deployment functionalities of helm to render kubernetes YAML.
- Deploy the chart some times to create releases
- Look at the history and secrets created

### Step by step instructions

<details>
<summary>More Details</summary>

**Use the different non-deployment functionalities of helm to render kubernetes YAML**

- `helm install --dry-run`
- `helm template`

**Deploy the chart some times to create releases**

- `helm upgrade --install myapp sentence-app/`
- `helm upgrade --install myapp sentence-app/` again to see the different behaviour in output. 
- `helm list` to see the releases
- `helm history myapp`
- `helm get values myapp`
- `helm upgrade --install myapp sentence-app/ --set=replicas=3`
- `helm get values myapp`

**Look at the history and secrets created**

- `kubectl get secrets` to see the releases as secrets files.
- `kubectl describe secrets sh.helm.release.v1.myapp.v2 `
- `kubectl get secret sh.helm.release.v1.myapp.v2 -o jsonpath="{ .data.release }" | base64 -d | base64 -d | gunzip -c`
- `helm rollback`
- `helm history myapp`

</details>

### Clean up

- `helm uninstall myapp`
