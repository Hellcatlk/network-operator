#!/usr/bin/sh

set -ue

exitcode=0

for file in $(find . -path ./.git -prune -o -type f | grep -E ".*\.go$")
do
    result=$(gofmt -s -d -w "$file")
    if [ "$result" != "" ];then
        echo "$result"
        exitcode=1
    fi
done

exit $exitcode
