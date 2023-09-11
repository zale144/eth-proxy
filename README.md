# Ethereum Balance Service

## Overview

This service provides an API to fetch the balance of a given Ethereum address. It's built using Go and utilizes the `go-ethereum` package for Ethereum interactions.

## Table of Contents

1. [Features](#features)
2. [Prerequisites](#prerequisites)
3. [Installation](#installation)
4. [Configuration](#configuration)
5. [Usage](#usage)
6. [API Endpoints](#api-endpoints)
7. [Monitoring](#monitoring)
8. [Testing](#testing)
9. [Swagger Documentation](#swagger-documentation)
10. [Contributing](#contributing)
11. [License](#license)

## Features

- Retrieve Ethereum balance in a fast and efficient manner.
- Health check and readiness endpoints.
- Prometheus metrics.
- Rate limiting.
- Retrying failed requests to Ethereum node.

## Prerequisites

- Go 1.21 or higher
- Ethereum node access (e.g., via Infura)

## Installation

Clone the repository:

```bash
git clone https://github.com/zale144/eth-proxy.git
```

Navigate to the project directory and build the service:

```bash
cd eth-proxy
make build
```

## Configuration

The service is configured using environment variables. The following variables are available:

- `ETHEREUM_NODE_URL (required)`: URL of the Ethereum node.
- `HTTP_PORT`: Port to listen on.
- `HTTP_TIMEOUT`: Timeout for HTTP requests to Ethereum node.
- `CLIENT_RETRIES`: Number of retries for failed requests.
- `CLIENT_RETRY_DELAY`: Delay between retries.
- `RATE_LIMIT`: Requests per minute per client.
- `RATE_BURST`: Maximum number of requests that can be sent at once.

Create a `.env` file or export these variables in your shell.

## Usage

Run the service:

```bash
 make run
```

The service will start, and you can access it at http://localhost:8080.

## API Endpoints

| Endpoint                | Method | Description                      |
|-------------------------| ------ |----------------------------------|
| `/eth/balance/:address` | GET    | Fetch the balance of an address. |
| `/metrics`              | GET    | Metrics endpoint.                |
| `/healthy`              | GET    | Health check endpoint.           |
| `/ready`                | GET    | Readiness endpoint.              |


## Monitoring

The service exposes a `/metrics` endpoint for Prometheus.

## Testing

Run the tests:

```bash
make test
```

## Swagger Documentation

Swagger documentation is available at http://localhost:8080/swagger/index.html.

## Contributing

For contributions, please create a new branch and submit a pull request. Make sure all tests pass and you've added tests for your feature or fix.

## License

This project is licensed under the MIT License - see the LICENSE file for details.

