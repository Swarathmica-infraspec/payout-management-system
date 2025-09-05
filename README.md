# PayoutManagementSystem


This project is about the payout management system built using golang.

# Project Setup

## Clone the repository

Clone this repo: <a href = "https://github.com/Swarathmica-infraspec/payoutManagementSystem"> source link  </a>

# Requirements

GO-VERSION: 1.25.0

The project contains payoutmanagementsystem/ <br>
    - .github/workflows/payoutManagementSystem.yml <br>
    - payoutmanagementsystem/
        - payee.go <br>
        - payee_test.go <br>
    - go.mod <br>
    - go.sum <br>
    - README.md <br>

NOTE: Only email ids with .com are supported.

# Run Tests

Test can be run by executing the below command in the terminal
  go test -v ./...

NOTE: this project is still under development and hence does not have HTTP API now.

# CI

The workflow is triggered on every push and pull request.
It runs the following checks automatically:
- Format with `test -z "$(gofmt -l .)"`
- Linting with `golangci-lint`
- Tests with `go test`
