package db

import (
	"gorm.io/gorm"
)

type (
	Container struct {
		Conn *gorm.DB
	}
)

var ErrNotFound = gorm.ErrRecordNotFound
