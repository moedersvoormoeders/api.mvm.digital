package main

import (
	"errors"
	"fmt"

	"github.com/moedersvoormoeders/api.mvm.digital/pkg/db"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
)

func init() {
	rootCmd.AddCommand(NewAddUserCmd())
}

type addUserCmdOptions struct {
	Username string
	Password string
	Name     string
}

// NewServeCmd generates the `serve` command
func NewAddUserCmd() *cobra.Command {
	a := addUserCmdOptions{}
	c := &cobra.Command{
		Use:     "add-user",
		Short:   "adds a user to the database",
		PreRunE: a.Validate,
		RunE:    a.RunE,
	}
	c.Flags().StringVarP(&a.Username, "username", "u", "", "Username for the user")
	c.Flags().StringVarP(&a.Password, "password", "p", "", "Password for the user")
	c.Flags().StringVarP(&a.Name, "name", "n", "", "Visible name fore the user")

	viper.BindPFlags(c.Flags())

	return c
}

func (a *addUserCmdOptions) Validate(cmd *cobra.Command, args []string) error {
	if a.Username == "" {
		return errors.New("need to set --username")
	}

	if a.Password == "" {
		return errors.New("need to set --password")
	}

	if a.Name == "" {
		return errors.New("need to set --name")
	}

	return nil
}

func (a *addUserCmdOptions) RunE(cmd *cobra.Command, args []string) error {
	dbConn := db.NewConnection()
	err := dbConn.Open(db.ConnectionDetails{
		Host:     "postgres",
		Port:     5432,
		User:     "postgres",
		Database: "postgres",
		Password: "moedersvoormoeders", //TODO: make flags
	})

	if err != nil {
		return fmt.Errorf("error opening database: %w", err)
	}
	defer dbConn.Close()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(a.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("error hashing password: %w", err)
	}

	err = dbConn.Add(&db.User{
		Name:     a.Name,
		Username: a.Username,
		Password: string(hashedPassword),
	})

	return err
}
