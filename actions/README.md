# Container Diff for Github Actions

This is a Github Action to allow you to run Container Diff in a 
[Github Actions](https://help.github.com/articles/about-github-actions/#about-github-actions)
workflow. The intended use case is to build a Docker container from the repository,
push it to Docker Hub, and then use container-diff to extract metadata for it that
you can use in other workflows (such as deploying to Github pages). You can also run
container diff to extract metadata for a container you've just built locally in the action.

## 1. Action Parameters

The action accepts the following parameters:

| Name | Description | Type| Default | Required |
|------|-------------|-----|---------|----------|
| command | main command for container-diff | string | analyze | false |
| args  | The full list of arguments to follow container-diff (see example below) | string | help | true |

See below for a simple example. Another interesting use case would be to generate metadata and upload
to an OCI registry using [OCI Registry As Storage](https://oras.land/).

## 2. Run Container Diff

Given an existing container on Docker Hub, we can run container diff
without doing any kind of build.

```yaml
name: Run container-diff

on: 
  pull_request: []

jobs:
  container-diff:
    name: Run container-diff
    steps:
      - name: Checkout
        uses: actions/checkout@v4
      - name: Run container-diff
        uses: GoogleContainerTools/container-diff/actions@master
        with:
          # Note this command is the default and does not need to be included
          command: analyze          
          args: vanessa/salad --type=file --output=./data.json --json
      - name: View output
        run: cat ./data.json
```

In the above, we run container-diff to output apt and pip packages, history,
and the filesystem for the container "vanessa/salad" that already exists on
Docker Hub. We save the result to a data.json output file. The final step in 
the workflow (list) is a courtesy to show that the data.json file is generated.
