# PayoutManagementSystem


This project is about the payout management system built using golang.

# Project Setup

## Clone the repository

Clone this repo: <a href = "https://github.com/Swarathmica-infraspec/payout-management-system"> source link  </a>

# Requirements

GO-VERSION: 1.25.0

The project contains payoutmanagementsystem/ <br>
<<<<<<< HEAD
=======

>>>>>>> a8ea28b (Introduce db for payee)
- .github/workflows/payoutManagementSystem.yml <br>
- payee/
  - payee.go <br>
  - payee_test.go <br>
  - payee_db.sql <br>
<<<<<<< HEAD
  - payeeDAO.go <br>
  - payeeDAO_test.go <br>
=======
>>>>>>> a8ea28b (Introduce db for payee)
- go.mod <br>
- go.sum <br>
- README.md <br>


NOTE: Only email ids with .com are supported.



# Database Setup

We use PostgreSQL(17.6-trixie) running inside Docker for persistant storage.

## 1. Start Postgres with Docker Compose

From the project root, open VS Code. Press F1: Dev Containers: Reopen in Dev Container

This will start PostgreSQL in a container.


## 2. Create Payees Table

Run the below command for the first time (or if db does not exist):
psql -h db -U $POSTGRES_USER -d $POSTGRES_DB -f payee_db.sql

It will prompt for password. Give your postgres password. (or refer to .env)

If 'command not found: psql' : run : apt-get install -y postgresql-client

# Data Access Object

payeeDAO contains database query for payee and payeeDAO_test contains relevant tests

# Run tests

To run tests: (inside docker)

go test -v ./...


NOTE: this project is still under development and hence does not have HTTP API now.

# NOTE:

To exit devcontainer: press F1: Dev Containers: Reopen folder locally


# CI

The workflow is triggered on every push and pull request.
It runs the following checks automatically:
- Format with `test -z "$(gofmt -l .)"`
- Linting with `golangci-lint`
- Tests with `go test`
