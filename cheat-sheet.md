
# List of Helm Commands

Use the commands listed below as a quick reference when working with Helm inside Kubernetes.

## Install and Uninstall Apps

Install an app:

`helm install [name] [chart]`

Install an app in a specific namespace:

`helm install [name] [chart] --namespace [namespace]`

Override the default values with those specified in a file of your choice:

`helm install [name] [chart] --values [yaml-file/url]`

Run a test install to validate and verify the chart:

`helm install [name] --dry-run --debug`

Uninstall a release:

`helm uninstall [release]`

## Perform App Upgrade and Rollback

Helm offers users multiple options for app upgrades, such as automatic rollback and upgrading to a specific version. Rollbacks can also be executed on their own.

Upgrade an app:

`helm upgrade [release] [chart]`

Instruct Helm to rollback changes if the upgrade fails:

`helm upgrade [release] [chart] --atomic`

Upgrade a release. If it does not exist on the system, install it:

`helm upgrade [release] [chart] --install`

Upgrade to a specified version:

`helm upgrade [release] [chart] --version [version-number]`

Roll back a release:

`helm rollback [release] [revision]`

## Download Release Information

The helm get command lets you download information about a release.

Download all the release information:

`helm get all [release]`

Download all hooks:

`helm get hooks [release]`

Download the manifest:

`helm get manifest [release]`

Download the notes:

`helm get notes [release]`

Download the values file:

`helm get values [release]`

Fetch release history:

`helm history [release] `

## Add, Remove, and Update Repositories

The helm command helm repo helps you manipulate chart repositories.

Add a repository from the internet:

`helm repo add [name] [url]`

Remove a repository from your system:

`helm repo remove [name]`

Update repositories:

`helm repo update`

## List and Search Repositories

Use the helm repo and helm search commands to list and search Helm repositories. Helm search also enables you to find apps and repositories in Artifact hub.

List chart repositories:

`helm repo list`

Generate an index file containing charts found in the current directory:

`helm repo index`

Search charts for a keyword:

`helm search [keyword]`

Search repositories for a keyword:

`helm search repo [keyword]`

Search Helm Hub:

`helm search hub [keyword]`

## Release Monitoring

The helm list command enables listing releases in a Kubernetes cluster according to several criteria, including using regular (Pearl compatible) expressions to filter results. Commands such as helm status and helm history provide more details about releases.

List all available releases in the current namespace:

`helm list`

List all available releases across all namespaces:

`helm list --all-namespaces`

List all releases in a specific namespace:

`helm list --namespace [namespace]`

List all releases in a specified output format:

`helm list --output [format]`

Apply a filter to the list of releases using regular expressions:

`helm list --filter '[expression]'`

See the status of a specific release:

`helm status [release]`

Display the release history:

`helm history [release]`

See information about the Helm client environment:

`helm env`

## Chart Management

Helm charts use Kubernetes resources to define an application.

Create a directory containing the common chart files and directories (Chart.yaml, values.yaml, charts/ and templates/):

`helm create [name]`

Package a chart into a chart archive:

`helm package [chart-path]`

Run tests to examine a chart and identify possible issues:

`helm lint [chart]`

Inspect a chart and list its contents:

`helm show all [chart]`

Display the chart’s definition:

`helm show chart [chart]`

Display the chart’s values:

`helm show values [chart]`

Download a chart:

`helm pull [chart]`

Download a chart and extract the archive’s contents into a directory:

`helm pull [chart] --untar --untardir [directory]`

Display a list of a chart’s dependencies:

`helm dependency list [chart]`