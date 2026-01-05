package gocqlxfx

import (
	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx/v3"
)

func New(session *gocql.Session) gocqlx.Session {
	return gocqlx.NewSession(session)
}
