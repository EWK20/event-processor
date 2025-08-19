# Event Processor

This project implements a `event processor` service that consumes events from a queue, validates, triages and persists them into a database, so that it can be consumed later to be delivered to clients.

## High Level Architecture

```
Producer -> SQS Queue -> Event Processor -> Postgres Database
```

### Producer

This is a simple message producer that sends messages to the SQS queue every `15 seconds` for testing purposes.
It takes in a set of environment variables:

```
AWS_REGION=xxxxxxx
AWS_ACCESS_KEY_ID=xxxxxxx
AWS_SECRET_ACCESS_KEY=xxxxxxx
SQS_ENDPOINT=xxxxxxx
SQS_QUEUE_NAME=xxxxxxx
SQS_DLQ_QUEUE_NAME=xxxxxxx
```

It depends on the SQS queue being created, so the docker compose will need to be up and running before. This will be covered later.

### Event Processor

- Continuously polls the SQS queue for any new events
- Validates each event against the defined struct located in the `models` package
- Persists valid events into PostgreSQL.
- Send invalid events to a DLQ

### Database

- PostgreSQL stores all events with indexes on client_id, event_type, and timestamp for fast lookups.

## Migrations

Migrations are managed using [Goose](https://github.com/pressly/goose).
Goose is a database migration tool that applies the migrations in `/processor/internal/db/migrations/`.
The migrations must be run before the events processor can work.

## Project Structure

```
├── localstack/
│   ├── init-aws.sh
├── processor/                      # Event Processor
│   ├── cmd/
│   │   ├── migrate.go       Run database migrations command
│   │   ├── process.go      Run events processor
│   │   ├── root.go
│   ├── internal
│   │   ├── config/              Specifies and Gathers environment variables
│   │   ├── db/                   Instantiates database connection and interacts with it
│   │   ├── models/           The event schema that is used to validate data being recieved from producers
│   │   ├── processor/       Processes the data by polling the SQS queue, receiving messages, validating them and persisting them for later consumption
│   ├── .env                       Stores all environment variables
│   ├── go.mod
│   ├── go.sum
│   ├── main.go
├── producer/                       Produces events
│   ├── config/                     Specifies and Gathers environment variables
│   ├── producer/                Produces messages of event type 'transaction_approved' with a slight variation every 15 seconds
│   ├── .env                       Stores all environment variables
│   ├── Dockerfile              Creates docker image for producer
│   ├── go.mod
│   ├── go.sum
│   ├── main.go
│   ├── Makefile                  Used to simplify build command input
├── .gitignore
├── docker-compose.yaml  Starts up crucial services for both applications to work
└── README.md
```

## Running Locally

The producer has a dependency on the SQS queue being available so placing it in the docker compose file was not effective. When adding `depends_on` it starts the producer service immediately the localstack container is available, but the SQS queue may not exist yet, causing the producer to fatal.

### Prerequisites

- Docker
- Golang 1.24

### Step 1 - Start Services

`docker compose up --build`

This will spin up the localstack server with SQS enabled and a postgres database.

### Step 2 - Start Producer

#### Set Env Vars

```
AWS_REGION=xxxxxxx
AWS_ACCESS_KEY_ID=xxxxxxx
AWS_SECRET_ACCESS_KEY=xxxxxxx
SQS_ENDPOINT=xxxxxxx
SQS_QUEUE_NAME=xxxxxxx
```

#### Run program

```
cd producer
go run .
```

### Step 3 - Start Events Processor

#### Set Env Vars

```
DB_USER=xxxxxxx
DB_PASSWORD=xxxxxxx
DB_HOST=xxxxxxx
DB_PORT=xxxx
DB_NAME=xxxxxxx
AWS_REGION=xxxxxxx
AWS_ACCESS_KEY_ID=xxxxxxx
AWS_SECRET_ACCESS_KEY=xxxxxxx
SQS_ENDPOINT=xxxxxxx
SQS_QUEUE_NAME=xxxxxxx
SQS_DLQ_QUEUE_NAME=xxxxxxx
```

#### Run migration

```
cd processor
go run . migrate
```

#### Run events processor

```
cd processor
go run . process
```

## Testing

The testing mainly consists of happy path tests, with some more time I would add some different edge cases to accomodate, such as:

- Malformed events are rejected and sent to the DLQ

To Run:

```
docker compose up
```

This will start the localstack server with sqs enables and a postgres database that the integration tests can use.

```
go test ./... -v
```

This starts all tests in project

## Design Decisions

### SQS over Kafka:

- Fits well with AWS stack.
- Easy to simulate with LocalStack.

### Postgres for persistence:

- Low latency with proper indexes.
- Well-known for reliability.

### DLQ handling:

- Guarantees no data loss.
- Easier debugging of “poison messages”.

### Reproducibility:

- Key systems run with docker-compose up.
- No external dependencies beyond Docker and Golang 1.24.
