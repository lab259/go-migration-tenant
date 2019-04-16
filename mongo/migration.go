package mongo

import (
	"errors"
	"time"

	"github.com/globalsign/mgo"
	"github.com/lab259/go-migration"
)

var ErrInvalidExecutionContext = errors.New("invalid execution context: `*ExecutionContext` expected")

type Handler func(executionContext *mgo.Database) error

func makeHandlers(handlers ...Handler) []migration.Handler {
	hs := make([]migration.Handler, len(handlers))
	for i, h := range handlers {
		func(i int, h Handler) {
			hs[i] = func(executionContext interface{}) error {
				ec, ok := executionContext.(*mgo.Database)
				if !ok {
					return ErrInvalidExecutionContext
				}
				return h(ec)
			}
		}(i, h)
	}
	return hs
}

// NewCodeMigration uses `NewCodeMigrationCustom` wrapping its handlers to a
// typed execution context.
func NewCodeMigration(handlers ...Handler) migration.Migration {
	hs := makeHandlers(handlers...)
	return migration.NewCodeMigrationCustom(1, hs...)
}

// NewMigration uses `migration.NewMigration` wrapping its handlers to a typed
// execution context.
func NewMigration(id time.Time, description string, handlers ...Handler) migration.Migration {
	hs := makeHandlers(handlers...)
	return migration.NewMigration(id, description, hs...)
}
