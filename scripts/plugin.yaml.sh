#!/bin/bash

cat <<EOF > "plugin.yaml"
---
name: "ksec"
version: "${1}"
usage: "Manage Kubernetes Secrets through Helm"
description: "Helm plugin that simplifies the management of Kubernetes Secrets"
ignoreFlags: false
useTunnel: false
command: "\$HELM_PLUGIN_DIR/ksec"
hooks:
  install: "cd \$HELM_PLUGIN_DIR && scripts/install-binary.sh"
EOF
