# ksec

A command line tool that simplifies the management of Kubernetes Secrets.
- Easily set and unset k8s Secret keys
- Handles base64 encoding
- Supports synchronization of `.env` files to k8s Secrets

## Installation

Compiled binaries can be found in the [GitHub releases](https://github.com/kanopy-platform/ksec/releases).

Install compiled binary as a Helm plugin (requires [Helm](https://docs.helm.sh/using_helm/#installing-helm)).

    helm plugin install https://github.com/kanopy-platform/ksec

Install from source (requires [golang](https://golang.org/doc/install#install)).

    go get github.com/kanopy-platform/ksec/cmd/...

### Trobleshooting

#### Windows: Spaces in `HELM_HOME`

`helm plugin install` may not work on Windows if you have spaces in your `HELM_HOME` path. You can instead download the windows executable from the [latest release](https://github.com/kanopy-platform/ksec/releases/latest).

#### Helm Install Errors

You may need to remove a plugin and clear your plugin cache if you have errors installing a plugin.

macOS
```
# rm cache
rm -rf ~/Library/Caches/helm/plugins/https-github.com-kanopy-platform-ksec*

# rm plugin
rm -rf ~/Library/helm/plugins/ksec/
```

## Usage
```
A tool for managing Kubernetes Secret data

Usage:
  ksec [command]

Available Commands:
  completion  Generate command completion scripts
  create      Create a Secret
  delete      Delete a Secret
  get         Get values from a Secret
  help        Help about any command
  list        List all secrets in a namespace
  pull        Pull values from a Secret into a .env file
  push        Push values from a .env file into a Secret
  set         Set values in a Secret
  unset       Unset values in a Secret

Flags:
      --config string      config file (Default: $HOME/.ksec.yaml)
  -h, --help               help for ksec
  -n, --namespace string   Operate in a specific NAMESPACE (Default: current kubeconfig namespace)
      --version            version for ksec

Use "ksec [command] --help" for more information about a command.
```

## Development

Run `make` to run all tests and create a new binary in `${GOPATH}/bin/`

Run `VERSION=<version_number> make dist` to cross compile binaries into a `./dist/` directory. These binaries can then be uploaded to a new GitHub release.
