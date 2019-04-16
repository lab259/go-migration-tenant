package mongo_test

import (
	"time"

	"github.com/globalsign/mgo"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/lab259/go-migration-tenant/mongo"
)

var _ = Describe("NewMigration", func() {
	var session *mgo.Session

	BeforeEach(func() {
		sess, err := mgo.DialWithTimeout("mongodb://localhost:27017", time.Second*3)
		Expect(err).ToNot(HaveOccurred())
		session = sess
		dbs, err := session.DatabaseNames()
		Expect(err).ToNot(HaveOccurred())
		for _, db := range dbs {
			switch db {
			case "admin", "config", "local", "test":
				continue
			}
			Expect(session.DB(db).DropDatabase()).To(Succeed())
		}
	})

	It("should call the handlers", func() {
		executionContext1 := session.DB("database1")
		executionContext2 := session.DB("database2")
		doExecuted, undoExecuted := false, false
		migration := mongo.NewMigration(time.Now(), "Description 01", func(db *mgo.Database) error {
			Expect(db.Name).To(Equal("database1"))
			doExecuted = true
			return nil
		}, func(db *mgo.Database) error {
			Expect(db.Name).To(Equal("database2"))
			undoExecuted = true
			return nil
		})
		Expect(migration.Do(executionContext1)).To(Succeed())
		Expect(migration.Undo(executionContext2)).To(Succeed())
		Expect(doExecuted).To(BeTrue(), "the migration was not migrated.")
		Expect(undoExecuted).To(BeTrue(), "the migration was not rewinded.")
	})

	It("should fail trying to migration with an invalid *ExecutionContext", func() {
		migration := mongo.NewMigration(time.Now(), "Description 01", func(executionContext *mgo.Database) error {
			Fail("the Do method should not be called")
			return nil
		})
		Expect(migration.Do(map[string]interface{}{
			"db": "database1",
		})).To(Equal(mongo.ErrInvalidExecutionContext))
	})
})
