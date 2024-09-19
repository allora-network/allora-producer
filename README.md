# Allora Chain Producer

This repository implements a service to monitor Allora Chain blocks, parse specific transactions and events, and send them to Kafka server topics.

## Features

- Monitors Allora Chain blocks in real-time
- Parses configurable transactions and events
- Sends parsed data to specified Kafka topics
- Offers flexible configuration for data processing and topic associations
- Optimized for efficient processing and near real-time data handling

## Configuration

Transactions, events, and their associations with Kafka topics are configurable.

The configuration is managed via the `config.yaml` file, which supports overriding through environment variables. We use [Viper](https://github.com/spf13/viper) to help loading the configurations.

### Configuring Transactions and Events

Events and transactions can be configured at the `filter_event` and `filter_transaction` sections of the configuration file. Example:

```yaml
filter_event:
  types:
    - "emissions.v4.EventNetworkLossSet"
    - "emissions.v4.EventRewardsSettled"

filter_transaction:
  types:
    - "emissions.v4.RemoveStakeRequest"
    - "emissions.v4.InsertReputerPayloadRequest"
```

You must use the fully qualified name of the protobuf type.

### Configuring association with Kafka topics

Events and Transactions can be associated with Kafta topics at the kafka_topic_router section:

```yaml
kafka_topic_router:
  - name: "domain-event.input.allora.staking"
    types:
      - "emissions.v4.AddStakeRequest"
      - "emissions.v4.RemoveStakeRequest"
  - name: "domain-event.input.allora.delegation"
    types:
      - "emissions.v3.DelegateStakeRemovalInfo"
```

`name` refers to the topic name and `types` is a list of fully qualified names of the types.

### Other configurations

Postgres, Kafka and Allora RPC configurations can also be performed at the config file:

```yaml
database:
  url: "postgres://devuser:devpass@localhost:5432/dev-indexer?sslmode=disable"

kafka:
  seeds:
    - "localhost:9092"
  user: "user"
  password: "pass"

allora:
  rpc: "https://allora-rpc.testnet.allora.network"
  timeout: 10

```

### Overriding configurations using environment variables

To override configurations using environment variables, you:
1. Join the different levels of the yaml config using underscore (`_`).
2. Change it to uppercase
3. Add this prefix: `ALLORA_`

Example:

```yaml
database:
  url: "postgres://devuser:devpass@localhost:5432/dev-indexer?sslmode=disable"
```

Becomes:

```yaml
ALLORA_DATABASE_URL: "postgres://devuser:devpass@localhost:5432/dev-indexer?sslmode=disable"
```

## Setup and Usage

[Add instructions on how to set up and run the service]

## Project Structure

- `cmd/`: Contains the main entry point for the application (producer).
- `app/`: Implements the business logic, including the use cases for monitoring transactions and events.
- `app/domain/`: Defines core abstractions such as interfaces for repositories and services.
- `infra/`: Handles interactions with external infrastructure like Kafka, PostgreSQL, and the Allora blockchain.
- `util/`: Contains utility functions like logging execution time.
- `config/`: Manages configuration loading and validation using Viper.

## Technologies

- Go 1.22+
- Kafka (Franz-Go)
- Postgres (pgx)
- Cosmos SDK Go client
- Mockery

## Testing

1. Unit tests:
```bash
make test
```

