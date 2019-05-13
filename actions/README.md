# Container Diff for Github Actions

This is a Github Action to allow you to run Container Diff in a 
[Github Actions](https://help.github.com/articles/about-github-actions/#about-github-actions)
workflow. The intended use case is to build a Docker container from the repository,
push it to Docker Hub, and then use container-diff to extract metadata for it that
you can use in other workflows (such as deploying to Github pages). In
the example below, we will show you how to build a container, push
to Docker Hub, and then container diff.  Here is the entire workflow:

## Example 1: Run Container Diff

Given an existing container on Docker Hub, we can run container diff
without doing any kind of build.

```
workflow "Run container-diff isolated" {
  on = "push"
  resolves = ["list"]
}

action "Run container-diff" {
  uses = "GoogleContainerTools/container-diff/actions@master"
  args = ["analyze vanessa/salad --type=file --output=/github/workspace/data.json --json"]
}

action "list" {
  needs = ["Run container-diff"]
  uses = "actions/bin/sh@master"
  runs = "ls"
  args = ["/github/workspace"]
}
```

In the above, we run container-diff to output apt and pip packages, history,
and the filesystem for the container "vanessa/salad" that already exists on
Docker Hub. We save the result to a data.json output file. The final step in 
the workflow (list) is a courtesy to show that the data.json file is generated.

## Example 2: Build, Deploy, Run Container Diff

This next example is slightly more complicated in that it will run container-diff
after a container is built and deployed from a Dockerfile present in the repository.

```
workflow "Run container-diff after deploy" {
  on = "push"
  resolves = ["Run container-diff"]
}

action "build" {
  uses = "actions/docker/cli@master"
  args = "build -t vanessa/salad ."
}

action "login" {
  uses = "actions/docker/login@master"
  secrets = ["DOCKER_USERNAME", "DOCKER_PASSWORD"]
}

action "push" {
  uses = "actions/docker/cli@master"
  args = "push vanessa/salad"
}

action "Run container-diff" {
  needs = ["build", "login", "push"]
  uses = "GoogleContainerTools/container-diff/actions@master"
  args = ["analyze vanessa/salad --type=file --output=/github/workspace/data.json --json"]
}

action "list" {
  needs = ["Run container-diff"]
  uses = "actions/bin/sh@master"
  runs = "ls"
  args = ["/github/workspace"]
}
```

The intended use case of the above would be to, whenever you update your
container, deploy its metadata to Github pages (or elsewhere).
