# Helm Chart with Parameters

## Learning goal

- Learning how to use parameters in Helm such as
  - Configurable Kubernetes resource naming (using
    Helm built-in values)
  - Configurable number of POD replicas in
    deployments (using single-value parameters
    with defaults)
  - Configurable resource settings (using
    list-type value parameters)
  - Optional definitions (using if/else
    constructs)

## Introduction

This exercise will demonstrate adding parameters
to a Helm chart to allow customization of the
installed application.

This exercise extends the Helm chart created in
the [create-a-helm-chart](create-a-helm-chart.md)
exercise.

when we initially created the Helm chart we did it
with an empty values.yaml file. Now we want to
inject parameters into the chart.

For that we need to update the chart template
files such that values are inserted at the
appropriate places. We do that using the `{{` and
`}}` template language constructions.

<details>
      <summary>More information</summary>
Helm templateting can both help us with inline
replacement of values, and whole sections to be
inserted into our Yaml files.

Kubernetes allows us to define which port to use
for services of type NodePort. In the exercise
below we will customize the Kubernetes YAMl for
this scenario.

</details>

## Exercise

### Overview

- Add content to `values.yaml` for replicas,
  image, and service
- Add parameter references to the chart
- Add conditional rendering of `nodePort` type
- Add a `values-resources.yaml` values files for
  your customizations of the helm chart.

### Step-by-step

<details>
      <summary>More information</summary>

**Add content to `values.yaml`**

- Replace the `values.yaml` with the following
  content:

```yaml
## sentence configured the main sentences micro-service
sentences:
  ## replicas is the number of POD replicas
  replicas: 1
  ## image contains the POD container image parameters
  image:
    repository: releasepraqma/sentence
    tag: latest
  ## resource requests and limits
  resources: {}
  service:
    ## type is the type of the service fronting the sentence application
    type: NodePort
    ## nodePort is the actual port used for NodePort-type services
    #nodePort:
```

<details>
      <summary>:bulb: This `values.yaml` file shows:</summary>

> - Using double `##` for documentation comments
>   and single `#` for commented-out fields. This
>   is not an official Helm standard but a very
>   common approach.
> - Image name and tag are separately
>   configurable. This is a common pattern because
>   the tag typically is configured independently
>   from the image name.
> - POD resource settings are not defaulted in
>   `values.yaml`. We cannot know the proper
>   settings and instead we provide an empty map
>   and allows users to set appropriate resource
>   requests.
> - The `nodePort` parameter (which is only
>   relevant in case the Kubernetes service is of
>   type `NodePort`) has no default. Again, we
>   cannot know the proper default but instead
>   indicate in the `values.yaml` file that there
>   is a parameter that can be configured. Another
>   typically used alternative for string-type
>   values with no good default is to leave them
>   as empty strings.

</details>

<details>
      <summary>:bulb: helm best practice</summary>
> :bulb: Note that the Helm best practices
> suggests using a flat values file vs. a nested
> one as shown above. Real-world charts however,
> very often use nested values and this is
> probably the defacto standard. For more info
> look
> [here](https://helm.sh/docs/chart_best_practices/values/#flat-or-nested-values)
</details>

**Adding Parameters to the Chart**

In the template file
`sentence-app/templates/sentences-deployment.yaml`
file:

- locate the first `spec` of the deployment
  manifest, i.e. the part that starts with:

```yaml
spec:
  selector:
    matchLabels:
```

- add a line with a `replicas` specification as
  follows:

```
spec:
  replicas: {{ .Values.sentences.replicas }}
  selector:
    matchLabels:
```

- Verify the rendering as follows:
  `helm template sentence-app/ --show-only templates/sentences-deployment.yaml`

> :bulb: Since `values.yaml` have a default
> replica count of 1 that is what we see in the
> rendered `sentences-deployment.yaml` template.

- Try explicitly override the value in the helm
  invocation as follows:

- `helm template sentence-app/ --show-only templates/sentences-deployment.yaml --set sentences.replicas=3`

- In the same yaml file, change the container
  image specification:

```
      - image: {{ .Values.sentences.image.repository }}:{{ .Values.sentences.image.tag }}
```

- Change the Deployment name as follows:

```
  name: {{ .Release.Name }}-sentences
```

> :bulb: The resource section is slightly
> different since this is not a single value but
> instead a full YAML map. Also, we do not know
> whether the user will specify limits or requests
> for either CPU or memory since all this is
> related to the actual usage.
>
> Instead we simply insert the full YAML as given
> by the user. To do this we use a Helm function
> and pipeline as follows.

- Change the hard-coded resource settings from:

```
        resources:
          requests:
            cpu: 0.25
          limits:
            cpu: 0.25
```

to:

```yaml
        resources:
{{ toYaml .Values.sentences.resources | indent 10 }}
```

- Validate the resource setting:
  `helm template sentence-app/ --set sentences.resources.requests.cpu=0.25 --show-only templates/sentences-deployment.yaml`

> :bulb: Pay particularly attention to indentation
> on the rendered YAML.

**Adding Conditional Rendering of Template**

In the template file
`sentence-app/templates/sentences-svc.yaml` file,
locate the specification of the service type:

```
  type: NodePort
```

- Change this line and add nodeport specification
  as follows:

```
  type: {{ .Values.sentences.service.type }}
```

- Locate the specification of the port mapping:

```
  ports:
  - port: 8080
    protocol: TCP
    targetPort: 8080
```

- Add the conditional nodeport specification as
  follows:

```
  ports:
  - port: 8080
    protocol: TCP
    targetPort: 8080
    {{ if and (eq .Values.sentences.service.type "NodePort") .Values.sentences.service.nodePort -}}
    nodePort: {{ .Values.sentences.service.nodePort }}
    {{- end }}
```

Note the post-fix used by Helm notation in the
above specification.

- Test that helm will not write a `nodePort:`
  section under ports, if you do not specify one:
  `helm template sentence-app/ --show-only templates/sentences-svc.yaml --set sentences.service.type=NodePort`

- Test that helm _will_ write it if you do by
  setting notePort to 30000:
  `helm template sentence-app/ --show-only templates/sentences-svc.yaml --set sentences.service.nodePort=30000,sentences.service.type=NodePort`

- Test that helm will not write a `nodePort`
  section, if the type is not NodePort:
  `helm template sentence-app/ --show-only templates/sentences-svc.yaml --set sentences.service.nodePort=30000,sentences.service.type=ClusterIP`

**Values Files**

When multiple variables needs to be set, its more
convenient to have them collected in files. Try
this by creating a file called
`values-resources.yaml` with the following
content:

```yaml
sentences:
  resources:
    requests:
      cpu: 0.25
    limits:
      cpu: 0.25
```

- Test template rendering with:

```shell
$ helm template sentence-app/ --values values-resources.yaml --show-only templates/sentences-deployment.yaml
```

- Validate the chart `helm lint sentence-app/`

- Install the sentences application using the new
  chart `helm install sentences sentence-app/`

This will install the chart with the default
values.

- Use `kubectl get pods` to see the running pods

- Try upgrading the chart using an increased
  replica count:
  `helm upgrade sentences sentence-app/ --set sentences.replicas=3`

- inspect the chart status and actual values:

```shell
$ helm list
$ helm get all sentences
$ helm get values sentences
```

</details>

## Cleanup

Delete the application installed with Helm:

```shell
$ helm delete sentences
```
