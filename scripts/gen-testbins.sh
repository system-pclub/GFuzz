#!/bin/bash -e
EXCLUDE_PATHS="integration"
TESTDIRS="./..."
pkg_list=$(go list $TESTDIRS | grep -vE "($EXCLUDE_PATHS)")

for pkg in $pkg_list
do
    echo "generating test bin for $pkg"
    name=$(echo "$pkg" | sed "s/\//-/g")
    go test -c -o $1/$name.test $pkg
done