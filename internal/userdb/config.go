package userdb

import "time"

type Config struct {
	URLTemplate string

	ConnMaxIdleTime time.Duration
	ConnMaxLifetime time.Duration
	MaxOpenConns    int
	MaxIdleConns    int
}
