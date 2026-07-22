#!/bin/bash

##########################
#
# usage:
# ./test.sh -t <test-name>
#
# e.g. ./test.sh -t TestApiAccount
#
##########################

TEST_POOL="TestApiAccount|TestApiCategory|TestApiDevice|TestApiQrcode|TestApiSetting"

COMPOSE_FILE="docker-compose.yaml"

TIMEOUT=300 # 5 minutes

TARGET_TEST=""

while [[ $# -gt 0 ]]; do
    case "$1" in
        -t|--test)
            TARGET_TEST="$2"
            shift 2
            ;;
        *)
            break
            ;;
    esac
done

# check if the test name is in the allowed test pool
if [[ ! "$TARGET_TEST" =~ ^($TEST_POOL)$ ]]; then
    echo "Error: test name '$TARGET_TEST' is not in the allowed test pool"
    echo "Allowed tests: $TEST_POOL"
    exit 1
fi

# remove remaining test db
rm -rf db-test

# create test db
mkdir db-test

# Up the containers using the selected compose file
if ! docker compose -f "$COMPOSE_FILE" up -d --wait --wait-timeout "$TIMEOUT"; then
    echo "Error: Failed to start containers using $COMPOSE_FILE"
    exit 1
fi

sleep 3

# run test
echo "Running test... $TARGET_TEST"
cd goTest
go mod tidy

go test -v -vet=off -run $TARGET_TEST
exit_code=$?
cd ..

# Cleanup: Stop and remove the containers after the test
if ! docker compose -f "$COMPOSE_FILE" down; then
    echo "Warning: Failed to stop and remove containers using $COMPOSE_FILE"
fi

echo "Test completed with exit code: $exit_code"
exit $exit_code