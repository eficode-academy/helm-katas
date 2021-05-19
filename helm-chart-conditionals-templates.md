# Helm Chart Conditionals and Templates

## Learning Goals

- Helm control flow and conditionals
- Use helm templates
- Use context aware templates
- Use helm templates in pipelines

## Introduction

Sometimes we to do different things in our templates depending on how our deployment is to be configured.

For example we might want to specify a specific port to use for a NodePort type service, but only want to that if the service is of type NodePort.

To achieve this we can use conditionals and control flow in our helm code.

Further we can define templates and use these customize sections of code that are often reused.

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

We can use the `eq` or 'equals' function to check if if value matches anothre predefined value:
```
{{ eq .Values.myVal "MatchThis" }}
```

We can even use these as arguments to other functions like the `and` function, which returns true if both of it's argumetns are true:

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

There are a number of functions available that you can use control the flow of your templates: [Control flow functions documentation](https://helm.sh/docs/chart_template_guide/function_list/#logic-and-flow-control-functions)

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

When invoking the template, we must pass an object containing the values, if we want to make the entire `.Values` avaialble, we specify the entire context with a 'dot' `.`:

```
{{ template "myTemplateWithArgs" . }}
```

We can also specify a specific subset for the template to use:

```
{{ template "myTemplateWithArgs" .Values.myArgs }}
```

You can place helm templates anywhere in your `templates/` directory, but by convention, templates are usually placed in `templates/_helpers.tpl`, which helm will not try to render as part of your chart.

[Templates Documentation](https://helm.sh/docs/chart_best_practices/templates/#helm)

You can use templates in pipelines, but in that case you must use the `include` keyword instead of `template`:

```
{{ include "myTemplateWithArgs" .Values.myArgs | indent 4 }}
```

When using `include` you must specify the context for the template to use: `{{ include <templateName> <context> }}`

> :bulb: This is a limitation of golang templates, you can read more about in the [documentation](https://helm.sh/docs/howto/charts_tips_and_tricks/#using-the-include-function).

## Exercise

In this exercise we will first conditionally specify which nodePort for our sentences service to use.
Then we will template the resource maps of the deployments, using conditional overrides.

### Overview

- Make the sentences service type parameterized
- Conditionally specify nodePort
- Template resources map for deployments
- Conditionally override the template

You can use your helm chart from the previous exercise as the starting point for this exercise.
Alternatively there is a helm chart that picks up from the last exercise in `helm-katas/helm-chart-conditionals-templates/start` that you can use.
If you get stuck, or you want to see how the final chart looks, there is a solved version of the chart in `helm-katas/helm-chart-conditionals-templates/done`.

### Step-by-Step

<details>
<summary>Steps:</summary>

**Make the sentences service type parameterized**

First let's have a look at the sentences service template, the file is located in `sentence-app/templates/sentences-svc.yaml`

```yaml
apiVersion: v1
kind: Service
metadata:
  labels:
    app: sentences
    component: main
  name: sentence
spec:
  ports:
  - port: 8080
    protocol: TCP
    targetPort: 8080
  selector:
    app: sentences
    component: main
  type: NodePort
```

As we can see the the type and ports for the service are hard-coded in the service template.

Let's make the type a parameter:

```yaml
apiVersion: v1
kind: Service
...
spec:
  ...
  type: {{ .Values.sentence.service.type }}
```

Add the type your `values.yaml`:
```yaml
sentence:
  ...
  service:
    type: ClusterIP
```

Let's try to render the service template:

```sh
$ helm template sentence-app --show-only templates/sentences-svc.yaml
---
# Source: sentence-app/templates/sentences-svc.yaml
apiVersion: v1
kind: Service
...
spec:
  ...
  type: ClusterIP
```

Sweet, that works.
Let's try to change the `type` in your `values.yaml` to `NodePort`.

**Conditionally specify nodePort**

Render the file again.

When using the `NodePort` service type Kubernetes allows us to specify which port we would like to use.
This argument is only relevant when using the `NodePort` service type, so let's make a conditional that only adds the port if the type is NodePort.

Add the port to your `values.yaml`:

```yaml
sentence:
  ...
  service:
    type: NodePort
    nodePort: 31234
```

We add the conditional to the `ports` map of the service spec in your sentences service template:

```yaml
apiVersion: v1
kind: Service
...
spec:
  ports:
  - port: 8080
    protocol: TCP
    targetPort: 8080
    {{ if and (eq .Values.sentences.service.type "NodePort") .Values.sentences.service.nodePort -}}
    nodePort: {{ .Values.sentences.service.nodePort }}
    {{- end }}
  ...
  type: {{ .Values.sentences.service.type }}
```
> :bulb: Notice the use of `{{-` and `-}}` to remove whitespace around the injected value in the if statement.

The above `if` statement is true if both the `type` is set to `NodePort` and the value `nodePort` has been specified.

Try to render the template:

```sh
$ helm template sentence-app --show-only templates/sentences-svc.yaml
---
# Source: sentence-app/templates/sentences-svc.yaml
apiVersion: v1
kind: Service
...
spec:
  ports:
  - port: 8080
    protocol: TCP
    targetPort: 8080
    nodePort: 31234
  ...
  type: NodePort
```

Now let's try to change the type back to `ClusterIP` in values file, and render the template again:

```sh
$ helm template sentence-app --show-only templates/sentences-svc.yaml
---
# Source: sentence-app/templates/sentences-svc.yaml
apiVersion: v1
kind: Service
...
spec:
  ports:
  - port: 8080
    protocol: TCP
    targetPort: 8080
  ...
  type: ClusterIP
```

So that we can verify that the `nodePort` key is only added when the `type` is set to `NodePort`.

**Template resources map for deployments**

In the previous exercise we learned how to parameterize the `resoruces` map of our deployments.

Now we would like to have a sensible default for our pod resources, with the ability to override the default on a per service basis.

To do this we will use a `template`.

Let's create a template file:

```sh
touch helm-chart/sentence-app/templates/_helpers.tpl
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

But now our sentences deployment will always use the resources map specified in the template, so lets a condition so that we can override it:

Edit your `_helpers.tpl` and the following if statement below the `define` line:

```yaml
{{ if .resources -}}
resources:
{{- .resources | toYaml | nindent 2 -}}
{{ else }}
```

We also need to add a `{{ end }}` to delimit the `if` statement:

```yaml
    ...
    memory: "1000Mi"
{{- end -}}
{{- end -}}
```

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

Now we setup our template so that it expects to be passed a context that potentially contains a `resources` map.

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

The neat thing here is that we can use the same template in the `age` and `name` deployments, such if we add a `resources` map for those functions, then those would override the template, which means that we can specify custom resources requests and limits for each deployment, just by setting the values in our `values.yaml`.

</details>

## Cleanup
