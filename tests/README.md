# Testing

## Methodology

All tests run through some kind of integration testing.
Each test has a `run.sh` script which will execute that test.

## Status

Currently there are two sets of tests: the original tests and specific ones in scenarios.

The Scenario Tests are intended to work.
The Original Tests are being deprecated and will be removed shortly.

## Scenario Tests

These come from the items in `designdocs/20190216-scenarios.md` and are based on behavior driven scenarios of the user use cases.

## Original Tests

The Original Tests were done as a bit of scaffolding before the conductor component was put into place.
They invoke some internal calls that are not meant to really be used by users.
Given some removal of that scaffolding, the Original Tests are being deprecated and will be removed shortly.
They will probably come back in some form depending on how to test the internal functions.