package gocqlfx

import (
	"fmt"

	"github.com/gocql/gocql"
)

func New(config Config) (*gocql.Session, error) {
	cluster := gocql.NewCluster(config.Hosts...)
	cluster.Keyspace = config.Keyspace
	cluster.Consistency = gocql.Quorum
	cluster.Authenticator = gocql.PasswordAuthenticator{
		Username:              config.Username,
		Password:              config.Password,
		AllowedAuthenticators: []string{},
	}

	s, err := cluster.CreateSession()
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	return s, nil
}
