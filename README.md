## Background

Each Monday at 05:00 UTC, a new DIMO issuance week begins; week 0 started on 2022-01-31. At the end of each issuance week, this service runs a cron job that calculates how many points each vehicle earned that week. A simple REST API provides vehicle owners with a view of when and how their vehicles earned.

The formulas are available on the [docs site](https://docs.dimo.zone/dimo-overview/token/demand-signal), but it's worth summarizing the steps here with an emphasis on the implementation:

**Activity.** A vehicle's integration is considered active in a given week if during that time it has transmitted some non-trivial signal. Vehicles without active integrations do not earn points.

**Integration.** A Smartcar integration earns 1000, Tesla 4000, and AutoPi 6000. The only valid combination of two of these is Smartcar together with AutoPi, which earns the sum 7000.

**Streak.** The program maintains a "weeks connected" counter for each vehicle. Ideally, this starts at 0 and simply increments by 1 every week, with the vehicle staying continuously connected.

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
