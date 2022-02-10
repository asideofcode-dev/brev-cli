<p align="center">
<img width="230" src="https://raw.githubusercontent.com/brevdev/assets/main/logo.svg"/>
</p>

# Brev.dev

[Brev.dev](https://brev.dev) connects your local computer to a cloud mesh so you can code locally on remote machines, without needing any new tools, and never in a browser.

## Install

if brew is not already installed on your computer install it with
```
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
```

then add Brev's tap and install `brev` with

```
brew install brevdev/homebrew-brev/brev
```

# Usage


create a workspace in your current org
```
brev start https://github.com/brevdev/brev-cli
```

list all workspaces in an org
```
brev ls
```

stop a workspace
```
brev stop brevdev/brev-cli
```

delete a workspace from an org
```
brev delete brevdev/brev-cli
```

## Completion


### zsh


```
mkdir -p ~/.zsh/completions && brev completion zsh > ~/.zsh/completions/_brev && echo fpath=~/.zsh/completions $fpath >> ~/.zshrc && fpath=(~/.zsh/completions $fpath) && autoload -U compinit && compinit
```

### bash

```
sudo mkdir -p /usr/local/share/bash-completion/completions
brev completion bash | sudo tee /usr/local/share/bash-completion/completions/brev
source /usr/local/share/bash-completion/completions/brev
```

### fish

```
mkdir -p ~/.config/fish/completions && brev completion fish > ~/.config/fish/completions/brev.fish && autoload -U compinit && compinit
```

# Development

## Build

`make build` runs a full release build
`make fast-build` builds a binary for your current machine only

## example .env file

```
VERSION=unknown
BREV_API_URL=http://localhost:8080
# BREV_API_URL=https://ade5dtvtaa.execute-api.us-east-1.amazonaws.com
```


## adding new commands

`pkg/cmd/logout/logout.go` is a minimal command to go off of for adding new commands.

commands for the cli should follow `<VERB>` `<NOUN>` pattern.

Don't forget to add a debug command to `.vscode/launch.json`


### Terminal

- `make` - execute the build pipeline.
- `make help` - print help for the [Make targets](Makefile).

### Visual Studio Code

`F1` → `Tasks: Run Build Task (Ctrl+Shift+B or ⇧⌘B)` to execute the build pipeline.

## Release

The release workflow is triggered each time a tag with `v` prefix is pushed.

_CAUTION_: Make sure to understand the consequences before you bump the major version. More info: [Go Wiki](https://github.com/golang/go/wiki/Modules#releasing-modules-v2-or-higher), [Go Blog](https://blog.golang.org/v2-go-modules).

get the latest tag git

```
git describe --tags --abbrev=0
```

when releasing make sure to

1. [ ] run `full-smoke-test` before cutting release to run through some common commands and make sure that they work

2. [ ]  release new version of [workspace-images](https://github.com/brevdev/workspace-images)

3. update [brev's homebrew tap](https://github.com/brevdev/homebrew-brev)


## Maintainance

Remember to update Go version in [.github/workflows](.github/workflows), [Makefile](Makefile) and [devcontainer.json](.devcontainer/devcontainer.json).

Notable files:

- [devcontainer.json](.devcontainer/devcontainer.json) - Visual Studio Code Remote Container configuration,
- [.github/workflows](.github/workflows) - GitHub Actions workflows,
- [.github/dependabot.yml](.github/dependabot.yml) - Dependabot configuration,
- [.vscode](.vscode) - Visual Studio Code configuration files,
- [.golangci.yml](.golangci.yml) - golangci-lint configuration,
- [.goreleaser.yml](.goreleaser.yml) - GoReleaser configuration,
- [Dockerfile](Dockerfile) - Dockerfile used by GoReleaser to create a container image,
- [Makefile](Makefile) - Make targets used for development, [CI build](.github/workflows) and [.vscode/tasks.json](.vscode/tasks.json),
- [go.mod](go.mod) - [Go module definition](https://github.com/golang/go/wiki/Modules#gomod),
- [tools.go](tools.go) - [build tools](https://github.com/golang/go/wiki/Modules#how-can-i-track-tool-dependencies-for-a-module).

## Contributing

Simply create an issue or a pull request.
