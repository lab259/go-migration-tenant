package mongo_test

import (
	"time"

	"github.com/globalsign/mgo"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/lab259/go-migration-tenant/mongo"
)

var _ = Describe("MongoConnector", func() {
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

	It("should create a new target", func() {
		_, err := mongo.Connector(session.DB("database1"))
		Expect(err).ToNot(HaveOccurred())
	})

	It("should fail creating a new target", func() {
		_, err := mongo.Connector(12345)
		Expect(err).To(Equal(mongo.ErrInvalidMgoDatabase))
	})
})
