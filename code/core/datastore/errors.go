package datastore

import (
	"github.com/MurashovVen/bandwidth-marketplace/code/core/errors"
)

var (
	// DBOpenError represent common.Error occurs while opening the db.
	DBOpenError = errors.NewError("db_open_error", "Error opening the DB connection")
)
