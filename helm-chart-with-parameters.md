# Helm Chart with Parameters and Handling Whitespace

## Learning Goals

- Create Helm chart with parameterized values
- Modify Existing template files to use parameterized values
- Create `values.yaml` file to setup defaults for parameterized values
- Render Helm chart with different values

## Introduction

So far the Helm chart we have created only contains static values, meaning that every time we install it we get the exact same result.

In order to make the chart customizable, so that we can modify the chart for a specific use-case when we install it, we can use parameters.

Helm uses `go templates` under the hood, which enables powerful text templating of the kubernetes yaml files in the chart.

[Helm docs on values](https://helm.sh/docs/topics/charts/#templates-and-values)

[Golang docs on templates](https://pkg.go.dev/text/template?utm_source=godoc)

### Parameterizing Helm Template Files

Values are parameterized in Helm by replacing the value you want to parameterize with `{{ .Values.<valueName> }}`.

For example, if we have a deployment that species the number of replicas:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  ...
spec:
  replicas: 2
  ...
```

And we want to make the number of replicas configureable, we would change the `2` to `{{ .Values.replicas }}`.

Like so:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  ...
spec:
  replicas: {{ .Values.sentences.replicas }}
  ...
```

Which allows us to customize the number of replicas, each time we deploy the chart.

<details>
<summary>More details</summary>

### Golang Templates and Actions

Parameters are injected into template files using the `go template` syntax.

Golang templates uses `actions` whenever you want to specify a value that can be parametized.
`actions` are written using a double curly brace syntax: `{{ }}` so that everything within the two curly braces is interpreted by the parser, and not treated as actual text.

A trivial example of an action that returns the text "kubernetes" would look like this:
```
{{ "kubernetes" }}
```

That's not very useful though, so instead we will reference the `.Values` object which contains all of the values that we make available to Helm to use:

```
{{ .Values.orchestrationTool }}
```

Where we imagine that the value of `orchestrationTool=kubernetes`, which would result in the string "kubernetes" being injected when we render the yaml template.

> :bulb: When referencing the `.Values` object in Helm, you cannot use dashes (`-`), instead the convention is to use camel case.

### Helm Built-in Objects

Helm has a number of [built-in objects](https://helm.sh/docs/chart_template_guide/builtin_objects/) that you might want to use values from, such as the `.Release` object, which provided metadata about the current release of the chart.

For example we might want to include the name of our release in the names of the resources that are deployed, such that we can differentiate which release they belong to.

We do this by referencing the `Name` key of the `.Release` object:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  ...
  name: {{ .Release.Name }}-sentence
...
```

If the name of this release is `mySentences`, then the resulting name of the deployment would be `mySentences-sentence`.

</details>

### Values

There are two main ways for specifying the values that Helm should use when rendering our templates:

- Using the imperative `--set key=value` option on Helm commands.
- Using declarative `values.yaml` values file, which specify each value that can be parameterized.

The imperative approach is good for experiments or one off commands, while the declarative approach is good for repeatable installations and upgrades.

By convention the default values file is named `values.yaml`, and will be included automatically.
You can add custom values files by using the `--values` option.

```sh
$ helm install my-chart my-chart/ --values myValues.yaml
```

## Exercise

In this exercise we will add parameters to the sentences deployment, the "frontend" so to speak of the sentences application.

### Overview

- Modify Sentences Deployment
- Render Sentences Deployment Template from Command Line
- Parameterize the Container Image
- Create values file
- Render the Template with the Values File

You can use your chart from the previous exercise [create a Helm chart](create-a-helm-chart.md), or if you want a clean starting point, you can use the files in `helm-katas/helm-chart-with-parameters/start`.
If you get stuck, or want to see how the chart looks after completing the exercise, look at the chart in `helm-katas/helm-chart-with-parameters/done`.

### Step-by-Step

<details>
<summary>Steps:</summary>

**Modify Sentences Deployment**

The sentences deployment should be in your Helm chart under the templates directory: `sentence-app/templates/sentences-deployment.yaml`.

- Open this file in your text editor.

There are a lot of arguments in this deployment that we might want to parameterize, like the number of replicas, the container repository and tag or the resources allocation.

Let start by parameterizing the replicas.
This key is currently not in the deployment specification, so we have to add it.

- add `replicas: {{ .Values.sentences.replicas }}` to the yaml:

> :bulb: Be careful here, a Deployment object has two `spec` keys, one for the Deployment itself, and one for the `pod spec`.
> You must add the `replicas` key to the `Deployment spec`, that is the outermost `spec` key.
> The outermost spec key should be on line 8, so you can add the replicas key on a new line below.

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  ...
spec:
  replicas: {{ .Values.sentences.replicas }}
  ...
```

> :bulb: we prefix the `replicas` key with the name of deployment, in this case the `sentences` deployment, so that if we want to have a replicas value for each of the different deployment we access these with different prefixes.

**Render Sentences Deployment Template from Command Line**

- Try to render the yaml with a specified number of replicas:

```sh
$ helm template sentence-app/ --set sentences.replicas=2 --show-only templates/sentences-deployment.yaml

---
# Source: sentence-app/templates/sentences-deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app: sentences
    component: main
  name: sentences
spec:
  replicas: 2
  ...
```

As we can see the deployment would now create 2 replicas, you can try a few different number of replicas if you want.

**Parameterize the Container Image**

Next let's also parameterize the container repository and the tag.

- change:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  ...
spec:
  ...
  template:
    ...
    spec:
      containers:
      - image: releasepraqma/sentence:latest
```

To:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  ...
spec:
  ...
  template:
    ...
    spec:
      containers:
      - image: {{ .Values.sentences.repository }}:{{ .Values.sentences.tag }}
```

- Render the template file, and observe the new values getting reflected:

```sh
$ helm template sentence-app/ --set sentences.replicas=2 --set sentences.repository=myimage --set sentences.tag=mytag --show-only templates/sentences-deployment.yaml

---
# Source: sentence-app/templates/sentences-deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  ...
spec:
  replicas: 2
  selector:
    ...
  template:
    ...
    spec:
      containers:
      - image: myimage:mytag
        ...
```

Install the release to your cluster in order to verify that the image and tag used are rendered correct.

- `helm install sentences sentence-app/ --set sentences.replicas=2 --set sentences.repository=releasepraqma/sentence --set sentences.tag=latest`
- Verify that the deployment is healthy with `kubectl get deployments`

You can try a few different values for the repository and tag if you want.

**Create values file**

In the previous step we parameterized some of the values of the sentences deployment, and used cli options to specify the values.
As you can imagine when you have a lot values to parameterize, specifying all of them from the command line does not scale well.
Instead we will create a file `values.yaml` which will contain all of our values we want to use.

- Create a file named `values.yaml` in the root of your repository:

```sh
$ touch sentence-app/values.yaml
```

> :bulb: You can create the file any way you want to, just make sure that it is in the right location!

- Open it in your editor and add:

```yaml
sentences:
  replicas: 2
  repository: releasepraqma/sentence
  tag: latest
```

> :bulb: The structure of the yaml file defines the scope of the values.
> So to reference the replicas key, we would prefix it with the parent key, sentences: `sentence.replicas` and in the full helm object notation: `.Values.sentence.replicas`, just like we did above.

**Render the Template with the Values File**

- Render the sentences deployment again using the values from `values.yaml`:

```sh
$ helm template sentence-app --show-only templates/sentences-deployment.yaml
---
# Source: sentence-app/templates/sentences-deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  ...
spec:
  replicas: 2
  ...
  template:
    ...
    spec:
      containers:
      - image: releasepraqma/sentences:latest
        ...
```

> Note: Using a `values.yaml` file scales much better for larger charts.
>The values file also allows you to provide sensible defaults for all of the parameters that your chart has, as well as help the user to understand what values they should provide for each parameter.
>Since the values are kept in a file, the values file can be versioned with git or other tools, and can be used in for example a GitOps workflow.

> :bulb: It is convention to call values file `values.yaml`, though you can name it anything that you want.
> Helm will automatically use the values file named `values.yaml` if it exists, and other value files can be used with the option `--values myvalues.yaml`.
> If you use multiple values files, these will be merged by helm.

</details>

### Extra Exercise (optional)

If you have more time, or want to practice using values a bit more, then here are a couple of extra exercises:

<details>
<summary>Extra steps:</summary>

**Parameterize the two other Deployments**

Now we will add the same parameters to the two other deployments in the sentence application.

We will make the same changes that you made to `sentence-app/templates/sentences-deployment.yaml` to the other deployments:

- `sentence-app/templates/sentences-age-deployment.yaml`
- `sentence-app/templates/sentences-name-deployment.yaml`

We need to do one thing differently though, and that is that we need to specify which of the deployment the value belongs to, so that we can differentiate between them.

In the previous steps we referenced the values of the `sentences` value map, now we will be creating two new maps, `sentencesAge` and `sentencesName`.

You must use the appropriate map when making your changes to the deployment template files.

Instead of `{{ .Values.sentences.replicas }}` we would use `{{ .Values.sentencesAge.replicas }}` and `{{ .Values.sentencesName.replicas }}` respectively.

- Make the changes for the `replicas`, `repository` and `tag` values to the files `sentences-age-deployment.yaml` and `sentences-name-deployment.yaml`.

**Add new parameters to values.yaml**

In order to render our newly edited deployment templates we have to also provide values for them:

- Edit your `values.yaml` and add values for `sentencesAge` and `sentencesName`:

```yaml
sentences:
  replicas: 2
  repository: releasepraqma/sentence
  tag: latest

sentencesAge:
  replicas: 1
  repository: releasepraqma/age
  tag: latest

sentencesName:
  replicas: 1
  repository: releasepraqma/name
  tag: latest
```

- Render our to two modified deployment templates:

```sh
$ helm template sentence-app --show-only templates/sentences-age-deployment.yaml --show-only templates/sentences-name-deployment.yaml
---
# Source: sentence-app/templates/sentences-age-deployment.yaml
apiVersion: apps/v1
kind: Deployment
...
spec:
  replicas: 1
  ...
  template:
    ...
    spec:
      containers:
      - image: releasepraqma/age:latest
      ...
---
# Source: sentence-app/templates/sentences-name-deployment.yaml
apiVersion: apps/v1
kind: Deployment
...
spec:
  replicas: 1
  ...
  template:
    ...
    spec:
      containers:
      - image: releasepraqma/name:latest
        ...
```

**Use a Global Organization**

When we have values that repeat themselves, we can cut down on redundancy by parameterizing those as well, for example the organization in our container repository `releasepraqma`.

Instead of adding the organization name to each of our image repository values, we could instead use a Helm global value to set the organization name, and then prefix that to each of our instances.

> :bulb: Global values have some extended functionality when developing charts that include multiple sub charts, you can read more about it in the [documentation](https://helm.sh/docs/topics/charts/#global-values)

- Add the following section to your `values.yaml`:

```yaml
global:
  organization: releasepraqma
```

> :bulb: Global values can be referenced from the values object: `.Values.global.organization` for example.

- Edit each of the repository tags, such that they only include the unique name of each micro service:

```yaml
global:
  organization: releasepraqma

sentences:
  replicas: 2
  repository: sentences
  tag: latest

sentencesAge:
  replicas: 1
  repository: age
  tag: latest

sentencesName:
  replicas: 1
  repository: name
  tag: latest
```

- Edit each of your deployment template files to use the global organization name:

From:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  ...
spec:
  ...
  template:
    ...
    spec:
      containers:
      - image: {{ .Values.sentences.repository }}:{{ .Values.sentences.tag }}
```

To:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  ...
spec:
  ...
  template:
    ...
    spec:
      containers:
      - image: {{ .Values.global.organization }}/{{ .Values.sentences.repository }}:{{ .Values.sentences.tag }}
```

- Make the same change for the `age` and `name` deployments.

- Render the templates and verify that the repository names are correctly templated:

```sh
$ helm template sentence-app --show-only templates/sentences-deployment.yaml --show-only templates/sentences-age-deployment.yaml --show-only templates/sentences-name-deployment.yaml
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
      - image: releasepraqma/sentences:latest
      ...
---
# Source: sentence-app/templates/sentences-age-deployment.yaml
apiVersion: apps/v1
kind: Deployment
...
spec:
  ...
  template:
    ...
    spec:
      containers:
      - image: releasepraqma/age:latest
        ...
---
# Source: sentence-app/templates/sentences-name-deployment.yaml
apiVersion: apps/v1
kind: Deployment
...
spec:
  ...
  template:
    ...
    spec:
      containers:
      - image: releasepraqma/name:latest
        ...
```

</details>

### Clean up

- `helm uninstall sentences`
