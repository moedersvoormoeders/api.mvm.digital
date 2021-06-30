package main

import (
	"errors"

	"github.com/moedersvoormoeders/api.mvm.digital/pkg/db"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(NewAddRoleBindingCmd())
}

type addRoleBindingCmdOptions struct {
	User string
	Role string

	postgresHost     string
	postgresPort     int
	postgresUsername string
	postgresDatabase string
	postgresPassword string
}

// NewAddRoleBindingCmd generates the `add-rolebinding` command
func NewAddRoleBindingCmd() *cobra.Command {
	a := addRoleBindingCmdOptions{}
	c := &cobra.Command{
		Use:     "add-rolebinding",
		Short:   "adds a rolebinding to the database",
		PreRunE: a.Validate,
		RunE:    a.RunE,
	}
	c.Flags().StringVar(&a.postgresHost, "postgres-host", "", "PostgreSQL hostname")
	c.Flags().IntVar(&a.postgresPort, "postgres-port", 5432, "PostgreSQL hostname")
	c.Flags().StringVar(&a.postgresUsername, "postgres-username", "", "PostgreSQL hostname")
	c.Flags().StringVar(&a.postgresPassword, "postgres-password", "", "PostgreSQL hostname")
	c.Flags().StringVar(&a.postgresDatabase, "postgres-database", "", "PostgreSQL hostname")

	c.Flags().StringVar(&a.User, "user", "", "User name")
	c.Flags().StringVar(&a.Role, "role", "", "Role name")

	viper.BindPFlags(c.Flags())

	return c
}

func (a *addRoleBindingCmdOptions) Validate(cmd *cobra.Command, args []string) error {
	if a.User == "" || a.Role == "" {
		return errors.New("--user and --name are required")
	}
	return nil
}

func (a *addRoleBindingCmdOptions) RunE(cmd *cobra.Command, args []string) error {
	dbConn, err := db.NewConnection(db.ConnectionDetails{
		Host:     a.postgresHost,
		Port:     a.postgresPort,
		User:     a.postgresUsername,
		Database: a.postgresDatabase,
		Password: a.postgresPassword,
	})

	if err != nil {
		return err
	}

	user := db.User{}
	res := dbConn.Where(db.User{Name: a.User}).Find(&user)
	if res.Error != nil {
		return res.Error
	} else if res.RowsAffected == 0 {
		return errors.New("No user found")
	}

	role := db.Role{}
	res = dbConn.Where(db.Role{Name: a.Role}).Find(&role)
	if res.Error != nil {
		return res.Error
	} else if res.RowsAffected == 0 {
		return errors.New("No role found")
	}

	roleb := db.RoleBinding{
		User: user,
		Role: role,
	}

	res = dbConn.Create(&roleb)

	return res.Error
}
