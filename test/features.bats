#!/usr/bin/env bats

# Acceptance test suite, made with BATS.
# https://github.com/bats-core/bats-core
# Run:
#     make test-acceptance


@test "Running without parameters shows the help" {
  run "${mj}"
  assert_success
  assert_output --partial 'Rank proposals of a Majority Judgment poll.'
}

@test "Run with an input CSV (example.csv)" {
  run "${mj}" example/example.csv
  assert_success
  assert_output --partial '#2   Pizza'
  assert_output --partial '#1   Chips'
  assert_output --partial '#3   Pasta'
}

# WIP: write more tests !

# ----------------------------------------------------------------------------

setup() {
    load 'test_helper/bats-support/load'
    load 'test_helper/bats-assert/load'
    export TZ="Europe/Paris"

    TESTS_DIR="$( cd "$( dirname "$BATS_TEST_FILENAME" )" >/dev/null 2>&1 && pwd )"
    PROJECT_DIR="$( dirname "$TESTS_DIR" )"
    COVERAGE_DIR="${PROJECT_DIR}/test-coverage"
    mj="${PROJECT_DIR}/mj"

    cd "${PROJECT_DIR}" || exit

    if [ "$MJ_COVERAGE" == "1" ] ; then
      echo "Setting up coverage in ${COVERAGE_DIR}"
      mkdir -p "${COVERAGE_DIR}"
      export GOCOVERDIR=${COVERAGE_DIR}
      mj="${mj}-coverage"
    fi
}

teardown() {
    true
    #rm -f test.log
}
