#!/usr/bin/bash
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null && pwd )"

NITROGEN_BIN="${NITROGEN_BIN:-$DIR/../bin/nitrogen}"
export TESTDATA_DIR="$DIR/../testdata"

echo "Run Nitrogen source test suite"
echo

shopt -s globstar
for test in $DIR/**/*.ni; do
    rel_path="$(realpath --relative-to="${DIR}" "$test")"
    echo -n -e "$rel_path - \e[31m"

    "$NITROGEN_BIN" -M $DIR/../nitrogen -M $DIR/../built-modules "$test"

    if [ $? -ne 0 ]; then
        echo -e "\e[0m"
        exit 1
    fi

    echo -e "\e[32mpassed\e[0m"
done
