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

This exercise is deliberately made vague, with only the overview, and not step by step.

But we will provide a section with hints in :bulb:.

### Overview

The exercise resides in the `helm-lint/start` folder.

- Try to deploy the helm chart and see what happens
- Run `helm lint` to help you identify the problems.
- Install the helm plugin [helm kubeval](https://artifacthub.io/packages/helm-plugin/kubeval/kubeval): `helm plugin install https://github.com/instrumenta/helm-kubeval`.
- Run `helm kubeval` to help you investigate further.
- Deploy the fixed chart

### Hints

<details>
<summary> :bulb: First hint</summary>

> The first error you get when running helm lint is that the chart.metadata.type is wrong. That property is set in the Chart.yaml. Look at the spelling of application.

</details>

<details>
<summary> :bulb: Second hint</summary>

> Next output tells you that there is something wrong at line 12 in templates/sentences-name-svc.yaml. This is not exactly correct, but has to do with the way YAML gets linted. The real problem comes a few lines above and has something to do with indentation.

</details>

<details>
<summary> :bulb: Third hint</summary>

> By now, helm lint should not give you any more errors, and you need to use helm kubeval. The error it gives you are aimed at the kubernetes objects. Therefore the mistypo in "port" should be fairly easy to spot in `sentence-app/templates/sentences-age-svc.yaml`

</details>

### Clean up

If anything needs cleaning up, here is the section to do just that.
https://artifacthub.io/packages/helm-plugin/kubeval/kubeval