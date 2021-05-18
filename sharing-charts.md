# Packaging and sharing charts

Note: This exercise requires that you have forked the GitHub repository to your own user, so you have administrative access.

## Learning goal

- Be able to package your chart
- Use GitHub as a simple chart registry
- Share your chart through GitHub

## Introduction

> :bulb: Know that a Helm repository is nothing more than a simple HTTP server that can host Helm charts .tar file and an index.yaml file. As simple as that.

## Exercise

The chart we will work with in this exercise is located in the sharing charts folder.

### Overview

- clone down the repository
- package your chart
- create index.html and push `gh-pages` branch to github
- create index.yaml and push `gh-pages` branch to github
- add the new repo to your helm cli

### Step by step

<details>
      <summary>More details</summary>

**clone down the repository**

- In the helm-katas folder, make a new folder with your github handle as your name.
- Open a terminal into that folder, and clone down your forked repository.

> :bulb: make sure that the repo you are cloning down is your own, and not the eficode-academy one. Yours will have a URL like the following: https://github.com/yourusername/helm-katas where `yourusername` is replaced with your username.

**Package your chart**

- Package your chart with `helm package sentence-app`
- Move the package out to the root folder of your cloned repo

> :bulb: the path would be something like this: `/home/ubuntu/helm-katas/YourGHName/helm-katas`

**create index.html and push `gh-pages` branch to github**

- Add an empty index.html file in the root: `touch index.html`

> :bulb: this is for github to recognize this as a website and start serving it as content.

- Create a branch named `gh-pages` from the main branch, and check it out.
- Commit and push your new branch

<details>
      <summary>git help</summary>

- make sure you are in the folder `/home/ubuntu/helm-katas/YourGHName/helm-katas`
- create and check out a new branch called `gh-pages` by running: `git checkout -b gh-pages`
- type `git status` to see that you have two new files, your index file and the app in compressed format.
- add both files to git with with `git add index.html` and `git add sentence-app-0.1.0.tgz`
- commit both files with `git commit -m "made first gh page"`
- push them to Github with `git push --set-upstream origin gh-pages`

</details>

- go to 

</details>

> :bulb: if you have multiple charts in the same repo added at different times, you can merge new versions into the same index.yaml file using `--merge` flag. For more info visit the [documentation](https://helm.sh/docs/topics/chart_repository/#add-new-charts-to-an-existing-repository)

> :bulb: there is a new way of sharing charts now; using Open Container Initiative format (OCI). In that way, your chart is saved in the same repository as your images. It is an experimental feature for now, but you can read up upon it (and instructions to try it out) in the [documentation](https://helm.sh/docs/topics/registries/#enabling-oci-support)

> :bulb: this exercise has been adapted from this medium blogpost by [Ravindra Singh](https://medium.com/xebia-engineering/how-to-share-helm-chart-via-helm-repository-4cbfc7b1df90).