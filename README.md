# GitHub Repo Labeller

Label all of your repositories within an organization with labels defined in
`labels.json` from the current directory.

This is a small tool to consistently label **all** repositories in a Github org
with the specified labels.  This takes advantage of the [go-github
library](https://github.com/google/go-github).

While [Github labels can be set at the organization
level](https://docs.github.com/en/github/setting-up-and-managing-organizations-and-teams/managing-default-labels-for-repositories-in-your-organization),
existing repositories are not updated to reflect those labels.

## Install

```sh
GO111MODULE=off go get github.com/displague/github-labeller
```

## Usage

**First**, make sure `labels.json` exists and has all of the labels you want to
add to your repositories.

```console
GITHUB_TOKEN=... github-labeller ORG
```

## TODO

* [ ] Add the labels to the Org (if permitted)?
* [ ] Update existing labels?
