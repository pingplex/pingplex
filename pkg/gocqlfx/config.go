package gocqlfx

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
