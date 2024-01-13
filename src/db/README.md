# DB

This service is only a thin wrapper on postgres, which executes our init.sql on startup.
This is used to generate our testdata.

## How to use the DB-Service

Start a postgres instance and execute our init.sql or use docker and execute the following commands:

```bash
docker build -t db .
docker run -dit -p 5432:5432 -e POSTGRES_USER=<username> -e POSTGRES_PASSWORD=<password> -e POSTGRES_DB=<db-name> db
```
