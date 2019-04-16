package mtnt

import (
	"github.com/lab259/go-migration"
)

// Connector should create a new `migration.Target` from the provided execution
// context.
//
// All default implementation expect the `executionContext` itself to be the
// database connection reference. So, if you need something more complex, you
// should create your own.
//
// See also `MongoConnector`.
type Connector func(executionContext interface{}) (migration.Target, error)
