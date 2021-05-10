#!/bin/sh

set -ue

diff="$(git diff)"

if [ "$diff" != "" ]; then
    echo "$diff"
    exit 1
fi

exit 0
