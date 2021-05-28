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

- In bullets, what are you going to solve as a student

### Step by step instructions

- `helm install --dry-run`
- `helm template`
- `helm upgrade --install myapp sentences-app/`
- `helm upgrade --install myapp sentences-app/` again to see the different behaviour in output. 
- `helm list` to see the releases
- `kubectl get secrets` to see the releases as secrets files.
- `kubectl describe secrets`
- `helm rollback`

<details>
<summary>More Details</summary>

**take the same bullet names as above and put them in to illustrate how far the student have gone**

- all actions that you believe the student should do, should be in a bullet

> :bulb: Help can be illustrated with bulbs in order to make it easy to distinguish.

</details>

### Clean up

If anything needs cleaning up, here is the section to do just t

helm install mysite bitnami/drupal --set drupalUsername=admin
helm list
