## Background

Each Monday at 05:00 UTC, a new DIMO issuance week begins; week 0 started on 2022-01-31. At the end of each issuance week, this service runs a cron job that calculates how many points each vehicle earned that week. The results from all weeks are visible to vehicle owners through a simple REST API.

## Database modifications

Create a new Goose migration file:
```
goose -dir migrations create MIGRATION_TITLE sql
```
This will create a file in the `migrations` folder named something like `$TIMESTAMP_MIGRATION_TITLE.sql`. Edit this with your new innovations. To run the migrations:
```
go run ./cmd/rewards-api migrate
```
And then to generate the models:
```
sqlboiler psql --no-tests --wipe
```

## Generate API documentation

```
swag init -g cmd/rewards-api/main.go --parseDependency --parseInternal --generatedTime true --parseDepth 2
```
