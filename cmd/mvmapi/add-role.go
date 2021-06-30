package main

import (
	"errors"

	"github.com/moedersvoormoeders/api.mvm.digital/pkg/db"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(NewAddRoleCmd())
}

type addRoleCmdOptions struct {
	Name      string
	Verbs     []string
	Endpoints []string

	postgresHost     string
	postgresPort     int
	postgresUsername string
	postgresDatabase string
	postgresPassword string
}

// NewAddRoleCmd generates the `add-role` command
func NewAddRoleCmd() *cobra.Command {
	a := addRoleCmdOptions{}
	c := &cobra.Command{
		Use:     "add-role",
		Short:   "adds a role to the database",
		PreRunE: a.Validate,
		RunE:    a.RunE,
	}
	c.Flags().StringVar(&a.postgresHost, "postgres-host", "", "PostgreSQL hostname")
	c.Flags().IntVar(&a.postgresPort, "postgres-port", 5432, "PostgreSQL hostname")
	c.Flags().StringVar(&a.postgresUsername, "postgres-username", "", "PostgreSQL hostname")
	c.Flags().StringVar(&a.postgresPassword, "postgres-password", "", "PostgreSQL hostname")
	c.Flags().StringVar(&a.postgresDatabase, "postgres-database", "", "PostgreSQL hostname")

	c.Flags().StringVar(&a.Name, "name", "", "Role name")
	c.Flags().StringSliceVar(&a.Verbs, "verbs", nil, "Role allowed verbs")
	c.Flags().StringSliceVar(&a.Endpoints, "endpoints", nil, "Role allowed endpoints")

	viper.BindPFlags(c.Flags())

	return c
}

func (a *addRoleCmdOptions) Validate(cmd *cobra.Command, args []string) error {
	if a.Name == "" || len(a.Endpoints) == 0 || len(a.Verbs) == 0 {
		return errors.New("--name, --endpoints and --verbs are required")
	}
	return nil
}

func (a *addRoleCmdOptions) RunE(cmd *cobra.Command, args []string) error {
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

	role := db.Role{
		Name: a.Name,
	}

	res := dbConn.Create(&role)
	if res.Error != nil {
		return err
	}

	for _, verb := range a.Verbs {
		res := dbConn.Create(&db.RoleVerb{
			Content:   verb,
			RoleRefer: role.ID,
		})
		if res.Error != nil {
			return err
		}
	}

	for _, endpoint := range a.Endpoints {
		res := dbConn.Create(&db.RoleEndpoint{
			Content:   endpoint,
			RoleRefer: role.ID,
		})
		if res.Error != nil {
			return err
		}
	}

	return res.Error
}
