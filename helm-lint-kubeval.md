# Helm Lint and Kubeval
## Learning Goals

- Debug helm chart through `helm lint`
- Debug kubernetes objects renderes through helm with the `helm kubeval` plugin

## Helm linting and Kubeval plugin

Linting and evaluating the kubernetes objects that helm produces through `helm lint` and `kubeval` can help you speed up you development of helm charts, by catching your typos and indentation mistakes.

We are going to use two different tools here:

- `helm lint` will lint your helm chart, catching mandatory fields not set in the helm objects, and validate all of the YAML containing actions.
- `helm kubeval` is a plugin to render your chart into kubernetes objects, and then run the `kubeval` kubernetes YAML linter on them. It functions similarly to `helm lint`, but for kubernetes objects.

Both tools help you, but with different scope, one for the template files themselves, and the other for the rendered templates.

## Exercise

This exercise is about using the linters to find mistakes in a chart we have prepared.
Therefore there is no step by step walk through in this exercise, as you will have to use the tools to find the mistakes and fix them.

But we will provide a section with hints in :bulb:.

### Overview

The exercise resides in the `helm-lint/start` folder.

- Run `helm lint sentence-app/` to help you identify the problems of the chart.
- Run `helm kubeval sentence-app/` to help you investigate further.
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

> By now, helm lint should not give you any more errors, and you need to use helm kubeval. The error it gives you are aimed at the kubernetes objects. Therefore the typo in "port" should be fairly easy to spot in `sentence-app/templates/sentences-age-svc.yaml`

</details>

### Clean up

If anything needs cleaning up, here is the section to do just that.

### Resources

https://artifacthub.io/packages/helm-plugin/kubeval/kubeval