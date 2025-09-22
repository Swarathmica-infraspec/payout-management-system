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
  - payeeDAO.go <br>
  - payeeDAO_test.go <br>
  - payeeAPI.go <br>
  - payeeApi_test.go <br>
- go.mod <br>
- go.sum <br>
- README.md <br>

NOTE: Only email ids with .com are supported.

# Database Setup

We use PostgreSQL(17.6-trixie) running inside Docker for persistant storage.

Install devcontainer extension in vs code, or from the terminal using
npm i -g @devcontainers/cli

## 1. Start Postgres with Docker Compose

From the project root, open VS Code. Press F1: Dev Containers: Reopen in Dev Container

This will start PostgreSQL in a container.

To start devcontainer using terminal:
in your project root, run,
devcontainer up --workspace-folder

To get into the container:
devcontainer exec --workspace-folder . bash


## 2. Create Payees Table

Run the below command for the first time (or if db does not exist):
psql -h db -U $POSTGRES_USER -d $POSTGRES_DB -f payee/payee_db.sql

It will prompt for password. Give your postgres password. (or refer to .env)

If 'command not found: psql' : run : apt-get update
                                     apt-get install -y postgresql-client

## 3. Data Access Object

payeeDAO contains database query for payee and payeeDAO_test contains relevant tests

## 4. HTTP API Usage

since postgres is run from docker, 

docker exec -it devcontainer-app-1 bash

cd /workspaces/payoutManagementSystem

then run: go run main.go #entry point

payeeApi.go has the code for API while payeeAPI_test.go has test code



1. POST request 
curl -X POST http://localhost:8080/payees \
  -H "Content-Type: application/json" \
  -d '{
    "name":"Abc",
    "code":"123",
    "account_number":1234567890,
    "ifsc":"CBIN0123456",
    "bank":"CBI",
    "email":"abc@example.com",
    "mobile":9876543210,
    "category":"Employee"
  }'

expected response: {'id':1}

2. GET request
curl -X GET http://localhost:8080/payees/list \
  -H "Content-Type: application/json"

3. GET by id request
curl -X GET http://localhost:8080/payees/1 \
  -H "Content-Type: application/json"

4. PUT request

<!-- SUPPOSE THE ROW GIVEN IN POST IS PRESENT IN DB -->
curl -X PUT http://localhost:8080/payees/update/1 \
  -d '{
    "name":"ABCD",
    "code":"123",
    "account_number":1234567890,
    "ifsc":"CBIN0123456",
    "bank":"CBI",
    "email":"abc@example.com",
    "mobile":9876543210,
    "category":"Employee"
  }'

  expected response: {"status":"updated"}

  5. DELETE request
<!-- SUPPOSE THE ROW GIVEN IN POST IS PRESENT IN DB -->
  
curl -X DELETE http://localhost:8080/payees/delete/1

expected response: {"status":"deleted"}


# Run Tests
To run tests:
go test -v ./...

# CI

The workflow is triggered on every push and pull request.
It runs the following checks automatically:
- Format with `test -z "$(gofmt -l .)"`
- Linting with `golangci-lint`
- Tests with `go test`
