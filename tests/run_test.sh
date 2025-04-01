#!/usr/bin/env bash
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null && pwd )"
TEST="$1"

if [ -z "$TEST" ]; then
    echo "No test file"
    exit 1
fi

NITROGEN_BIN="${NITROGEN_BIN:-$DIR/../bin/nitrogen}"
export TESTDATA_DIR="$DIR/../testdata"

echo "Running Nitrogen source test suite"
echo

rel_path="$(realpath --relative-to="${DIR}" "$TEST")"
echo -n -e "$rel_path - \e[31m"

"$NITROGEN_BIN" -M $DIR/../nitrogen -M $DIR/../built-modules -nonibs "$TEST"

if [ $? -ne 0 ]; then
    echo -e "\e[0m"
    exit 1
fi

echo -e "\e[32mpassed\e[0m"
