package repository

import (
	"github.com/driver005/database"
)

var (
	// ErrNotFound is a convenience reference for the actual GORM error
	ErrNotFound = database.ErrRecordNotFound
)
