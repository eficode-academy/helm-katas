# Helm Release and Rollback

## Learning Goals

- Learn the helm commands `install`,`upgrade` and `rollback`.
- See the differences from `install --dry-run` and `helm template`
- Learn how helm stores a release in Kubernetes

## Introduction

With the release of Helm3, the way Helm interacts with the Kubernetes cluster has changed dramatically!
Helm speaks directly with the cluster, and stores any state of the releases as Kubernetes native objects.

### Producing Kubernetes Output

Helm charts consists of Kubernetes YAML with [golang template](https://helm.sh/docs/chart_template_guide/functions_and_pipelines/) functionality.
The more complex your templating gets, the harder it can be to keep track of how the rendered YAML will look.
To help you check that the YAML will render correctly Helm has two different commands:

- `helm <install/upgrade> --dry-run` will produce the all the output (including the helm status), but redirect all Kubernetes YAML to standard-out instead of sending it to the cluster.
- `helm template` will only produce the Kubernetes YAML, in the version that the helm client was compiled for.

### Install and Upgrade

Helm has two separate commands for installing a new release, or upgrading existing release.
The reason for this division is that in production, you would like to be sure that you only install a new service if it is not installed, and only update a service that is installed.

If you do not want those guards, then you can simply use the `helm upgrade` command with the `--install` option.
The `helm upgrade --install` command will create a release if it does not exist in the namespace, or upgrade a release if a release by that name is found.

>:bulb: What will this upgrade change?
> There is a helm plugin called [helm diff](https://artifacthub.io/packages/helm-plugin/diff/diff).
> The helm diff plugin gives your a preview of what a helm upgrade would change.
> It basically generates a diff between the latest deployed version of a release and a `helm upgrade --debug --dry-run`. This can also be used to compare two revisions/versions of your helm release.

### Rollbacks

Each helm release will store a secret inside the namespace the given release has been deployed to, which contains data about the release.
This enables, helm to issue a rollback if the release is faulty in any way.

Now, if something does not go as planned during a release, it is easy to roll back to a previous release using `helm rollback [RELEASE] [REVISION]`.

```sh
$ helm rollback myapp 1
Rollback was a success! Happy Helming!
```

The above rolls back myapp release to its very first release version.
A release version is an incremental revision. Every time an install, upgrade, or rollback happens, the revision number is incremented by 1. The first revision number is always 1.

We can use `helm history [RELEASE]` to see revision numbers for a certain release.

## Exercise

### Overview

- Use the different non-deployment functionalities of helm to render kubernetes YAML.
- Deploy the chart multiple times to create a revision history
- Look at the history and secrets created

### Step by step instructions

<details>
<summary>More Details</summary>

**Use the different non-deployment functionalities of helm to render kubernetes YAML**

- Navigate to the folder with the exercise content `cd release-rollback`

Let us start with seeing the difference between `install --dry-run` and `template`.

- Run install and observe the three parts of the output: `helm install --dry-run myapp sentence-app/`

> :bulb: Can you find out where the NOTES: section of the output is generated from?

- In order only to get the rendered kubernetes YAML, helm template will be a better fit: `helm template myapp sentence-app/`

**Deploy the chart multiple times to create a revision history**

- Run `helm upgrade --install myapp sentence-app/` _twice_, and observe the different behaviour in output.
- Run `helm list` to see the release, and the revision state.
- Run `helm history myapp` and look at the two revisions made. Output should look like below.

```sh
$ helm history myapp
REVISION        UPDATED                         STATUS          CHART                   APP VERSION     DESCRIPTION
1               Fri May 28 11:38:49 2021        superseded      sentence-app-0.1.0      1.16.0          Install complete
2               Fri May 28 11:39:13 2021        deployed        sentence-app-0.1.0      1.16.0          Upgrade complete
```

- Run `helm get values myapp` and observe the output.

> :bulb: why are there no user-supplied values, even though we have a values.yaml file outside the chart?

<details>
<summary>Solution</summary>

The reason is that none of our two upgrade commands took the values.yaml file as a parameter.

</details>

> :bulb: If do not have the `diff` plugin installed you can install it with: `$ helm plugin install https://github.com/databus23/helm-diff`

- Run `helm diff upgrade myapp sentence-app -f values.yaml` to see what changes the values in `values.yaml` would have to our release.
- Apply the values, run: `helm upgrade --install myapp sentence-app/ -f values.yaml` to create a new release with the values applied as well.
- Rerun `helm get values myapp` and observe the changed output.

**Look at the history and secrets created**

Let us have a look deeper inside how helm stores the release data.

- `kubectl get secrets` to see the releases as secrets files.

> Note: the release secrets are a specific helm type `helm.sh/release.v1` of secret. Different types vary in terms of the validations performed and the constraints Kubernetes imposes on them.

- Look closer at one of the secrets with: `kubectl describe secrets sh.helm.release.v1.myapp.v2`

```sh
$ kubectl describe secrets sh.helm.release.v1.myapp.v2
Name:         sh.helm.release.v1.myapp.v2
Namespace:    user1
Labels:       modifiedAt=1622202110
              name=myapp
              owner=helm
              status=superseded
              version=2
Annotations:  <none>

Type:  helm.sh/release.v1

Data
====
release:  4444 bytes
```

This gives us some overall information about the release, but does not really tell us where the release data is stored.

- `kubectl get secret sh.helm.release.v1.myapp.v2 -o json` will give us the entire secret, including the `release` field.

> :bulb: the release field is not human readable as it is now.
> So in order to decode the data, you have to:
>
>- base64 decode - Kubernetes secrets encoding
>- base64 decode (again) - Helm encoding
>- gzip decompress - Helm zipping

- Run `kubectl get secret sh.helm.release.v1.myapp.v2 -o jsonpath="{ .data.release }" | base64 -d | base64 -d | gunzip -c` to see the all kubernetes object that was deployed with this helm release.

Now that we have seen how helm stores the data for a release, let us try to make a rollback to the initial release.

- Run `helm rollback myapp 1` to issue a rollback.

- Run `helm history myapp` to see the list of releases.

> :bulb: why does a rollback create a new release?

</details>

## Extra: 3-way merges and hand-edited releases

When doing rollbacks and upgrades, Helm3 tries to perform a 3 way _strategic_ merge.

The three states it takes are:

- The latest revision state of the release stored in the Helm secret.
- The current "live" versions of the Kubernetes objects running in the cluster.
- The new revision either being upgraded to or rolled back.

If the live state has added something to the objects, Helm tries to add them to the new revision as well. If something has changed, say an environment variable, it will be overwritten to the revision currently getting applied.

<details>
<summary>:bulb: More information on merge strategies.</summary>

The strategic-merge approach attempts to “do the right thing” when combining the provided spec with the existing spec. More specifically, it attempts to merge both objects and arrays, meaning changes tend to be additive. For example, providing a patch that contains a single, new environment variable in a pod container spec results in that environment variable being added to the existing environment variables, not overwriting them. To delete a property with this approach, you need to specifically set its value to null in the provided spec.

With a strategic merge, a list is either replaced or merged depending on its patch strategy. The patch strategy is specified by the value of the patchStrategy key in a field tag in the Kubernetes source code. You can also see it in the [Kubernetes API documentation](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.21/#podspec-v1-core) under `patch strategy`.

</details>


> :bulb: credit to [this awesome blogpost](https://blog.atomist.com/kubernetes-apply-replace-patch/) for invaluable way of explaining this complex part of Kubernetes and Helm.

## Exercise

### Overview

- Helm install the chart
- Add new label to podspec template and apply
- Helm upgrade with the original one
- Change value of existing label and apply
- Helm upgrade with the original one

### Step by step instructions

<details>
<summary>Detailed instructions</summary>

- Navigate to the folder with the exercise content `cd release-rollback`
- Make sure that you have a release running in your cluster: `helm upgrade --install myapp sentence-app/`
- See that the pods are deployed: `kubectl get pods`
- Note down the revision number: `helm ls`
- Add a label to the deployment located in `extra/sentences-age-deployment.yaml`
- Apply the hand-edited deployment (its safe to ignore the warning about a missing annotation) `kubectl apply -f extra/sentences-age-deployment.yaml`
- See that the revision is still the same `helm ls`
- See the added label in the cluster `kubectl describe deployments.apps sentence-age `
- Make an upgrade back to the original version  `helm upgrade myapp sentence-app/`
- See that the new label is still persisted. `kubectl describe deployments.apps sentence-age`
- Edit one of the existing labels in `extra/sentences-age-deployment.yaml`
- Apply the hand-edited deployment `kubectl apply -f extra/sentences-age-deployment.yaml`
- See that the label has changed value in the cluster `kubectl describe deployments.apps sentence-age`
- Apply our original release `helm upgrade myapp sentence-app/`
- See that the label value have been reverted to the original value:`kubectl describe deployments.apps sentence-age`


</details>

### Clean up

- `helm uninstall myapp`
