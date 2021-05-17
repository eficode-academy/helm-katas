# Helm Chart with Parameters and Handling Whitespace

## Learning Goals

- Create helm chart with parameterized values
- Modify Existing template files to use injected values
- Create `values.yaml` file to setup defaults for injected values
- Render Helm chart with different values
- Manage whitespace around injected values

## Introduction

So far the helm chart we have created only has static values, meaning that every time we install it we get the exact same result.

In order to make the chart customizable, so that we can modify the chart for a specific use-case when we install it, we can use parameters.

Helm uses `go templates` under the hood, which enables powerful text templating of the kubernetes yaml files in the chart.

<!-- TODO add link to doc -->

Since we are templating `yaml` files we have to be careful with getting the indentation of parameterized values right, so therefore the templating includes functionality for handling whitespace.

### Helm Template Files

Parameters can be injected into template files using the `go template` syntax.
Golang templates uses `actions` whenever you want to specify a value that can be parametized.
`actions` are written using a double curly brace syntax: `{{ }}` so that everything within the two curly braces is interpreted by the parser, and not treated as actual text.

A trivial example of an action that returns the text "kubernetes" would look like this:
```
{{ "kubernetes" }}
```

That's not very useful though, so instead we will reference the `.Values` object which contains all of the values that we make available to helm to use:

```
{{ .Values.orchestration-tool }}
```
Where we imagine that the value of `orchestration-tool=kubernetes`, which would result in the string "kubernetes" being injected when we render the yaml template.

If we, for example, wanted to template the name of the resources of our sentence application:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  ...
  name: sentences
...
```
Then we would change the value of the `name` key to an action that inserts the name of the release in front of the name of the deployment:
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  ...
  name: {{ .Release.Name }}-sentence
...
```


### Values

There are two main ways for specifying the values that Helm should use when rendering our templates:
- Using the imperative `--set key=value` option on `helm install` and `helm upgrade` commands.- Or using declarative values file, which are yaml files that specify each value that can be parametized.

The imperative approach is good for experiments or one off commands, while the declarative approach is good for repeatable installations and upgrades.





## Exercise

### Overview

- Modify sentences deployment
- Render the sentences deployment with different cli arguments
- Create values file
- Render the sentences deployment with the values file
- Modify sentences service and update values file
- Render both templates with the values file
- Modify the sentences deployment to trim whitespace

You can use your chart from the previous exercise, or if you want a clean starting point, you can use the files in `helm-katas/helm-chart-with-parameters/chart-with-parameters-start`.
If you get stuck, or want to see how the chart looks after completing the exercise, look at the chart in `helm-katas/helm-chart-with-paramters/chart-with-parameters-done`.

### Step-by-Step

<!-- <details> -->
      <!-- <summary>Steps:</summary> -->
<!-- </details> -->

TODO

## Cleanup

TODO
