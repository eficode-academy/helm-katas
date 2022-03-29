# Testing Helm Deployments

## Learning Goals

- Write Helm deployment tests
- Test deployments with the `helm test` command

## Introduction

This exercise will deal with two kinds of testing helm charts.

Use helm to verify that the injected values in your chart are correct, and testing that deployments work as expected.

Helm provides functionality for orchestrating tests, using the `helm test` command.

Tests are defined as a number of `pod specs` which include the `test annotation` in their metadata dictionary.

[Helm Documentation](https://helm.sh/docs/topics/chart_tests/)

## Writing Helm Tests

A helm test consists of a `pod spec` with a specific annotation: `helm.sh/hook: test`.
As the annotation implies, this is actually a `helm hook`, meaning that all of the pod specs with this annotation are executed when the `helm test` command is issued.

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

<details>
<summary>:bulb: tip on test pod placement</summary>

> Test pod specs can be located anywhere in the in the `<chart>/templates` directory, though it is convention to place tests in a separate directory called `tests`, eg. `<chart>/templates/tests`.

</details>

[Further Reading](https://helm.sh/docs/topics/chart_tests/)

<details>
<summary>Detailed Helm test template</summary>

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

> :bulb: Note that we set the `restartPolicy` to `Never`.
> If we do not specify a restart policy, Kubernetes will try to be helpful, and will keep restarting our test pods, which will eventually fail the test once it reaches it's timeout.
> Therefore make sure to specify the `restartPolicy`.

You can use all of the functionality of normal pod specs when writing tests.

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
      image: ghcr.io/eficode-academy/network-multitool:latest
      command: ["curl", "-s", "sentence:8080"]
```

### command vs. args

When writing helm tests, you are likely to want to override the original `ENTRYPOINT` or `CMD` defined in the Dockerfile of the image used in the test.
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

### Helm Hooks / Automatically Removing Test Pods

You can use Helm hooks in your test pod specs to do useful things.
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

<details>
<summary>:bulb: (Don't put) Multiple Test-Containers in the same Pod</summary>

> Best practice when writing helm tests is to have each test container in it's own pod, but you can technically add as many containers to your test pods as you want.
> Having multiple containers in the same pod, will mean that the pod will only succeed if all of the containers exit successfully, and the pod will fail if just one of the containers exit unsuccessfully.
> This can be a useful pattern in certain cases, but you should know that if you do so, the `helm test --logs` command will not work, as helm will not know which of the containers in the pod to get logs from, and it will be up to you to gather the logs some other way.
> Therefore best practice is to put each test into it's own pod, such that all test logs can be viewed easily.

</details>

</details>

## Executing Helm Tests

Helm tests are executed with the `helm test` command followed by the name of the release name and optional flags.

```sh
$ helm test [RELEASE] [flags]
```

This will run all tests specified in the helm chart.

[Further Reading](https://helm.sh/docs/helm/helm_test/)

<details>
<summary>More Details</summary>

### Viewing Test Logs

The `stdout` and `stderr` of the test pods can be conveniently viewed when running tests by using the `--logs` flag on the test command.

```sh
$ helm test --logs [RELEASE]
```

The above command will run all of the tests and print the logs of each of the tests.

### Waiting for all Chart Resources to be Ready

If you are testing a newly deployed helm release, you might end up with errors because the release has not been completely deployed yet.

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

To do the exercises you can use the sentence helm chart you have already created.
Otherwise, the directory `test-helm-deployments/start` contains a clean starting point if you need it.

If you get stuck on any of the exercises, or want to see how the finished chart should look, you can look at the finished chart `test-helm-deployments/done`.

### Step by step instructions

<details>
<summary>More Details</summary>

**Add a helm test which checks that the sentence service is reachable**

- add a `tests` directory to the `templates` directory of your sentences helm chart:

```sh
$ mkdir sentence-app/templates/tests
```

- Create a file in the new tests directory called `sentence-svc-test.yaml`

- Add the following podspec to the file:

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
      image: ghcr.io/eficode-academy/network-multitool:latest
      command: ["curl", "-s", "sentence:8080"]
```

<details>
      <summary>:bulb: what does the podspec do?</summary>

> This helm test will run a pod with a single container, which will use the curl command to make a HTTP request to the sentence service.
> If the curl command receives a 200 response, then the container will exit with code 0, indicating a success.
> If the curl command does not receive a 200 response, the container will exit with a code that is greater than 0, indicating a failed test.
> Thus we can use this simple test to verify that after we have installed our chart, that our services are actually responding!

</details>

**Execute the sentence service test**

We have to deploy the test to the Kubernetes cluster, so that Kubernetes knows what to do when we issue the test command.

- Deploy (or upgrade) the existing deployment:

```sh
$ helm upgrade --install sentences sentence-app
Release "sentences" has been upgraded. Happy Helming!
NAME: sentences
LAST DEPLOYED: Wed Apr 28 08:42:36 2021
NAMESPACE: default
STATUS: deployed
REVISION: 2
```

- Verify that all resources are correctly deployed with `kubectl get pods`.

> It is important that all pods have the status `Running`, since otherwise we might get a false negative when we run the test.

- Execute the test: `$ helm test sentences`

- Verify that your output is successful like the below example:

```sh
$ helm test sentences
NAME: sentences
LAST DEPLOYED: Wed Apr 28 08:42:36 2021
NAMESPACE: default
STATUS: deployed
REVISION: 2
TEST SUITE:     sentences-sentence-svc-test
Last Started:   Wed Apr 28 08:42:41 2021
Last Completed: Wed Apr 28 08:42:45 2021
Phase:          Succeeded
```

- Use `kubectl` to list the pods, notice the test pod:

```sh
$ kubectl get pods
NAME                             READY   STATUS      RESTARTS   AGE
sentence-age-7c948b5d88-vrmbp    1/1     Running     0          3m27s
sentence-name-5687d74d64-mmhzs   1/1     Running     0          3m27s
sentences-668bd45d9-t5gn4        1/1     Running     0          3m27s
sentences-sentence-svc-test      0/1     Completed   0          2m58s
```

- Use `kubectl logs` to see the output of the test pod:

```sh
$ kubectl logs sentences-sentence-svc-test
Michael is 17 years
```

- Clean up the test pod:

```sh
$ kubectl delete pod sentences-sentence-svc-test
```

**Improve your test by using helm templating to make sure that values are injected correctly in the sentence service**

- Change the following lines in your sentence service template `templates/sentences-svc.yaml`:

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

- Change the test to use the same service name and port:

Change `templates/tests/sentence-svc-test.yaml` from:
```yaml
...
spec:
  ...
  containers:
    ...
      command: ["curl", "-s", "sentence:8080"]
```

To:

```yaml
...
spec:
  ...
  containers:
    ...
      command: ["curl", "-s", "{{ .Values.sentences.service.name }}:{{ .Values.sentences.service.port }}"]
```

> :bulb: now both places refers to the same value, meaning that both service and test will change when you change the value.

Next we add the service name and port values to the `values.yaml`.
Edit `sentence-app/values.yaml`, and add the `name: sentence` and `port: 9090` values under the sentence service:

```yaml
sentences:
  ...
  service:
    ...
    port: 9090
    name: sentence
```

This change enables us to template the service name and port that the sentence service will use.
The cool thing is that we can use the same templating in our test specification.
This is cool because we can use it to test that the service is actually using the values we have specified.

**Execute the improved sentence service test**

Upgrade the helm installation like you did before, and run the test the same way as before.

<details>
      <summary>:bulb: How did I do that?</summary>

> You can always go back and search the text for the commands we wanted you to perform. But a more direct way could be to use bash build-in history of all commands issued. To try it out, type `history` and a list of all commands you have issued will appear. Try to see if you can remember which ones you need to use.

</details>

The test should succeed.

- Clean up the test pod after the test has run with `kubectl delete pod sentences-sentence-svc-test`.

**Add a new test which uses regex to check that the returned body of the sentence service is correct**

Helm test pod specs can contain any container executing arbitrary commands.

Therefore we can create containers with custom code to test our deployments.

For this test we will use `regex` to test that the body returned form the sentences application is valid.
Which means that we can test that the deployment is not only responding, but that it is responding correctly.

We have prepared a small golang program that will query the endpoint and check the response using regex.

The program has already been packaged in a [docker image](https://hub.docker.com/r/releasepraqma/sentence-regex-test) so that we can use it a test spec.

<details>
      <summary>More details about the regex-tester</summary>

The sentence application returns a response that looks like this:

`Terry is 89 years`

We can break that into a pattern with four sections: a capitalized name, the word 'is', a number and finally the word 'years'.

We can create a regex statement to match this:

```regex
^[A-Z][a-z]+\ is\ \d+\ years$
```

If you are not sure how regex works, then don't worry, the important part is that this statement will verify that a response from the service follows the pattern outlined above.

We could verify the regex using shell commands, but that can get messy and hard to maintain, so let's use a programming language to write our test in.

The golang code is located in `test-helm-deployments/sentence-regex-test/sentence_regex.go`, but the implementation is not important for the purpose of this exercise.
The program will return an exit code 0 if the regex matches, and 1 if it does not.

</details>

We add a new test spec:

- Create a new file: `sentence-app/templates/tests/sentence-regex-test.yaml`

- Add the code:

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: "{{ .Release.Name }}-sentence-regex-test"
  annotations:
    "helm.sh/hook": test
spec:
  restartPolicy: Never
  containers:
    - name: "{{ .Release.Name }}-sentence-regex-test"
      image: releasepraqma/sentence-regex-test:latest
      args: ["http://{{ .Values.sentences.service.name }}:{{ .Values.sentences.service.port }}"]
```

The above pod spec should look familiar.
What is interesting to note is that it uses the image with the regex golang test, and that it takes the templated endpoint as it's argument.

**Execute both tests**

- Upgrade the Helm chart to install the new test.

- Execute the test command.

This time Helm will execute both of our tests sequentially:

```sh
$ helm test sentences
NAME: sentences
LAST DEPLOYED: Wed Apr 28 09:30:59 2021
NAMESPACE: default
STATUS: deployed
REVISION: 3
TEST SUITE:     sentences-sentence-regex-test
Last Started:   Wed Apr 28 09:31:08 2021
Last Completed: Wed Apr 28 09:31:13 2021
Phase:          Succeeded
TEST SUITE:     sentences-sentence-svc-test
Last Started:   Wed Apr 28 09:31:13 2021
Last Completed: Wed Apr 28 09:31:14 2021
Phase:          Succeeded
```

- Verify the logs from the regex test pod:

```sh
$ kubectl logs sentences-sentence-regex-test
2021/04/28 07:31:13 response: ' Michael is 13 years ' is valid.
```

> :bulb: You can add as many tests as you need to your helm chart, and the `test` command will execute all of them.

</details>

### Clean up

To clean up the release together with the test pods:

- `helm uninstall sentences`
- `kubectl delete pod sentences-sentence-regex-test`
- `kubectl delete pod sentences-sentence-svc-test`
