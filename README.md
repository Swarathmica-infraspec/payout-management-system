# PayoutManagementSystem


This project is about the payout management system built using golang.

# Project Setup

## Clone the repository

Clone this repo: <a href = "https://github.com/Swarathmica-infraspec/payout-management-system"> source link  </a>

# Requirements

GO-VERSION: 1.25.0

The project contains payoutmanagementsystem/ <br>

- .github/workflows/payoutManagementSystem.yml <br>
- expense/
  - expense.go <br>
  - expense_test.go <br>
  - expense_db.sql <br>
  - expenseDAO.go <br>
  - expenseDAO_test.go <br>
- go.mod <br>
- go.sum <br>
- README.md <br>

# Database Setup

We use PostgreSQL running inside Docker for persistant storage.

## 1. Start Postgres with Docker Compose

From the project root, run:

docker compose up -d db


This will:

Start a container named devcontainer-db-1 (from .devcontainer/docker-compose.yml)


## 2. Create Payees Table

Copy the SQL file into the container:

docker cp expense/expense_db.sql devcontainer-db-1:/expense_db.sql


Then apply it:

docker exec -it devcontainer-db-1 psql -U postgres -d postgres -f /expense_db.sql


# Data Access Object

1. expenseDAO contains database query for expense and expenseDAO_test contains relevant tests

To run tests:

docker exec -it devcontainer-app-1 bash

cd /workspaces/payoutManagementSystem

go test -v ./...


NOTE: Only email ids with .com are supported.



# Database Setup

We use PostgreSQL running inside Docker for persistant storage.

## Start Postgres with Docker Compose

From the project root, open VS Code. Press F1: Dev Containers: Reopen in Dev Container

This will start PostgreSQL in a container.


# Run tests

To run tests:
go test -v ./...

<<<<<<< HEAD
Test can be run by executing the below command in the docker terminal
  go test -v ./...
=======
>>>>>>> origin/master

NOTE: this project is still under development and hence does not have HTTP API now.

## To come out of devcontainer:

press F1: Dev Containers: Reopen Folder Locally


# CI

The workflow is triggered on every push and pull request.
It runs the following checks automatically:
- Format with `test -z "$(gofmt -l .)"`
- Linting with `golangci-lint`
- Tests with `go test`
