# Template headline

## Learning Goals

- provide a list of goals to learn here

## Helm linting and Kubeval plugin

Linting and evaluating the kube objects that helm produces through kubeval can help you speed up you development of helm charts, by pointing towards your typos and indentation mistakes.

We are going to use two different tools here:

- `helm lint` will lint your helm chart, catching mandatory fields not set in the helm objects, and YAML validate all your kubernetes objects with actions.
- `helm kubeval` is a plugin to render your chart into kubernetes objects, and run `kubeval` on them. It performs the same as `helm lint` does, but for you kubernetes objects.

So both tools help you, but with different scope.

## Exercise

This exercise is deliberately made a bit vague, but we will provide hints along the way in :bulb: drop down sections.

### Overview

The exercise resides in the `helm-lint/start` folder.

- Try to deploy the helm chart and see what happens
- Run `helm lint` to help you identify the problems.
- Install the helm plugin `helm kubeval`
- Run `helm kubeval` to help you investigate further.
- Deploy the fixed chart

### Hints

<details>
<summary>First hint</summary>

**take the same bullet names as above and put them in to illustrate how far the student have gone**

- all actions that you believe the student should do, should be in a bullet

> :bulb: Help can be illustrated with bulbs in order to make it easy to distinguish.

</details>

<details>
<summary>Second hint</summary>

**take the same bullet names as above and put them in to illustrate how far the student have gone**

- all actions that you believe the student should do, should be in a bullet

> :bulb: Help can be illustrated with bulbs in order to make it easy to distinguish.

</details>

<details>
<summary>Third hint</summary>

**take the same bullet names as above and put them in to illustrate how far the student have gone**

- all actions that you believe the student should do, should be in a bullet

> :bulb: Help can be illustrated with bulbs in order to make it easy to distinguish.

</details>

### Clean up

If anything needs cleaning up, here is the section to do just that.
https://artifacthub.io/packages/helm-plugin/kubeval/kubeval