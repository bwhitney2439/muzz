# Muzz Dating App 

## IDE Development

### Visual Studio Code

Use the following plugins, in this boilerplate project:
  - Name: Go
  - ID: golang.go
  - Description: Rich Go language support for Visual Studio Code
  - Version: 0.29.0
  - Editor: Go Team at Google
  - Link to Marketplace to VS: https://marketplace.visualstudio.com/items?itemName=golang.Go


## Features

- **User Creation**: Sign up new users with autogenerated details for testing purposes.
- **Login**: Authenticate users and generate JWT tokens for session management.
- **Discover**: Allow users to discover potential matches based on preferences. 
- **Swipe**: Users can swipe right or left on potential matches to express their interest or lack thereof.

## Getting Started

### Prerequisites

Ensure you have Go installed on your machine and your workspace is set up correctly. This package depends on several external libraries:

- Fiber v2 for the web framework.
- Golang JWT v5 for handling JWT tokens.
- Pioz/faker for generating fake data for users.
- Custom Muzz libraries for database operations and models.

### Installation
```bash
make requirements
```


### Start the application 

```bash
make run
```

### Use local container

```
# Shows all commands
make help

# Clean packages
make clean-packages

# Generate go.mod & go.sum files
make requirements

# Generate docker image
make build

# Generate docker image with no cache
make build-no-cache

# Run the projec in a local container
make up

# Run local container in background
make up-silent

# Run local container in background with prefork
make up-silent-prefork

# Stop container
make stop

# Start container
make start
```

## Production

```bash
make build
make up
```

Go to http://localhost:3000:

## Application Routes Documentation

This document outlines the API routes available in the Muzz application. The application uses the Fiber web framework for Go, with a focus on performance and ease of use. Routes are organized under the `/api/v1` prefix and include both public and protected endpoints. 

### Getting Started

Before diving into the routes, ensure your development environment is set up:

1. Ensure Go is installed.
2. Clone the Muzz repository.
3. Install the necessary Go packages.
4. Start the server using `go run app.go` with optional flags `-port` for specifying the port and `-prod` to enable prefork in production mode.

#### Middleware

The application uses the following middleware:

- **Recover**: To recover from panics anywhere in the stack and handle them gracefully.
- **Logger**: Logs every request to the console.
- **Protected**: A custom middleware to protect routes that require authentication.

### Routes

#### Public Routes

These routes do not require authentication.

- **POST `/api/v1/user/create`**: Register a new user. The user's details are autogenerated for testing purposes.
- **POST `/api/v1/login`**: Authenticate a user and return a JWT token for session management.

#### Protected Routes

These routes require a valid JWT token to be accessed, provided through the `Protected` middleware.

- **GET `/api/v1/discover`**: Discover potential matches based on user preferences. This endpoint accepts the following query parameters:
  - `age`: The age to filter potential matches. Must be a positive integer.
  - `gender`: The gender to filter potential matches. Acceptable values are typically "Male", "Female", or "Other", but check application settings for available options.
  - `orderBy`: Specifies the order of the results. Options may include parameters such as "age" or "location", but specific implementations may vary.

Example request: `/api/v1/discover?age=25&gender=Female&orderBy=age`
- **POST `/api/v1/swipe`**: Swipe right or left on potential matches to express interest or disinterest.

## Starting the Server

To start the server on a custom port (default is 3000), use the `-port` flag:

```sh
go run app.go -port=:XXXX



