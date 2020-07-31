package db

import (
	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	Name     string
	Username string
	Password string
}
