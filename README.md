# rewards-api

## Background

Each Monday at 05:00 UTC, a new DIMO issuance week begins; week 0 started on 2022-01-31. At the end of each issuance week, this service runs a cron job that calculates how many points each vehicle earned that week. A simple REST API provides vehicle owners with a view of when and how their vehicles earned.

The formulas are available on the [docs site](https://docs.dimo.zone/dimo-overview/token/demand-signal), but it's worth summarizing the steps here with an emphasis on the implementation:

**Activity.** A vehicle's integration is considered active in a given week if during that time it has transmitted some non-trivial signal. Vehicles without active integrations do not earn points.

**Integration.** Better integrations earn more points, as shown in the table below. The only valid combination of two of these is Smartcar and AutoPi together, which earns the sim 7000.

| Integration | Points |
| ----------- | ------ |
| Smartcar    | 1000   |
| Tesla       | 4000   |
| AutoPi      | 6000   |

**Streak.** The program maintains a "weeks connected" counter for each vehicle. Ideally, this starts at 0 and simply increments by 1 every week, with the vehicle staying continuously connected. Being disconnected for one or two weeks in a row merely pauses this counter.

**Level.** The connection streak at the end of a week determines a level for the vehicle. An active vehicle then earns points based on the level:

| Level | Min streak | Points |
| ----- | ---------- | ------ |
| 1     | 0          | 0      |
| 2     | 4          | 1000   |
| 3     | 20         | 2000   |
| 4     | 36         | 3000   |

Being inactive for three weeks straight drops your connection streak to the minimum one for the previous level.

## Contributing

### Database modifications

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

### Generate API documentation

```
swag init -g cmd/rewards-api/main.go --parseDependency --parseInternal --parseDepth 2
```

### Simulating a production run

1. Get production values for the Clickohouse instance to place in `settings.yaml`. You will need to be on VPN for this to work.
2. Set `DEFINITIONS_API_GRPC_ADDR` and `DEVICES_API_GRPC_ADDR` to local ports. For the sake of an example, let's say these are `localhost:8086` and `localhost:8087`, respectively.
3. Port-forward these two services through. In our example:
   ```sh
   kubectl port-forward -n prod services/device-definitions-api-prod 8086:8086
   kubectl port-forward -n prod services/devices-api-prod 8087:8086
   ```
4. Run the job with, e.g., `go run ./cmd/rewards-api calculate 93`
