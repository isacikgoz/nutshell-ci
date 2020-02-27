
# CI-fail finder

A toy tool to parse output of failed tests in [mattermost-server](https://github.com/mattermost/mattermost-server/). Requires `GITHUB_TOKEN` and `CIRCLECI_TOKEN` env vars to be set.

install:

```bash
go get github.com/isacikgoz/nutshell-ci/cmd/nutshell-ci
```

## Run

```bash
nutshell-ci https://github.com/mattermost/mattermost-server/pull/13045
```

Also you can pipe regular `go test -v` output

```bash
make test | nutshell-ci
```
