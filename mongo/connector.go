package mongo

import (
	"errors"

	"github.com/globalsign/mgo"
	"github.com/lab259/go-migration"
)

var ErrInvalidMgoDatabase = errors.New("executionContext is not a *mgo.Database")

// MongoConnector returns a new MongoDBTarget from the execution context.
//
// If the executionContext is not a *mgo.Database, it returns an
// `ErrInvalidMgoDatabase`.
func Connector(executionContext interface{}) (migration.Target, error) {
	db, ok := executionContext.(*mgo.Database)
	if !ok {
		return nil, ErrInvalidMgoDatabase
	}
	return migration.NewMongoDB(db), nil
}
