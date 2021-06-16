# Helm Chart Conditionals and Control Flow

## Learning Goals

- Use Helm control flow and conditionals
- Use if statements
- Use equality operators

## Introduction

Sometimes we want to do different things in our templates depending on how our deployment is to be configured.

For example we might want to specify a specific port to use for a NodePort service, but only if the service is of type NodePort.

To achieve this we can use conditionals and control flow in our Helm code.

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

## Exercise

In this exercise we will conditionally specify which nodePort our sentences service should use, and only insert the value, if the service type is NodePort.

### Overview

- Make the sentences service type parameterized
- Conditionally specify which nodePort to use
- Only insert the nodePort if service type is NodePort

You can use your Helm chart from the previous exercise as the starting point for this exercise.
Alternatively there is a Helm chart that picks up from the last exercise in `helm-katas/helm-chart-conditionals/start` that you can use.
If you get stuck, or you want to see how the final chart looks, there is a solved version of the chart in `helm-katas/helm-chart-conditionals/done`.

### Step-by-Step

<details>
<summary>Steps:</summary>

**Make the sentences service type parameterized**

First let's have a look at the sentences service template, the file is located in `sentence-app/templates/sentences-svc.yaml`

The type and ports for the service are hard-coded in the service template.

- Let's make the type a parameter:

```yaml
apiVersion: v1
kind: Service
...
spec:
  ...
  type: {{ .Values.sentences.service.type }}
```

- Add the type to your `values.yaml`:

```yaml
sentences:
  ...
  service:
    type: ClusterIP
```

- Let's try to render the service template:

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

- Let's try to change the `type` in your `values.yaml` to `NodePort`.

- Render the template again, and verify that it is now set to `NodePort`.

**Conditionally specify which nodePort to use**

When using the `NodePort` service type, Kubernetes allows us to specify which port we would like to use.
This argument is only relevant when using the `NodePort` service type, so let's make conditional that only adds the `nodePort` key if the value is set.

> :bulb: Note the difference in case for `NodePort` and `nodePort`.

- Choose a random number between 30000 and 32767 and add it as the `nodePort` in your `values.yaml`:

> :bulb: the range `30000-32767` is the default range for NodePorts in Kubernetes.
> Only one service can occupy a specific port at a time, therefore if you are multiple people doing the exercises together, then everyone must choose a unique port, so as not to conflict.
> [You click here if you need inspiration for your unique port number](https://www.randomlists.com/random-numbers?dup=false&qty=1&max=32767&min=30000)
> In the example we will use `31234`, but you should choose a different one.

```yaml
sentences:
  ...
  service:
    type: NodePort
    nodePort: 31234
```

- Add the conditional to the `ports` map of the service spec in your sentences service template:

```yaml
apiVersion: v1
kind: Service
...
spec:
  ports:
  - port: 8080
    protocol: TCP
    targetPort: 8080
    {{- if .Values.sentences.service.nodePort }}
    nodePort: {{ .Values.sentences.service.nodePort }}
    {{- end }}
  ...
  type: {{ .Values.sentences.service.type }}
```

> :bulb: Notice the use of `{{-` to remove whitespace around the injected value in the if statement.

Now the `nodePort` key will be inserted if the key is set.

- Render the template:

```sh
helm template sentence-app --show-only templates/sentences-svc.yaml
---
# Source: sentence-app/templates/sentences-svc.yaml
apiVersion: v1
kind: Service
metadata:
...
spec:
  ports:
  - port: 8080
    protocol: TCP
    targetPort: 8080
    nodePort: 31234
  ...
```

- Try to comment out the `nodePort` in your values file and render the template again:

```yaml
sentences:
  ...
  service:
    type: ClusterIP
    # nodePort: 31234
```

When we render the template again the `nodePort` will not be shown, as the value is not set, which means the if statement is false.

```sh
helm template sentence-app --show-only templates/sentences-svc.yaml
---
# Source: sentence-app/templates/sentences-svc.yaml
apiVersion: v1
kind: Service
metadata:
...
spec:
  ports:
  - port: 8080
    protocol: TCP
    targetPort: 8080
  ...
```

- Uncomment the `nodePort` line in your values file.

**Only insert the nodePort if service type is NodePort**

As we stated in the beginning of the exercise the `nodePort` key is only used if the service type is `NodePort`.
So we can add an extra condition to our service, so that the `nodePort` key is only used if the service is indeed set to `NodePort`.

To do this we will use the `and` function.

The `and` function takes two arguments, both of which must be true:

```yaml
{{ and <pipeline> <pipeline> }}
```

An argument can be the result of another function or pipeline of functions.
These must be placed in parentheses.

We can use the `eq` or 'equals' function to check that the service is type is "NodePort":

```yaml
{{ eq .Values.MyValue "MyValue" }}
```

When we put it all together we get an if statement, where the condition is the `and` function, where the first of the arguments to the `and` function is the `eq` function that checks if the service type is `NodePort` and the second argument is a check that the `nodePort` value has been set.

```yaml
{{ if and (eq .Values.sentences.service.type "NodePort") .Values.sentences.service.nodePort }}
```

Thus only when both conditions are met, the value will be inserted.

- Edit your service file with the new if statement:

```yaml
apiVersion: v1
kind: Service
...
spec:
  ports:
  - port: 8080
    protocol: TCP
    targetPort: 8080
    {{- if and (eq .Values.sentences.service.type "NodePort") .Values.sentences.service.nodePort }}
    nodePort: {{ .Values.sentences.service.nodePort }}
    {{- end }}
  ...
  type: {{ .Values.sentences.service.type }}
```

- Try to render the template:

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

- Try to change the type back to `ClusterIP` in the values file, and render the template again:

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

</details>
