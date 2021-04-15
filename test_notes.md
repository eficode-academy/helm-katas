# Testing helm Charts

Chart tests validate that a deployed helm chart is working as expected.

Tests live in the `/templates` directory and are kubernetes `job` definitions.
Each job definition runs a container and specific command in that container.
The command in the container must exit successfully with the exit code 0, for the job to be successfully and the test to be considered passed.
Each job definition must contain the helm test annotation: `helm.sh/hook: test`.

## Examples of helm tests

- Test that values from `values.yaml` were successfully injected into the rendered yaml.
- Test that credentials work as expected.
- Test that invalid credentials do not work.
- Test that services are correctly serving traffic on the endpoints expected.
- Test that services are correctly load balanced.

## Writing Tests

Docs: https://helm.sh/docs/topics/chart_tests/

...

## Running Tests

Docs: https://helm.sh/docs/helm/helm_test/

Tests are executed with the `helm test <release-name>` command:

```sh
helm test [RELEASE] [flags]
```

Tests should be run once a `helm install` command has been executed, and all of the components have been deployed, as helm will not wait for this.

Each test is a `pod` job spec, in one or several `.yaml` files.

Test spec files should be located in the `/templates` directory, you may put them in a sub directory like `/templates/tests` for a clean code structure.

A helm test is actually a `helm hook` (https://helm.sh/docs/topics/charts_hooks/), so you can use other annotations in conjunction with test resources for more advanced test behaviour.
