# Welcome to heimdall's contributing guide

Thank you for investing your time in contributing to this project! Any contribution you will make, will make heimdall better for everyone :sparkles:!

Before you continue, here are some important resources:

* To keep the community around heimdall approachable and respectable, please read our [Code of Conduct](./CODE_OF_CONDUCT.md).
* [heimdall's Discussions](https://github.com/dadrus/heimdall/discussions) can help you in getting answers to your questions


## How can I contribute?

There are many ways you can contribute to heimdall. Here are some ideas:

* **Give this project a star**: It may not seem like much, but it really makes a difference. This is something that everyone can do to help. Github stars help the project gaining visibility and stand out.
* **Join the community**: Helping people can be as easy as by just sharing your own experience. You can also help by listening to issues and ideas of other people and offering a different perspective, or providing some related information that might help. Take a look at [heimdall's Discussions](https://github.com/dadrus/heimdall/discussions).  Bonus: You get GitHub achievements for answered discussions :wink:.
* **Review documentation**: Most [documentation](https://dadrus.github.io/heimdall/) just needs a review for proper spelling and grammar. If you think a document can be improved in any way, feel free to do so by opening a PR.
* **Help with open issues**: There are many [open issues](https://github.com/dadrus/heimdall/issues). Some of them may lack necessary information, some may be duplicates of older issues. Most are waiting for being implemented.
* **You spot a problem**: Search if an [issue already exists](https://github.com/dadrus/heimdall/issues). If a related issue doesn't exist, please open a new issue using a relevant [issue form](https://github.com/dadrus/heimdall/issues/new/choose).

## Disclosing vulnerabilities

Please disclose vulnerabilities exclusively to [dadrus@gmx.de](mailto:dadrus@gmx.de). Do not use GitHub issues.

## Contribute code

Unless you are fixing a known bug, we strongly recommend discussing it with the core team via a GitHub issue or in [heimdall's Discussions](https://github.com/dadrus/heimdall/discussions) before getting started.

The general process is as follows:

Set up your local development environment to contribute to heimdall:

1. [Fork](https://github.com/dadrus/heimdall/fork), then clone the repository.
  
   ```bash
   > git clone https://github.com/your_github_username/heimdall.git
   > cd heimdall
   > git remote add upstream https://github.com/dadrus/heimdall.git
   > git fetch upstream
   ```

2. Install required tools:
  * [Just](https://github.com/casey/just/releases), which is used to lint, build and test the project via CLI.
  * Recent [Golang](https://go.dev/dl/) version.
  * [golangci-lint](https://golangci-lint.run/usage/install/#local-installation) to lint the code.
  * [hadolint](https://github.com/hadolint/hadolint/releases) to lint the Dockerfiles.
  * [Helm](https://github.com/helm/helm/releases) to lint and work with Helm Chart provided with heimdall.
  * [go-licenses](https://github.com/google/go-licenses) to ensure all dependencies have allowed licenses.
  * [kubeconform](https://github.com/yannh/kubeconform/releases) to validate the manifests generated by Helm.
  * [Docker](https://docs.docker.com/desktop/) if you want to build and work with docker images.

3. Verify that tests and other checks pass locally.
   ```bash
   > git pull
   > git checkout main
   > just lint
   > just test
   ```

Make your changes:

1. Create a new feature branch.
   ```bash
   > git checkout -b cool_new_feature
   ```

2. Make your changes, and verify that all tests and lints still pass.
   ```bash
   > just lint
   > just test
   ```

3. When you're satisfied with the change, push it to your fork and make a pull request.
   ```bash
   > git push origin cool_new_feature
   # Open a PR at https://github.com/dadrus/heimdall/compare
   ```

When creating the PR, please follow the guide, you'll see in the PR template to make the review process more smoothly. 

At this point, you're waiting on us to review your changes. We *try* to respond to issues and pull requests within a few days, and we may suggest some improvements or alternatives. Once your changes are approved, one of the project maintainers will merge them.

## Contribute documentation

TODO

## Contribute examples

TODO