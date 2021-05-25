# Packaging and sharing charts

Note: This exercise requires that you have forked the GitHub repository to your own user, so you have administrative access.

## Learning goal

- Packaging Helm charts
- Use GitHub as a simple chart registry
- Share your chart through GitHub

## Introduction

> :bulb: Know that a Helm repository is nothing more than a simple HTTP server that can host Helm charts.
> Helm charts consists of a .tar file and an index.yaml file.
> As simple as that.

## Exercise

The chart we will work with in this exercise is located in the sharing charts folder.

### Overview

- Clone down the repository
- Package your chart
- Create index.html and push `gh-pages` branch to github
- Create index.yaml and push `gh-pages` branch to github
- Add the new repository to your Helm cli

### Step by step

<details>
      <summary>More details</summary>

**Clone down the repository**

- In the helm-katas folder, make a new folder with your github handle as your name.
- Open a terminal into that folder, and clone down your forked repository.

> :bulb: make sure that the repo you are cloning down is your own, and not the eficode-academy one. Yours will have a URL like the following: https://github.com/yourusername/helm-katas where `yourusername` is replaced with your username.

**Package your chart**

- Open a terminal in the `<yourusername>/helm-katas/sharing-charts` directory
- Package your chart with `helm package sentence-app`
- Move the package out to the root folder of your cloned repository

> the path would be something like this: `/home/ubuntu/helm-katas/YourGHName/helm-katas`

**Create index.html and push `gh-pages` branch to github**

- Add an empty index.html file in the root: `touch index.html`

> Note: this is for github to recognize this as a website and start serving it as content.

- create a branch named `gh-pages` from the main branch, and check it out.
- add the helm chart `.tgz` and the `index.html` to git, make a commit and push it to your new `gh-pages` branch.

> Note: the VSCode instances used would like to login for you with OAuth. We therefore recommend you to use a dummy github account for this due to security considerations.

<details>
      <summary>:bulb: git help</summary>

- Make sure you are in the folder `/home/ubuntu/helm-katas/YourGHName/helm-katas`
- Create and check out a new branch called `gh-pages` by running: `git checkout -b gh-pages`
- Type `git status` to see that you have two new files, your index file and the app in compressed format.
- Add both files to git with with `git add index.html` and `git add sentence-app-0.1.0.tgz`
- Commit both files with `git commit -m "made first gh page"`
- Push them to Github with `git push --set-upstream origin gh-pages`

</details>

- Go to the Settings tab of your Github repository and the `Pages` tab on the left. Here you will see a link, in the format like: https://UserName.github.io/helm-katas/.
- Click the link to see a blank webpage, making sure that the page is served through github.

**Create index.yaml and push `gh-pages` branch to github**

- Use the `helm repo index` command to generate the `index.yaml` file: `helm repo index . --url https://<yourGitHubUsername.github.io/helm-katas`
- Open the newly created file to see that the content matches the below example:

```yaml
apiVersion: v1
entries:
  sentence-app:
  - apiVersion: v2
    appVersion: 1.16.0
    created: "2021-05-18T11:05:16.935664314Z"
    description: A Helm chart for Kubernetes
    digest: 2125cd363e6f764472cb70c7afab5e35170c64ae06630bc7fd15577a40afaef4
    name: sentence-app
    type: application
    urls:
    - https://sofusalbertsen.github.io/helm-katas/sentence-app-0.1.0.tgz
    version: 0.1.0
generated: "2021-05-18T11:05:16.93437855Z"
```

- Add, commit and push the index file to the `gh-pages` branch.

Congratulations! You have now made your first chart repository.

**Add the new repository to your Helm cli**

To test out your newly created repo, try to add it to your helm CLI.

- Add the repository to Helm: `helm repo add my-repo https://<yourGitHubUsername.github.io/helm-katas`
- List your Helm repositories to see the newly added repo: `helm repo list`

```sh
$ helm repo list
NAME    URL
my-repo https://sofusalbertsen.github.io/helm-katas/
```

Great! Your helm repo is live and now you can fetch your helm chart or install your helm chart.

```sh
$ helm install sentence-app my-repo/sentence-app
NAME: sentence-app
LAST DEPLOYED: Tue May 18 11:23:54 2021
NAMESPACE: user1
STATUS: deployed
REVISION: 1
```

- Watch the Kubernetes object gets created with `kubectl get pods,svc`
- Clean up by uninstalling the chart: `helm uninstall sentence-app`

</details>

> :bulb: if you have multiple charts in the same repo added at different times, you can merge new versions into the same index.yaml file using `--merge` flag. For more info visit the [documentation](https://helm.sh/docs/topics/chart_repository/#add-new-charts-to-an-existing-repository)

> :bulb: there is a new way of sharing charts now; using Open Container Initiative format (OCI). In that way, your chart is saved in the same repository as your images. It is an experimental feature for now, but you can read up upon it (and instructions to try it out) in the [documentation](https://helm.sh/docs/topics/registries/#enabling-oci-support)

### Extra (optional)

This is the "Manual" way of doing a helm chart repo, and it has several downsides:

* It stores all versions of your charts in a packaged (binary) file in your git repo, creating a large repository to clone over time.
* It right now is manually done, so it needs to be CI'ed in a pipeline to become really usefull.

But there is another way, using the [releaser](https://helm.sh/docs/howto/chart_releaser_action/) tool.

The guide linked to describes how to use Chart Releaser Action to automate releasing charts through GitHub pages. Chart Releaser Action is a GitHub Action workflow to turn a GitHub project into a self-hosted Helm chart repo, using helm/chart-releaser CLI tool.

Have a look at how to set this up.

### Credits

This exercise has been adapted from a medium blogpost by [Ravindra Singh](https://medium.com/xebia-engineering/how-to-share-helm-chart-via-helm-repository-4cbfc7b1df90).
