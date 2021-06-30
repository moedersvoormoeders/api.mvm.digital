package main

import (
	"fmt"

	"github.com/moedersvoormoeders/api.mvm.digital/pkg/db"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(NewVoedingCmd())
}

type voedingCmdOptions struct {
	postgresHost     string
	postgresPort     int
	postgresUsername string
	postgresDatabase string
	postgresPassword string
}

func NewVoedingCmd() *cobra.Command {
	a := voedingCmdOptions{}
	c := &cobra.Command{
		Use:     "voeding",
		Short:   "adds voeding to the database",
		PreRunE: a.Validate,
		RunE:    a.RunE,
	}
	c.Flags().StringVar(&a.postgresHost, "postgres-host", "", "PostgreSQL hostname")
	c.Flags().IntVar(&a.postgresPort, "postgres-port", 5432, "PostgreSQL hostname")
	c.Flags().StringVar(&a.postgresUsername, "postgres-username", "", "PostgreSQL hostname")
	c.Flags().StringVar(&a.postgresPassword, "postgres-password", "", "PostgreSQL hostname")
	c.Flags().StringVar(&a.postgresDatabase, "postgres-database", "", "PostgreSQL hostname")

	viper.BindPFlags(c.Flags())

	return c
}

func (a *voedingCmdOptions) Validate(cmd *cobra.Command, args []string) error {
	return nil
}

func (a *voedingCmdOptions) RunE(cmd *cobra.Command, args []string) error {
	// TODO: make function to actually be flexible

	dbConn, err := db.NewConnection(db.ConnectionDetails{
		Host:     a.postgresHost,
		Port:     a.postgresPort,
		User:     a.postgresUsername,
		Database: a.postgresDatabase,
		Password: a.postgresPassword,
	})

	dbConn.DoMigrate()

	if err != nil {
		return fmt.Errorf("error opening database: %w", err)
	}

	toAddCategories := []db.Categorie{
		{
			Naam:       "Voeding",
			PerPersoon: false,
			OpMaat:     false,
			Order:      0,
			Afdeling:   "Voeding",
		},
	}

	for _, obj := range toAddCategories {
		fmt.Println(obj)
		err = dbConn.Add(&obj)
		if err != nil {
			return err
		}
	}

	catVoeding := db.Categorie{}
	dbConn.GetWhereIs(&catVoeding, "naam", "Voeding")

	objectsToAdd := []db.Object{
		db.Object{
			Naam:      "Voeding",
			Categorie: catVoeding,
		},
		db.Object{
			Naam:      "Vuilzakken",
			Categorie: catVoeding,
		},
		db.Object{
			Naam:      "Verjaardag",
			Categorie: catVoeding,
		},
		db.Object{
			Naam:      "Snoep Sinterklaas",
			Categorie: catVoeding,
		},
		db.Object{
			Naam:      "Paaspakket",
			Categorie: catVoeding,
		},

		db.Object{
			Naam:      "Melkpoeder Nan 1",
			Categorie: catVoeding,
		},
		db.Object{
			Naam:      "Melkpoeder Nan 2",
			Categorie: catVoeding,
		},
		db.Object{
			Naam:      "Melkpoeder Nutrilon 1",
			Categorie: catVoeding,
		},
		db.Object{
			Naam:      "Melkpoeder Nutrilon 2",
			Categorie: catVoeding,
		},
		db.Object{
			Naam:      "Melkpoeder Novalec 2",
			Categorie: catVoeding,
		},
		db.Object{
			Naam:      "Melkpoeder Anderen",
			Categorie: catVoeding,
		},
	}

	for _, obj := range objectsToAdd {
		err = dbConn.Add(&obj)
		if err != nil {
			return err
		}
	}

	return err
}
