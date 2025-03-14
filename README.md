# SWIFT Code API

This is my first project created in GO, it imports data from CSV file and provides a REST API for managing SWIFT codes. It is built using Go with the `chi` router and utilizes MySQL for database storage. The application is containerized using Docker for easy deployment.

## Getting Started

### Prerequisites

Ensure you have Docker installed on your system, the best way is to install Docker Desktop.

[https://www.docker.com/get-started/]

You need docker desktop running to start the program.

## How to Run

1. **Download the repository**:
   - Clone the repository using Git:  
     ```sh
     git clone https://github.com/your-username/SWIFT-Parser.git
     ```
   - Or download the repository as a ZIP and extract it.

2. **Navigate to the project directory**:
   ```sh
   cd SWIFT-Parser
   ```

## Running the Application

To start the application, run the following command inside the project directory:

```sh
docker compose up --build
```

This will build and start the containers, including the API and the MySQL database.

Once started, the API should be accessible at:
```
http://localhost:8080
```

## Running Tests

To run unit and integration tests, use the following command:

```sh
docker compose -f docker-compose.test.yml up --build
```

This will build and run the test environment, executing all tests automatically.

## Database

The application uses MySQL as the database engine.
Application will automacally create necessary tables and import data from CSV file, if it's not already imported.

> **Important**: I assumed the address is not mandatory, and the program imports records with empty address fields.
Every other column cannot be empty, and if it is, the record will be ommited.

If you need to reset the database, you can remove the volume and recreate the containers:
```sh
docker compose down -v
```

## API Endpoints

### Get SWIFT Code
**Endpoint:** `GET /v1/swift-codes/{swift-code}`


### Get Country SWIFT Codes
**Endpoint:** `GET /v1/swift-codes/country/{country-iso2-code}`


### Create a SWIFT Code Entry
**Endpoint:** `POST /v1/swift-codes`

**Request Body:**
```json
{
    "address": "123 New Bank Street",
    "bankName": "New Bank",
    "countryISO2": "PL",
    "countryName": "Poland",
    "isHeadquarter": true,
    "swiftCode": "ABCDEFGHXXX"
}
```


### Delete SWIFT Code
**Endpoint:** `DELETE /v1/swift-codes/{swift-code}`




## Stopping and Cleaning Up
To stop the application, press `CTRL + C` in the terminal or run:
```sh
docker compose down
```

