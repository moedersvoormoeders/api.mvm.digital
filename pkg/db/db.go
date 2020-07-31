package db

import (
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var ErrorNotFound = errors.New("Not found")

type Connection struct {
	db *gorm.DB
}

type ConnectionDetails struct {
	Host     string
	Port     int
	User     string
	Database string
	Password string
}

func NewConnection() *Connection {
	conn := Connection{}

	return &conn
}

func (c *Connection) Open(details ConnectionDetails) error {
	var err error
	c.db, err = gorm.Open("postgres", fmt.Sprintf(
		"host=%s port=%d user=%s dbname=%s password=%s sslmode=disable",
		details.Host, details.Port, details.User, details.Database, details.Password))

	return err
}

func (c *Connection) Close() error {
	return c.db.Close()
}

func (c *Connection) AutoMigrate() error {
	err := c.db.AutoMigrate(&User{})
	if err != nil {
		return err.Error
	}

	return nil
}

func (c *Connection) Add(obj interface{}) (uint, error) {
	res := c.db.Create(obj)
	if res.Error != nil {
		return 0, res.Error
	}

	return res.Value.(*gorm.Model).ID, nil
}

func (c *Connection) GetID(obj interface{}, id uint) error {
	res := c.db.First(obj, id)

	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrorNotFound
	}

	return nil
}

func (c *Connection) GetWhereIs(obj interface{}, property string, where interface{}) error {
	res := c.db.First(&obj, fmt.Sprintf("%s = ?", property), where)

	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrorNotFound
	}

	return nil
}
