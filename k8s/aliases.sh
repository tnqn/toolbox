#!/bin/bash

# curl -L https://raw.githubusercontent.com/tnqn/toolbox/main/k8s/aliases.sh | bash

set -e

src_url=https://raw.githubusercontent.com/tnqn/toolbox/main/k8s/aliases
dst_path=~/.bash_aliases.k8s

wget -q -4 "$src_url" -O $dst_path
echo "+++ Downloaded $dst_path"

if ! grep -Fxq ". $dst_path" ~/.bashrc; then
cat >> ~/.bashrc <<EOF
. $dst_path
EOF
echo "+++ Added K8s aliases to ~/.bashrc"
else
echo "+++ K8s aliases is already in ~/.bashrc"
fi
