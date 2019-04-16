# go-tenant-migration

go-tenant-migration, or mtnt for short, is designed to run migrations on tenants
in parallel.

This lihbrary is an extention of the [go-migration](https://github.com/lab259/go-migration).

## TL;DR

_Account_ and _AccountProducer_ interfaces should be implemented. Where _Account_
should provide the account database reference and _AccountProducer_ should
implement listing the _Account_ instaces that should be migrated.

The rest should look like this:

```go
package main

import (
	"os"
	
	mtnt "github.com/lab259/go-migration-tenant"
	mtntMongo "github.com/lab259/go-migration-tenant/mongo"
	migration "github.com/lab259/go-migration"
	
	"yoursystem"
	_ "yoursystem/migrations"
)

func main() {
	executor := mtnt.NewMigrationExecutor(
		migration.NewDefaultReporter(), 
		mtnt.NewDefaultReporter(os.Stdout),
		mtntMongo.Connector,
		yoursystem.NewAccountProducer(), // YOU MUST implement the account producer
		migration.DefaultCodeSource(),
	)
	executor.Run(10, os.Args...) // <- number of workers that will execute migrations, in other
	                            // words how many accounts you want to migrate in parallel.
}
```

## Usage

### Listing accounts

In order to run the system, you need to implement 2 interfaces: `Account` and
`AccountProducer`.

**Account** _(interface)_

Account represents each tenant of your system. It should hold the database
reference and provides the execution context for the migration. The execution
context will be passed to each handler as it is called.

:exclamation: **Hence, it is up to you to produce an execution context that will 
provide the database connection for the handlers.** The default implementation
expects you to provide the database connection reference as execution context.
But, it can be easily changed by implementing a new `Connector`.

```go
Identification() string
ProvideMigrationContext(func(executionContext interface{}) error) error
```

**AccountProducer** _(interface)_

AccountProducer should list all the `Account`s that need to be migrated.

```go
Get() ([]Account, error)
HasNext() bool
Total() int
```

### Running a migration

Once implemented `Account` and `AccountProducer`, the library is able to run all
migrations.

```go
package main

import (
	"os"
	
	mtnt "github.com/lab259/go-migration-tenant"
	mtntMongo "github.com/lab259/go-migration-tenant/mongo"
	migration "github.com/lab259/go-migration"
	
	"yoursystem"
	_ "yoursystem/migrations"
)

func main() {
	executor := mtnt.NewMigrationExecutor(
		migration.NewDefaultReporter(), 
		mtnt.NewDefaultReporter(os.Stdout),
		mtntMongo.Connector,
		yoursystem.NewAccountProducer(), // YOU MUST implement the account producer
		migration.DefaultCodeSource(),
	)
	executor.Run(10, os.Args...) // <- number of workers that will execute migrations, in other
	                // words how many accounts you want to migrate in parallel.
}
```

Now, after compiling the system:

```bash
./bin/migration migrate
```

To know more about the migration commands you can:

```bash
./bin/migration --help
```

### What if I need to filter one specific account

Well, the `AccountProducer` implementation is up to you. So, you should filter
accounts while implementing the producer.

## Concurrency

The library uses [gp-prdcsm](https://github.com/lab259/go-prdcsm) to implement a
Producer & Consumer pattern to achieve concurrent migrations. Migrations at the 
account level are not ran in parallel. But, two, or more, accounts should run
in parallel.

## License

MIT