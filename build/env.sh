#!/bin/sh

set -e

if [ ! -f "build/env.sh" ]; then
    echo "$0 must be run from the root of the repository."
    exit 2
fi

# Create fake Go workspace if it doesn't exist yet.
workspace="$PWD/build/_workspace"
root="$PWD"
krdir="$workspace/src/github.com/krypton"
if [ ! -L "$krdir/go-krypton" ]; then
    mkdir -p "$krdir"
    cd "$krdir"
    ln -s ../../../../../. go-krypton
    cd "$root"
fi

# Set up the environment to use the workspace.
# Also add Godeps workspace so we build using canned dependencies.
GOPATH="$krdir/go-krypton/Godeps/_workspace:$workspace"
GOBIN="$PWD/build/bin"
export GOPATH GOBIN

# Run the command inside the workspace.
cd "$krdir/go-krypton"
PWD="$krdir/go-krypton"

# Launch the arguments with the configured environment.
exec "$@"
