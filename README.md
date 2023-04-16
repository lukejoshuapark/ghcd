# ghcd

_"GitHub Change Detection_"

A CLI tool that enables change detection in GitHub workflows.

## Modes

There are two modes offered by `ghcd`:

- `FilesDiff` - This mode pulls a list of all files modified between two
commits.  Intended for CI e.g. building and testing.

- `EnvironmentDiff` - This mode compares all files changed between the commit
that is currently live in a GitHub environment and the commit of the currently
running GitHub workflow.  Intended for CD e.g. deployment.

## Configuration

Configuration is provided to `ghcd` as a YAML file.  An example `ghcd.yml` file
that demonstrates all features is provided below:

```yml
detect:

  build-api:
    mode: FilesDiff
    paths:
      - src/api

  deploy-api:
    mode: EnvironmentDiff
    environment: production-api
    paths:
      - src/api
      - deploy/terraform/api
```

Paths in the configuration are relative to the current working directory when
executing `ghcd`.

## CLI Reference

An example execution of `ghcd` in a GitHub workflow context is provided below.

```bash
ghcd detect -f ghcd.yml --token ${{ secrets.GITHUB_TOKEN }} --repository ${{ github.repository }} --start ${{ github.event.before }} --end ${{ github.event.after }}
```

There is currently only one command, `detect`, which accepts the following
flags.

|Flag|Required|Description|
|----|--------|-----------|
|f|No|The configuration file to use. If omitted, "ghcd.yml" is used.|
|token|If configuration uses `EnvironmentDiff` mode|The GitHub token to use to access the GitHub API.  Used for checking what commit is in an environment.|
|repository|If configuration uses `EnvironmentDiff` mode|The GitHub repository.|
|start|If configuration uses `FilesDiff` mode|The starting commit used for comparison in `FilesDiff` mode.  Changes made in the specified commit **are not** considered when checking for changes.|
|end|Yes|The end commit used for all comparisons.  Changes made in the specified commit **are** considered when checking for changes.|

## Output

It is intended for `ghcd` to be used to produce your job outputs.  Consider the
following segment of a GitHub workflow:

```yml
jobs:
  change-detection:
    runs-on: ubuntu-latest
    outputs:
      build-api: ${{ steps.ghcd.outputs.build-api }}
      deploy-api: ${{ steps.ghcd.outputs.deploy-api }}
    steps:
      - id: ghcd
        run: ghcd detect --token ${{ secrets.GITHUB_TOKEN }} --repository ${{ github.repository }} --start ${{ github.event.before }} --end ${{ github.event.after }} >> "$GITHUB_OUTPUT"
```

If we assume the example `ghcd.yml` file provided above is present in the root
of the repository, the `change-detection` job defined above will output
something like the following, including the new lines:

```bash
build-api=true
deploy-api=true

```

The value for each pair will differ based on whether or not any changes occured
in any of the provided paths, with respect to the operating mode.

For example, `build-api` will have `true` if any files changed between the
`start` and `end` commits had changes with that path prefix.

Further, `deploy-api` will have `true` if there have been any changes to files
with the provided path prefixes since the commit that was last deployed to
the `production-api` environment and the `end` commit.
