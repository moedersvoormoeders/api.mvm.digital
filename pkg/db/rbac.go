package db

import (
	"github.com/jinzhu/gorm"
)

type RoleVerb struct {
	gorm.Model
	Content   string `json:"content"`
	RoleRefer uint   `json:"roleRefer"`
}

type RoleEndpoint struct {
	gorm.Model
	Content   string `json:"content"`
	RoleRefer uint   `json:"roleRefer"`
}

type Role struct {
	gorm.Model
	Name      string         `json:"name"`
	Verbs     []RoleVerb     `json:"verbs" gorm:"foreignKey:RoleRefer"`
	Endpoints []RoleEndpoint `json:"endpoints" gorm:"foreignKey:RoleRefer"`
}

type RoleBinding struct {
	gorm.Model
	UserID int  `json:"userID"`
	User   User `json:"user"`
	RoleID int  `json:"roleID"`
	Role   Role `json:"role"`
}
