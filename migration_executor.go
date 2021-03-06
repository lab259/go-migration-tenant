package mtnt

import (
	"os"

	"github.com/lab259/go-migration"
	"github.com/lab259/go-prdcsm/v2"
)

// MigrationExecutor runs all migrations in all accounts listed by the `AccountProducer`.
type MigrationExecutor struct {
	connector         Connector
	reporter          Reporter
	migrationReporter migration.Reporter
	producer          AccountProducer
	source            migration.Source
	pool              prdcsm.Pool
	args              []string
}

// NewMigrationExecutor returns a new instance of a `MigrationExecutor`.
func NewMigrationExecutor(migrationReporter migration.Reporter, mtntnReporter Reporter, connector Connector, producer AccountProducer, source migration.Source) *MigrationExecutor {
	return &MigrationExecutor{
		connector:         connector,
		producer:          producer,
		source:            source,
		reporter:          mtntnReporter,
		migrationReporter: migrationReporter,
	}
}

// Run initializes the `go-prdcsm` starting the process.
func (e *MigrationExecutor) Run(workers int, args ...string) {
	e.args = args
	producer := newAccountProducerProxy(e, e.producer, workers*2)

	e.pool = prdcsm.NewPool(prdcsm.PoolConfig{
		Workers:  workers,
		Consumer: e.consumer,
		Producer: producer,
	})
	e.pool.Start()
	e.pool.Wait()
}

// Stop stops the running pool of workers.
func (e *MigrationExecutor) Stop() {
	if e.pool == nil {
		return
	}
	e.pool.Stop() // Stop already waits for the pool to finish gracefully.
	e.pool = nil
}

// consumer is a consumer according to the `go-prdcsm` approach.
func (e *MigrationExecutor) consumer(data interface{}) {
	account, ok := data.(Account)
	if !ok {
		panic("account provided is not a Account instance")
	}
	e.migrate(account)
}

// migrate calls the runner for one specific account.
func (e *MigrationExecutor) migrate(account Account) error {
	return account.ProvideMigrationContext(func(executionContext interface{}) error {
		target, err := e.connector(executionContext)
		if err != nil {
			return err
		}
		e.reporter.BeforeAccount(account)
		manager := migration.NewDefaultManager(target, e.source)
		runner := migration.NewArgsRunnerCustom(e.migrationReporter, manager, os.Exit, e.args...)
		runner.Run(executionContext)
		e.reporter.AfterAccount(account)
		return nil
	})
}
