<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [Contributing guidelines](#contributing-guidelines)
  - [Contributions](#contributions)
  - [Certificate of Origin](#certificate-of-origin)
  - [Issue and Pull Request Management](#issue-and-pull-request-management)
  - [Contributing A Patch](#contributing-a-patch)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# Contributing guidelines

## Contributions

All contributions to the repository must be submitted under the terms of the [Apache Public License 2.0](https://www.apache.org/licenses/LICENSE-2.0).

## Certificate of Origin

By contributing to this project you agree to the Developer Certificate of
Origin (DCO). This document was created by the Linux Kernel community and is a
simple statement that you, as a contributor, have the legal right to make the
contribution. See the [DCO](DCO) file for details.

## Issue and Pull Request Management

Anyone may comment on issues and submit reviews for pull requests. However, in
order to be assigned an issue or pull request, you must be a member of the
[open-cluster-management](https://github.com/open-cluster-management-io) GitHub organization.

Repo maintainers can assign you an issue or pull request by leaving a
`/assign <your Github ID>` comment on the issue or pull request.

## Contributing A Patch

1. Submit an issue describing your proposed change to the repo in question.
2. The [repo owners](OWNERS) will respond to your issue promptly.
3. Fork this repo and clone the forked repo to your `$GOPATH/src/open-cluster-management.io/api` directory.
4. After your code changes is ready to commit, please run following commands to check your code.
5. Check that the `GOPATH` environment variable has been set correctly. If not, this makefile will automatically get the value from `go env GOPATH` and set it.

   ```shell
   make update
   make verify
   ```

6. Submit a pull request.

Now, you can follow the [getting started guide](./README.md#getting-started) to work with the open-cluster-management API repository.
