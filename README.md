# Test Task FIO Enrichment Service

The FIO Enrichment Service is designed to enrich provided full names (FIO) using open APIs to estimate the most likely age, gender, and nationality. The enriched data is stored in a PostgreSQL database, and the service offers RESTful API endpoints for managing and querying this data.

## Overview

This service covers the following functionality:

1. **RESTful API Endpoints**: Implement several REST methods for managing person data, including filtering and pagination.

2. **Data Enrichment**: Enrich the data with the most likely age, gender, and nationality using the following APIs:
    - Age is enriched using the [Agify API](https://api.agify.io/?name=Dmitriy).
    - Gender is enriched using the [Genderize API](https://api.genderize.io/?name=Dmitriy).
    - Nationality is enriched using the [Nationalize API](https://api.nationalize.io/?name=Dmitriy).

3. **Database Storage**: Save the enriched data in a PostgreSQL database. The database schema is created through migrations.

4. **Logging**: Implement proper logging with debug and info messages to monitor the service's operation.

## REST API Endpoints

The service exposes the following RESTful endpoints:

1. `GET /users`: Retrieve enriched data with optional filters and pagination. You can use them with additional parameters in request: `gender`, `page`, `limit`
2. `PUT /update`: Update an existing person's information by their identifier.
3. `DELETE /delete`: Remove a person's record by their identifier. Use it with parameter `id` 
4. `POST /new`: Create new person records in a specific format.

## Data Enrichment

The service enriches the provided full name (FIO) by calling the following APIs:
- Age is enriched using the Agify API.
- Gender is enriched using the Genderize API.
- Nationality is enriched using the Nationalize API.

The enriched data is saved in the database for future reference.

## Dependencies:

- [chi router](https://github.com/go-chi/chi)
- [migrations](https://github.com/golang-migrate/migrate)
- [driver for DB](https://github.com/jackc/pgconn)
## Usage

To use this app you should go to project folder. Then execute this command

```
make up_build
```
to build the backend app.

After you did this, execute this 
```
make start
```
Make sure you have docker and make installed
