# Helm Chart Whitespace Handling, Pipelines and Functions

## Learning Goals

- Use functions to transform values in actions
- Use pipelines to combine different functions
- Use Helm features and functions to handle whitespace

## Introduction

When writing Helm charts you might need to do more elaborate templating than simply injecting values from parameters.

For example setting default values, ensuring that strings are quoted and that blocks are indented correctly.

`Functions` helps us do these things, and we can even combine different functions using `pipelines`.

Finally we will look at different ways of managing whitespace in Helm templates.

### Whitespace Handling with Helm

Helm provides different ways of handling whitespace around actions.
The actions delimiters `{{` and `}}` can be augmented with a dash `-`:
- `{{-` will consume all whitespace to left of the action, including newlines.
- `-}}` will consume all whitespace to the right of the action, including newlines.

Whitespace includes all `spaces`, `tabs` and `newlines`!

<details>
<summary>An example:</summary>

```
PRE
  {{- "mytext" -}}
    POST
```

Would render to:
```
PREmytextPOST
```

Because all of the whitespace around the action will be consumed by the `{{-` and `-}}`, until non-whitespace characters are encountered.
</details>

There are also functions for managing whitespace, like adding a configurable amount of indentation with the `indent` function:

```
{{ indent 4 .Values.myVal }}
```

Would add 4 spaces in front of the injected value.

Similarly `nindent` functions the same as `indent`, but also adds a newline before the line being indented.

> :bulb: Documentation links:
>
>[Documentation for controlling whitespaces](https://helm.sh/docs/chart_template_guide/control_structures/#controlling-whitespace)
>
>[Indent and nindent documentation](https://helm.sh/docs/chart_template_guide/function_list/#indent)

### Helm Functions

Helm has a number of `functions` available that enable more elaborate templating.

Functions are used in actions and take at least one argument:

```
{{ function argument1 argument2 }}
```

The result of applying the argument to the function will be returned by the action.

<details>
<summary>An example:</summary>

A useful and simple example of a function could be to add quotes to a string:

```yaml
shouldBeAString: {{ quote .Values.myString }}
```

We assume that `myString=FooBar`, thus the result of the function will be `shouldBeAString: "FooBar"`.

</details>

> :bulb: Documentation links:
>
> [Helm Documentation on using functions](https://helm.sh/docs/chart_template_guide/functions_and_pipelines/#helm)
>
> [Full list of available functions](https://helm.sh/docs/chart_template_guide/function_list/)

##### Custom Functions

Helm is distributed as a static binary, so it only includes the functions that the binary was compiled with.

If you need to use custom functions, you can 'attach' an external binary as a post-renderer, which will run on the templates after helm has templated them, but before installing them.

[Post rendering documentation](https://helm.sh/docs/topics/advanced/#post-rendering)

### Helm Pipelines

Pipelines allow us to use the output of one function as the input of another function:

```
{{ function1 | function2 }}
```

Where the result of function1 is used as the argument for function2, and the result of function2 is returned from the action.

<details>
<summary>More information:</summary>

Pipelines are written using the "pipe" character `|`.

We can rewrite our quoting example above with a pipeline:

```yaml
shouldBeAString: {{ .Values.myString | quote }}
```

> :bulb: Referencing a value is actually an implicit function!

Which will produce the exact same result.

We can use as many functions as we want to in a pipeline.

For example if we wanted to make sure that our string only contains lower case characters, we can use the `lower` function in our pipeline:

```yaml
shouldBeALowerCaseString: {{ .Values.myString | lower | quote }}
```

Which would first change the value of `myString=FooBar` to lowercase, and then add quotes.

The result would be: `shouldBeALowerCaseString: "foobar"`

> :bulb: Documentation links:
>
> [Documentation on using pipelines](https://helm.sh/docs/chart_template_guide/functions_and_pipelines/#pipelines)

</details>

## Exercise

In this exercise we will use functions, pipelines and whitespace handling to parameterize the resource section of the sentences deployment.

### Overview

- Make CPU and memory limits configurable for the sentences deployment
- Add default values for CPU resources requests and limits
- Use functions to render values map to yaml
- Use a pipeline to properly indent resources map
- Make the resources pipeline more readable by managing whitespace

You can use your Helm chart from the previous exercise as the starting point for this exercise.
Alternatively there is a Helm chart that picks up from the last exercise in `helm-katas/helm-chart-whitespace-functions-pipelines/start` that you can use.
If you get stuck, or you want to see how the final chart looks, there is a solved version of the chart in `helm-katas/helm-chart-whitespace-functions-pipelines/done`.

### Step-by-Step

<details>
<summary>Steps:</summary>

**Make CPU and Memory Limits Configurable for the Sentences Deployment**

So far our sentences deployment has a hard-coded definition of each pods resource limits, in this case the CPU request and the CPU limit:

<details>
<summary>:bulb: What is the difference between resource requests and limits?</summary>
In kubernetes each deployment can specify a request for a number resources to be allocated for a given pod.

This is used by the scheduler to ensure that there are enough resources available on a given node to run the pod.

For CPU this is specified as a decimal, where `1.0` is one CPU core.

The limit key specifies the maximum of a resource a pod may consume, and can be set to the same as the request, or higher to allow for the pod to consume more resources if needed.

For CPU, when a pod reaches it's limit, it will be throttled, if a pod reaches it's memory limit, it will be stopped, so configure these wisely!

You can read more about it in the Kubernetes [documentation](https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/).
</details>

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
          requests:
            cpu: 0.25
          limits:
            cpu: 0.25
```

Let's make the CPU request and limit configurable, we learned in the last exercise to use `actions` to accomplish this:

- Change your `sentences-deployment.yaml` to have parameterized CPU resource requests and limits like below.

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
          requests:
            cpu: {{ .Values.sentences.cpuRequest }}
          limits:
            cpu: {{ .Values.sentences.cpuLimit }}
```

Check that your parameters are working:

```sh
$ helm template sentence-app --show-only templates/sentences-deployment.yaml --set sentences.cpuRequest=0.25 --set sentences.cpuLimit=0.5
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
            cpu: 0.25
          limits:
            cpu: 0.5
```

**Add Default Values for CPU Resources Requests and Limits**

Maybe we don't always know what kind of limitations we want to put on our pods, but declaring a value like we do above means that we **have** to provide a value to render the template.

Luckily we can use the `default` [function](https://helm.sh/docs/chart_template_guide/function_list/#default) to specify a default value for our values:

```
{{ default "defaultValue" .optionalValue }}
```

Change your sentences deployment to use the `default` function:

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
          requests:
            cpu: {{ default 0.50 .Values.sentences.cpuRequest }}
          limits:
            cpu: {{ default 0.75 .Values.sentences.cpuLimit }}
```

Now try to render the template again, without specifying any argument for the values:

```sh
$ helm template sentence-app --show-only templates/sentences-deployment.yaml
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
          limits:
            cpu: 0.75
```

**Use Functions to Render Values Map to Yaml**

But what about memory requests and limits?

We could simply add parameterized, defaulted values for memory limits and requests:

```yaml
resources:
  requests:
    cpu: {{ default 0.50 .Values.sentences.cpuRequest }}
    memory: {{ default "100Mi" .Values.sentences.memoryRequest }}
  limits:
    cpu: {{ default 0.75 .Values.sentences.cpuLimit }}
    memory: {{ default "500Mi" .Values.sentences.memoryLimit }}
```

This is getting a bit hard to read, also we would be enforcing these defaults on anyone who installed the chart, thinking it might use the cluster defined resource request and limit defaults.

So instead let's make the entire `resources` map parameterized, but only for the values that are provided by the user.

- Add CPU resource values to our `values.yaml`:

```yaml
sentences:
  ...
  resources:
    requests:
      cpu: 0.25
    limits:
      cpu: 0.50
```

- Modify our `templates/sentences-deployment.yaml`:

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
          {{ .Values.sentences.resources }}
```

Sadly this doesn't quite do what we want to do, if we try to render the template:

```sh
$ helm template sentence-app --show-only templates/sentences-deployment.yaml
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
          map[limits:map[cpu:0.5] requests:map[cpu:0.25]]
```

Helm helpfully attempts to insert our `resources` map from the values file, but inserts it as a golang map of maps, which we cant use.

Fortunately we can use the `toYaml` function to render the golang map as yaml:

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
          {{ toYaml .Values.sentences.resources }}
```

- Add the `toYaml` function to the action like shown above
- Render the template:

```sh
$ helm template sentence-app --show-only templates/sentences-deployment.yaml
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
requests:
  cpu: 0.25
```

It looks better, the `resources` map from the values file is rendered as proper yaml, but the indentation is not correct.

**Use a Pipeline to Properly Indent Resources Map**

To fix the indentation we can use the `indent` function to add a number of spaces in front of our rendered yaml.

That means we have to take the result of our `toYaml` function and use it as the input of the `indent` function, so we will use a pipeline:

```
{{ toYaml .Values.sentences.resources | indent 10 }}
```

Let's edit our sentences deployment:

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
{{ toYaml .Values.sentences.resources | indent 10 }}
```

> :bulb: Notice that we remove all indentation in front of our action, as the `indent` function will handle creating all of the required whitespace.

> :bulb: The 10 argument for the indent function is the number of characters to indent using spaces.
> Your text editor likely has a character counter to allow you to see how many characters on the current line your caret is at, otherwise you can simply count the number spaces the block would have been indented.

Now let's try to render the template again:

```sh
$ helm template sentence-app --show-only templates/sentences-deployment.yaml
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
          requests:
            cpu: 0.25
```

Success! Our resources are now properly formatted and indented.

**Make the Resources Pipeline more Readable by Managing Whitespace**

While the resources parameterization we have created so far works, it looks a bit odd without any indentation in the `templates/sentences-deployment.yaml`.
We can fix that by controlling the whitespaces with functions.

- change the first function call to `toYaml` to a pipeline:

```yaml
{{ .Values.sentences.resources | toYaml | indent 10 }}
```

This is stylistic change, and produces the exact same result.
The pipeline syntax seems to be preferred, but you can use whichever style you prefer.

Next we use a `{{-` to consume all whitespace to the left of the action.

- Change the action adding the whitespace handling in the beginning like the example below:

```
{{- .Values.sentences.resources | toYaml | indent 10 }}
```

> :bulb: rendering this will result in an error because newlines are also considered "whitespace".
> This means that there will **not be any whitespace** before our rendered resource map, so we need to add a newline.

We can add a newline before our indented block by using the `nindent` function instead of the `indent` function.

Since we add the newline and all of the whitespace with functions, we can write the action at the logical indentation in the template yaml.

- Change the `indent` function to `nindent` like the example below

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

- Test that it works by letting Helm render it: `helm template sentence-app --show-only templates/sentences-deployment.yaml`

The resulting template is much cleaner and easier to read.

- Try to add memory specifications to your `values.yaml`:

```yaml
sentences:
  ...
  resources:
    requests:
      cpu: 0.25
      memory: "100Mi"
    limits:
      cpu: 0.50
      memory: "500Mi"
```

- And render the template:

```sh
$ helm template sentence-app --show-only templates/sentences-deployment.yaml
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

And we can see that the memory specifications are injected correctly!

</details>
