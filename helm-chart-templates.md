# Helm Chart Named Templates

## Learning Goals

- Helm control flow and conditionals
- Use Helm templates
- Use context aware templates
- Use Helm templates in pipelines

## Introduction

Sometimes we want to do different things in our templates depending on how our deployment is to be configured.

For example we might want to specify a specific port to use for a NodePort service, but only if the service is of type NodePort.

To achieve this we can use conditionals and control flow in our Helm code.

Further we can define templates for chunks of code we want to reuse, but with different parameters throughout the code.

### Helm Control Flow and Conditionals

Helm has support for simple control flow using conditionals:

The simplest form is an `if` statement:

```
{{ if <pipeline> }}
    # do something
{{ end }}
```
> :bulb: We write `<pipeline>` here to indicate that the argument for the conditional can be as simple or complicated as needed.

> :bulb: Referencing a value is implcitly a valid pipeline!

A pipeline is evaluated as false if it returns:
- a boolean `false`
- a numeric zero `0`
- an empty string `""`
- an empty or null `nil`
- an empty collection `[] {} ()`

**All other cases are evaluated as true!**

Let's look at an example:
```
{{ if .Values.myVal }}
    myValPresent: true
{{ end }}
```
In the above example the string `myValPresent: true` is returned if the value `.Values.myVal` has been specified.

We can add an `else` statement:

```
{{ if .Values.myVal }}
    myValPresent: true
{{ else }}
    myValPresent: false
{{ end }}
```
In which case, if `.Values.myVal` is not set, the string `myValPresent: false` will be returned.

We can also use `else if` if we have multiple conditions:

```
{{ if .Values.myVal }}
    myValPresent: true
{{ else if .Values.myOtherVal }}
    myOtherValPresent: true
{{ else }}
    myValPresent: false
{{ end }}
```

[Control flow documentation](https://helm.sh/docs/chart_template_guide/control_structures/)

Helm has a number of functions that can be used in the conditionals like `and` and `eq`.

We can use the `eq` or 'equals' function to check if a value matches another predefined value:
```
{{ eq .Values.myVal "MatchThis" }}
```

We can even use these as arguments to other functions like the `and` function, which returns true if both of it's arguments are true:

```
{{ and (eq .Values.myVal "MatchThis") .Values.myOtherVal }}
```

The above example will return true if both: the `eq` function returns true, and the `.Values.myOtherVal` is present.

We could even use the above `and` example as the conditional for an `if` statement:

```
{{ if and (eq .Values.myVal "MatchThis") .Values.myOtherVal }}
    # do something
{{ end }}
```

There are a number of functions available that you can use to control the flow of your templates: [Control flow functions documentation](https://helm.sh/docs/chart_template_guide/function_list/#logic-and-flow-control-functions)

### Helm Templates

Helm has support for templating chunks of code that you want to reuse.

Templates are defined with the `define` action and the name of the template, and delimited with an `end`:

```
{{ define "myFirstTemplate" }}
foo:bar
{{ end }}
```

Templates are injected into your yaml with the `template` action and the name of the template:

```
{{ template "myFirstTemplate" }}
```

Templates can have it's own values, but these must be injected when calling the template.

```
{{ define "myTemplateWithArgs" }}
foo: {{ .Values.bar }}
{{ end }}
```

When invoking the template, we must pass an object containing the values, if we want to make the entire `.Values` available, we specify the entire context with a 'dot' `.`:

```
{{ template "myTemplateWithArgs" . }}
```

We can also specify a specific subset for the template to use:

```
{{ template "myTemplateWithArgs" .Values.myArgs }}
```

You can place helm templates anywhere in your `templates/` directory, but by convention, templates are usually placed in `templates/_helpers.tpl`.
Helm will not try to render `templates/_helpers.tpl` as part of your chart.

[Templates Documentation](https://helm.sh/docs/chart_best_practices/templates/#helm)

You can use templates in pipelines, but to do so you must use the `include` keyword instead of `template`:

```
{{ include "myTemplateWithArgs" .Values.myArgs | indent 4 }}
```

When using `include` you must specify the context for the template to use: `{{ include <templateName> <context> }}`

> :bulb: This is a limitation of golang templates, you can read more about in the [documentation](https://helm.sh/docs/howto/charts_tips_and_tricks/#using-the-include-function).

## Exercise

In this exercise we will first conditionally specify which nodePort our sentences service should use.

Then we will template the resource maps of the deployments, using conditional overrides.

### Overview

- Make the sentences service type parameterized
- Conditionally specify which nodePort to use
- Template resources map for deployments
- Conditionally override the template

You can use your Helm chart from the previous exercise as the starting point for this exercise.
Alternatively there is a Helm chart that picks up from the last exercise in `helm-katas/helm-chart-conditionals-templates/start` that you can use.
If you get stuck, or you want to see how the final chart looks, there is a solved version of the chart in `helm-katas/helm-chart-conditionals-templates/done`.

### Step-by-Step

<details>
<summary>Steps:</summary>

**Template resources map for deployments**

In the [previous exercise](https://github.com/eficode-academy/helm-katas/blob/main/helm-chart-whitespace-pipelines-functions.md) we learned how to parameterize the `resources` map of our deployments.

Now we would like to have a sensible default for our pod resources, with the ability to override the default on a per service basis.

To do this we will use a `template`.

Let's create a template file:

```sh
touch sentence-app/templates/_helpers.tpl
```
> :bulb: You can create the file in any way you want, it just has to be placed in the `templates` directory.

Open the file in your text editor and the following code:

```yaml
{{- define "resources" -}}
resources:
  requests:
    cpu: 0.50
    memory: "500Mi"
  limits:
    cpu: 1.0
    memory: "1000Mi"
{{- end -}}
```

This is just a simple template that will insert the above yaml map.

Let's use it in your sentences deployment, edit the file `sentence-app/templates/sentences-deployment.yaml` and change:

```yaml
apiVersion: apps/v1
kind: Deployment
...
spec:
  ...
  template:
    ...
    spec:
      containers:
      - ...
        resources:
          {{- .Values.sentences.resources | toYaml | nindent 10 }}
```

To:

```yaml
apiVersion: apps/v1
kind: Deployment
...
spec:
  ...
  template:
    ...
    spec:
      containers:
      - ...
        {{ template "resources" }}
```

Render the template:

```sh
helm template sentence-app --show-only templates/sentences-deployment.yaml
---
# Source: sentence-app/templates/sentences-deployment.yaml
apiVersion: apps/v1
kind: Deployment
...
spec:
  ...
  template:
    ...
    spec:
      containers:
      - ...
        resources:
  requests:
    cpu: 0.50
    memory: "500Mi"
  limits:
    cpu: 1.0
    memory: "1000Mi"
```

So far so good, but we have to fix the indentation.

In order to that we will change the `template` to an `include` in your sentences deployment, so that we can use a pipeline and the `nindent` function:

```yaml
apiVersion: apps/v1
kind: Deployment
...
spec:
  ...
  template:
    ...
    spec:
      containers:
      - ...
        {{- include "resources" . | nindent 8 }}
```

> :bulb: Also notice that we use `{{-` to remove whitespace before the template is injected.

Render the template again:

```sh
helm template sentence-app --show-only templates/sentences-deployment.yaml
---
# Source: sentence-app/templates/sentences-deployment.yaml
apiVersion: apps/v1
kind: Deployment
...
spec:
  ...
  template:
    ...
    spec:
      containers:
      - ...
        resources:
          requests:
            cpu: 0.50
            memory: "500Mi"
          limits:
            cpu: 1.0
            memory: "1000Mi"
```

There we go!

**Conditionally override the template**

But now our sentences deployment will always use the resources map specified in the template, so lets add a condition so that we can override it:

Edit your `_helpers.tpl` and the following if statement below the `define` line:

```yaml
{{ if .resources -}}
resources:
{{- .resources | toYaml | nindent 2 -}}
{{ else }}
```

We also need to add a `{{- end -}}` to delimit the `if` statement:

```yaml
    ...
    memory: "1000Mi"
{{- end -}}
{{- end -}}
```
> :bulb: We end up having two `{{- end -}}` at the end of the file because we have to delimit both the template `define` and the `if` statement.

The final `_helpers.tpl` should look like this:

```yaml
{{- define "resources" -}}
{{ if .resources -}}
resources:
{{- .resources | toYaml | nindent 2 -}}
{{ else }}
resources:
  requests:
    cpu: 0.50
    memory: "500Mi"
  limits:
    cpu: 1.0
    memory: "1000Mi"
{{- end -}}
{{- end -}}
```

Now we have modified our template, so that it expects a `context` that potentially contains a `resources` map.

This means that if the context indeed contains a `resources` map, then it will be rendered to yaml and returned, if not the default resources map is returned.

Next we edit our sentences deployment to pass the `.Values.sentences` context to the template:

Change:

```yaml
apiVersion: apps/v1
kind: Deployment
...
spec:
  ...
  template:
    ...
    spec:
      containers:
      - ...
        {{- include "resources" . | nindent 8 }}
```

To:

```yaml
apiVersion: apps/v1
kind: Deployment
...
spec:
  ...
  template:
    ...
    spec:
      containers:
      - ...
        {{- include "resources" .Values.sentences | nindent 8 }}
```

Since we defined a `resources` map in the `values.yaml` a few exercises ago, when we render the template we should see these values being used instead of the template one:


```sh
helm template sentence-app --show-only templates/sentences-deployment.yaml
---
# Source: sentence-app/templates/sentences-deployment.yaml
apiVersion: apps/v1
kind: Deployment
...
spec:
  ...
  template:
    ...
    spec:
      containers:
      - ...
        resources:
          limits:
            cpu: 0.5
            memory: 500Mi
          requests:
            cpu: 0.25
            memory: 100Mi
```

You can try to delete the `resources` map from your `values.yaml` and render the template again:

```sh
helm template sentence-app --show-only templates/sentences-deployment.yaml
---
# Source: sentence-app/templates/sentences-deployment.yaml
apiVersion: apps/v1
kind: Deployment
...
spec:
  ...
  template:
    ...
    spec:
      containers:
      - ...
        resources:
          requests:
            cpu: 0.50
            memory: "500Mi"
          limits:
            cpu: 1.0
            memory: "1000Mi"
```

Which means that the template will be used.

The neat thing here is that we can use the same template in the `age` and `name` deployments.

Now if we do not specify any resource limitations, the defaults will be used, but we can override those simply by adding limitations to the `values.yaml`.

</details>

### Extra Exercises

If you have more time, or you want to practice using templates and conditionals a bit more then you can do the extra steps:

<details>
<summary>Extras</summary>

Now that we have created a parameterized, conditional template for the resources map of the deployments in our charts, let's also apply it to the two other deployments in the chart:

- Change the `resources` map of the sentences-age and sentences-name deployments to use the template:

By changing the:
```
      - ...
        resources:
          requests:
            cpu: 0.25
          limits:
            cpu: 0.25
```
To:

```yaml
      - ...
        {{- include "resources" .Values.sentences | nindent 8 }}
```

When you make the change you should use the values map of the two other deployments: `sentencesAge` and `sentencesName` instead of the `sentences` map in the when you pass the context to the template: `.Values.sentences`.

> :bulb: Note that we use camel case for the deployment names when referencing the values object, as dashes `-` are not allowed in golang object names.

- Render the templates and verify that the age and name deployments use the default values from the `resources` template.

Now we can try to add some custom values to the `values.yaml` for the age and name deployments:

- Add a map for the `sentencesAge` and `sentencesName` deployments in your values file:

```yaml
...

sentencesAge:
  resources:
    requests:
      cpu: 1.25
      memory: "200Mi"
    limits:
      cpu: 1.50
      memory: "1200Mi"

sentencesName:
  resources:
    requests:
      cpu: 2.25
      memory: "500Mi"
    limits:
      cpu: 2.50
      memory: "2200Mi"
```

- Render the template again, and note that all of the deployments should have different resources specifications.

You can play around with setting some different values for each of the deployments.

</details>

### Food for thought

- When should you use templates in Helm?
- What is a good scope for a template?
- Can you make a chart too configurable?
