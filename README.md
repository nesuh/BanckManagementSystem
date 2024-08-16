# Bank Management API

This project is a simple RESTful API for managing bank accounts, built with Go (Golang). The API allows you to create accounts, retrieve account details, and list all accounts.

## Table of Contents
- [Installation](#installation)
- [Usage](#usage)
- [API Endpoints](#api-endpoints)
- [Project Structure](#project-structure)
-

## Installation

To run this project locally, follow these steps:

1. Clone the repository:
    ```sh
    git clone https://github.com/yourusername/bank-management-api.git
    cd bank-management-api
    ```

2. Install the dependencies:
    ```sh
    go mod download
    ```

3. Run the application:
    ```sh
    go run main.go
    ```

The server will start on `http://localhost:8080`.

## Usage

You can interact with the API using tools like `curl`, `Postman`, or directly from your application. Below are examples of how to use each endpoint.

### Create a New Account

```sh
curl -X POST http://localhost:8080/accounts \
-H "Content-Type: application/json" \
-d '{
"first_name": "John",
 "last_name": "Doe",
"balance": 1000}'


                        API Endpoints
=================POST /accounts
Description: Create a new bank account.
Request Body:
first_name: First name of the account holder.
last_name: Last name of the account holder.
balance: Initial balance of the account.
Response:
Returns the created account object.
=================GET /accounts
Description: Retrieve all bank accounts.
Response:
Returns an array of account objects.
=================GET /accounts/{id}
Description: Retrieve a bank account by its ID.
Parameters:
id: The ID of the account to retrieve.
Response:
Returns the account object if found, otherwise an error.
================== Project Structure
Here's a brief overview of the project's structure:

main.go: The entry point of the application. It sets up the HTTP server and routes.
api.go: Contains the HTTP handlers for the API endpoints.
storage.go: Implements in-memory storage for accounts.
types.go: Defines the data structures used in the application.
