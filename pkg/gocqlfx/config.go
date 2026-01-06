package gocqlfx

import "fmt"

// Config holds the configuration for Cassandra/ScyllaDB database connection.
type Config struct {
	// Hosts is a list of database node addresses
	Hosts []string
	// Keyspace is the name of the keyspace to connect to
	Keyspace string
	// Username for database authentication
	Username string
	// Password for database authentication
	Password string
}

func (c Config) Validate() error {
	if len(c.Hosts) == 0 {
		return fmt.Errorf("%w: at least one host is required", ErrInvalidConfig)
	}
	return nil
}
