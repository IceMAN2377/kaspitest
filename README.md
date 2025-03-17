# EGOV MINI

## RUN INSTRUCTIONS

1. Launch docker-compose file with PostgreSQL container:

```
docker compose up -d
```

2. Set env variables:

```
export HTTP_PORT=8080
export POSTGRES_HOST="localhost"
export POSTGRES_PORT=5432
export POSTGRES_USER="egov"
export POSTGRES_PASSWORD="secret"
export POSTGRES_DB="egov"
```

3. Start app:
```
go run ./cmd/egov
```

## ENDPOINTS

- GET /iin_check/{iin}
- GET /people/info/iin/{iin}
- GET /people/info/phone/{search}
- POST POST /people/info
