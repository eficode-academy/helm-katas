# Testing Helm Deployments

## Learning Goal

- Write helm deployment tests
- Test deployments with the `helm test` command

## Introduction

There are multiple aspects of a helm chart that one might want to test, such as if the rendered Kubernetes yaml is correct, whether values are injected correctly into the rendered yaml, and that when a chart is deployed, that the services started are working as expected.
This Exercise will deal will deal with the two latter, verifying that the injected values are correct, and testing that deployments work as expected.

Helm provides functionality for orchestrating tests, using the `helm test` command.
Tests are defined as a number of `pod specs` which include the `test annotation` in their metadata dictionary.

## Writing Helm Tests

A helm test consists of a `pod spec` with a specific annotation: `helm.sh/hook: test`.
As the annotation implies, this is actually a `helm hook`, meaning that all of the pod specs with this annotation are executed when the helm test command is issued.

The hook is an annotation to the metadata of the pod:

```yaml
...
kind: Pod
metadata:
  ...
  annotations:
    # this annotation is what makes this pod spec a helm test!
    "helm.sh/hook": test
...
```

Test pod specs can be located anywhere in the in the `<chart>/templates` directory, though it is convention to place tests in a separate directory called `tests`, eg. `<chart>/templates/tests`.

[Further Reading](https://helm.sh/docs/topics/chart_tests/)

<details>
<summary>More Details</summary>

##### Helm Test Template

Below is an example of a complete boilerplate test pod spec:

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: "{{ .Release.Name }}-example-test"
  annotations:
    # this annotation is what makes this pod spec a helm test!
    "helm.sh/hook": test
spec:
  restartPolicy: Never
  containers:
    - name: "{{ .Release.Name }}-example-test"
      image: <container-image>:<tag>
      command: ["example-command", "example-argument"]
```

**Note** that we set the `restartPolicy` to `Never`.
If we do not specify a restart policy, Kubernetes will try to be helpful, and will keep restarting our test pods, which will eventually fail the test once it reaches it's timeout.
Therefore make sure to specify the `restartPolicy`.

You can of course use all of the functionality of normal pod specs when writing tests, as well as use variable injection to template the tests themselves.

Here is an example test that will check if the http endpoint of the sentences application responds to requests:

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: "{{ .Release.Name }}-sentence-svc-test"
  annotations:
    # this annotation is what makes this pod spec a helm test!
    "helm.sh/hook": test
spec:
  restartPolicy: Never
  containers:
    - name: "{{ .Release.Name }}-sentence-svc-test"
      image: praqma/network-multitool:minimal
      command: ["curl", "-s", "sentence:8080"]
```

##### command vs. args

When writing helm tests, you are likely to want to override the origin `ENTRYPOINT` or `CMD` defined in the Dockerfile of the image used in the test.
In kubernetes this is done, slightly unintuitively, by using the `command` key of the container spec to define the `ENTRYPOINT`, and the `args` key to define the `CMD` of the container.

An example of overwriting the entrypoint of container:
```yaml
spec:
  containers:
      ...
      command: ["curl", "-s", "sentence:8080"]
```

An example of overwriting both the entrypoint (with `command`) and the cmd (with `args`)
```yaml
spec:
  containers:
      ...
      command: ["curl"]
      args: ["-s", "sentence:8080"]
```
You can of course also use `args` by itself without modifying the `command`.

[Further Reading](https://kubernetes.io/docs/tasks/inject-data-application/define-command-argument-container/)

##### Helm Hooks / Automatically Removing Test Pods

You can use helm hooks in your test pod specs to do useful things.
An example could be to delete pods after they have completed successfully.
This is done with the `helms.sh/hook-delete-policy: hook-succeeded` hook.
The new hook is added to the annotations of the pod spec:

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: "{{ .Release.Name }}-example-test"
  annotations:
    # this annotation is what makes this pod spec a helm test!
    "helm.sh/hook": test
    # we can use this hook to automatically delete the pod
    # after the test has succussfully run, but keep the pod if it fails
    # such that we can inspect why it failed
    "helm.sh/hook-delete-policy": hook-succeeded
```

> Note: With the current version of helm, v3.5.4, when using this hook, pods are deleted immediately, which means that the `helm test --logs <release>` will not print the logs, as the pods are deleted too early.
> This is a [known issue](https://github.com/helm/helm/issues/9098) and will hopefully soon be fixed.


[Further Reading](https://helm.sh/docs/topics/charts_hooks/)

##### (Don't Put) Multiple Tests in the same Pod

Best practice when writing helm tests is to have each test container in it's own pod, but you can technically add as many containers to your test pods as you want.
Having multiple containers in the same pod, will mean that the pod will only succeed if all of the containers exit successfully, and the pod will fail if just one of the containers exit unsuccessfully.
Do note that if you do so, the `helm test --logs` command will not work, as helm will not know which of the containers in the pod to get logs from.
Therefore best practice is to put each test into it's own pod, such that all test logs can be viewed easily.

</details>

## Executing Helm Tests


Helm tests are executed with the `helm test` command followed by the name of the release name and optional flags.
```sh
$ helm test [RELEASE] [flags]
```
This will run all test specified in the helm chart.

[Further Reading](https://helm.sh/docs/helm/helm_test/)

<details>
<summary>More Details</summary>

##### Viewing Test Logs
The stdout/stderr of the test pods can be conveniently viewed when running tests by using the `--logs` flag on the test command.
```sh
$ helm test --logs [RELEASE]
```
The above command will run all of the tests and print the logs of each of the tests.

##### Waiting for all Chart Resources to be Ready
When deploying a helm chart that has tests, one will usually execute the tests each time a deployment of the chart is done.
One has to note though that executing a `$ helm test` immediately after a `$ helm install` might produce false failed tests, as the helm test command does not check if all of the chart resources are ready.
To alleviate this we can use the `--wait` flag on the install command to make helm wait for all of the chart resources to be ready before moving to the next command.
```sh
$ helm install --wait [RELEASE] [CHART] && helm test [RELEASE]
```
Hence the above command would first install the chart, then wait for all of the chart resources to be in the ready state, and then run the tests.


</details>

## Exercise

### Overview

- Add a helm test which checks that the sentence service is reachable
- Execute the sentence service test
- Improve your test by using helm templating to check that values are injected correctly in the sentence service
- Execute the improved sentence service test
- Add a new test which uses regex to check that the returned body of the sentence service is correct
- Execute both tests

### Step by Step

<!-- <details> -->
<!-- <summary>More Details</summary> -->

##### * Add a helm test which checks that the sentence service is reachable

Start by adding a `tests` directory to the templates of your sentences helm chart:
```sh
mkdir sentence-app/templates/tests
```
(assuming that your helm chart is called `sentence-app`)

Create a file in the new tests directory called `sentence-svc-test.yaml`, and add the following code to the file:
```yaml
apiVersion: v1
kind: Pod
metadata:
  name: "{{ .Release.Name }}-sentence-svc-test"
  annotations:
    # this annotation is what makes this pod spec a helm test!
    "helm.sh/hook": test
spec:
  restartPolicy: Never
  containers:
    - name: "{{ .Release.Name }}-sentence-svc-test"
      image: praqma/network-multitool:minimal
      command: ["curl", "-s", "sentence:8080"]
```

This helm test will run a pod, the pod will run a single container, which will curl the sentence service.
If the curl command recieved a 200 response, then the container will exit with code 0, indicating a success.
If the curl command does not recieve a 200 response, the container will exit with a code that is greater than 0, indicating a failed test.

Thus we can use this simple test to verify that after we have installed our chart, that our servies are actually responding!

##### * Execute the sentence service test

Now let's execute the new test that we have created.
First we have to deploy the test to the Kubernetes cluster, we can do that by upgrading the existing deployment.
```sh
$ helm upgrade sentences sentence-app
Release "sentences" has been upgraded. Happy Helming!
NAME: sentences
LAST DEPLOYED: Thu Apr 22 14:49:03 2021
NAMESPACE: default
STATUS: deployed
REVISION: 2
```

Verify that all resources are correctly deployed with `kubectl get`.

It is important that all pods are in the `ready` state, since otherwise we might get a false negative when we run the test.

Now execute the test:
```sh
$ helm test sentences
NAME: sentences
LAST DEPLOYED: Thu Apr 22 14:49:03 2021
NAMESPACE: default
STATUS: deployed
REVISION: 2
TEST SUITE:     sentences-sentence-svc-test
Last Started:   Thu Apr 22 14:49:33 2021
Last Completed: Thu Apr 22 14:49:38 2021
Phase:          Succeeded
```

As we can see from the output, the test executed successfully.

We can inspect the test pod:
```sh
$ kubectl get pods
NAME                             READY   STATUS      RESTARTS   AGE
sentence-age-7c948b5d88-wwhj7    1/1     Running     0          6m21s
sentence-name-5687d74d64-49skg   1/1     Running     0          6m21s
sentences-668bd45d9-cmx5v        1/1     Running     0          6m21s
sentences-sentence-svc-test      0/1     Completed   0          3m25s

$ kubectl logs sentences-sentence-svc-test
Michael is 17 years
```

Clean up the test pod:
```sh
$ kubectl delete pod sentences-sentece-svc-test
```

##### * Improve your test by using helm templating to check that values are injected correctly in the sentence service

Now let's test the templating functionality of helm.

Change the following lines in your sentence service template `templates/sentences-svc.yaml`:

From:
```yaml
...
metadata:
  ...
  name: sentence
spec:
  ports:
    - port: 8080
      ...
  ...
```
To:
```yaml
...
metadata:
  ...
  name: {{ .Values.sentences.service.name }}
spec:
  ports:
    - port: {{ .Values.sentences.service.port }}
      ...
  ...
```



```yaml
```

```yaml
```
##### * Execute the improved sentence service test
##### * Add a new test which uses regex to check that the returned body of the sentence service is correct
##### * Execute both tests

<!-- </details> -->


## Going Further

testing credentials -> bitnami mysql chart
test other aspects of helm
test that rendered yaml is correct
lint
