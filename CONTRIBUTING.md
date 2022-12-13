# Contributing guidelines

## Formatting

By default, formating is checked by gofmt tool, where almost format all the
cases.

The project does not have a line length limit, but the team try to be as small
as possible, keeping the code clean on reading. Having an 80-120 char limit is
an excellent number to split the line.

## Style

The project follows the Golang project guidelines that you can find in the
following link:

[https://github.com/golang/go/wiki/CodeReviewComments]()

## Developer environment

By default, the project rely on Makefile, where users can run all the
workflows operations.

To make it easier, the Makefile implement a help section to see the actions you
can run. Here is an example:

```
make help
```

### Testing:

TBD

## Pull Request process

1) Fork the project and commit to a local branch.
2) Submit a PR with all details. Small PRs are prefered; on larger ones, please
ask any maintainers before a significant change to be aligned with the project
roadmap. DCO is needed.
3) The PR to be approved should contain test cases on the new features added.
4) Maintainer will approve the GH actions checks.
5) If all checks are working, PR will be merged.  (Checks can be found on
`.github` folder)

### Contributor compliance with Developer Certificate Of Origin (DCO)

We require every contributor to certify that they are legally permitted to
contribute to our project.  A contributor expresses this by consciously signing
their commits, and by this act expressing that they comply with the [Developer
Certificate Of Origin](https://developercertificate.org/).

A signed commit is a commit where the commit message contains the following
content:

```
Signed-off-by: John Doe <jdoe@example.org>
```

This can be done by adding
[`--signoff`](https://git-scm.com/docs/git-commit#Documentation/git-commit.txt---signoff)
to your git command line.

## Documentation

TBD
