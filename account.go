package mtnt

// Account represents
type Account interface {
	Identification() string
	ProvideMigrationContext(func(executionContext interface{}) error) error
}
