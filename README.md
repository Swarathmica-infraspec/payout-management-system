# PayoutManagementSystem


This project is about the payout management system built using golang.

# Project Setup

## Clone the repository

Clone this repo: <a href = "https://github.com/Swarathmica-infraspec/payout-management-system"> source link  </a>

# Requirements

GO-VERSION: 1.25.0

The project contains payoutmanagementsystem/ <br>

- .github/workflows/payoutManagementSystem.yml <br>
- payee/
  - payee.go <br>
  - payee_test.go <br>
  - payee_db.sql <br>
- go.mod <br>
- go.sum <br>
- README.md <br>


NOTE: Only email ids with .com are supported.



# Database Setup

We use PostgreSQL running inside Docker for persistant storage.

## Start Postgres with Docker Compose

From the project root, open VS Code. Press F1: Dev Containers: Reopen in Dev Container

This will start PostgreSQL in a container.


# Run tests

To run tests:
go test -v ./...


NOTE: this project is still under development and hence does not have HTTP API now.

## To come out of devcontainer:

press F1: Dev Containers: Reopen Folder Locally


# CI

The workflow is triggered on every push and pull request.
It runs the following checks automatically:
- Format with `test -z "$(gofmt -l .)"`
- Linting with `golangci-lint`
- Tests with `go test`
