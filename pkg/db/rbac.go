package db

import (
	"github.com/jinzhu/gorm"
)

type RoleVerb struct {
	gorm.Model
	Content   string
	RoleRefer uint
}

type RoleEndpoint struct {
	gorm.Model
	Content   string
	RoleRefer uint
}

type Role struct {
	gorm.Model
	Name      string
	Verbs     []RoleVerb     `gorm:"foreignKey:RoleRefer"`
	Endpoints []RoleEndpoint `gorm:"foreignKey:RoleRefer"`
}

type RoleBinding struct {
	gorm.Model
	UserID int
	User   User
	RoleID int
	Role   Role
}
