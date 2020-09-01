package main

import (
	"fmt"

	"github.com/moedersvoormoeders/api.mvm.digital/pkg/db"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(NewMateriaalCmd())
}

type materiaalCmdOptions struct {
	Username string
	Password string
	Name     string

	postgresHost     string
	postgresPort     int
	postgresUsername string
	postgresDatabase string
	postgresPassword string
}

// NewServeCmd generates the `serve` command
func NewMateriaalCmd() *cobra.Command {
	a := materiaalCmdOptions{}
	c := &cobra.Command{
		Use:     "materiaal",
		Short:   "adds materiaal to the database",
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

func (a *materiaalCmdOptions) Validate(cmd *cobra.Command, args []string) error {
	return nil
}

func (a *materiaalCmdOptions) RunE(cmd *cobra.Command, args []string) error {
	// TODO: make function to actually be flexible

	dbConn := db.NewConnection()
	err := dbConn.Open(db.ConnectionDetails{
		Host:     a.postgresHost,
		Port:     a.postgresPort,
		User:     a.postgresUsername,
		Database: a.postgresDatabase,
		Password: a.postgresPassword,
	})

	if err != nil {
		return fmt.Errorf("error opening database: %w", err)
	}
	defer dbConn.Close()

	toAddCategories := []interface{}{
		db.MateriaalCategory{
			Naam:    "Kleding",
			OpMaat:  true,
			PerKind: true,
			Order:   1,
		},
		db.MateriaalCategory{
			Naam:    "Speelgoed",
			OpMaat:  false,
			PerKind: true,
			Order:   2,
		},
		db.MateriaalCategory{
			Naam:    "Babymateriaal",
			OpMaat:  false,
			PerKind: false,
			Order:   3,
		},
		db.MateriaalCategory{
			Naam:    "Voor Moeder",
			OpMaat:  false,
			PerKind: false,
			Order:   4,
		},
	}

	for _, obj := range toAddCategories {
		err = dbConn.Add(obj)
		if err != nil {
			return err
		}
	}

	catKleding := db.MateriaalCategory{}
	dbConn.GetWhereIs(&catKleding, "naam", "Kleding")
	catSpeelgoed := db.MateriaalCategory{}
	dbConn.GetWhereIs(&catKleding, "naam", "Speelgoed")
	catBabymateriaal := db.MateriaalCategory{}
	dbConn.GetWhereIs(&catKleding, "naam", "Babymateriaal")
	catVoorMoeder := db.MateriaalCategory{}
	dbConn.GetWhereIs(&catKleding, "naam", "Voor Moeder")

	objectsToAdd := []interface{}{
		db.MateriaalObject{
			Naam:      "Pakket Zomer",
			Categorie: catKleding,
		},
		db.MateriaalObject{
			Naam:      "Pakket Winter",
			Categorie: catKleding,
		},
		db.MateriaalObject{
			Naam:      "Schoenen Zomer",
			Categorie: catKleding,
		},
		db.MateriaalObject{
			Naam:      "Schoenen Winter",
			Categorie: catKleding,
		},
		db.MateriaalObject{
			Naam:      "Uniform",
			Categorie: catKleding,
		},
		db.MateriaalObject{
			Naam:      "Verjaardag",
			Categorie: catSpeelgoed,
		},
		db.MateriaalObject{
			Naam:      "Mini cadeau",
			Categorie: catSpeelgoed,
		},
		db.MateriaalObject{
			Naam:      "Carnaval",
			Categorie: catSpeelgoed,
		},
		db.MateriaalObject{
			Naam:      "Extra",
			Categorie: catSpeelgoed,
		},
		db.MateriaalObject{
			Naam:      "Ziekenhuispakket",
			Categorie: catVoorMoeder,
		},
		db.MateriaalObject{
			Naam:      "Kindskorf",
			Categorie: catVoorMoeder,
		},
		db.MateriaalObject{
			Naam:      "Zwangerschapskleding",
			Categorie: catVoorMoeder,
		},
		db.MateriaalObject{
			Naam:      "Gelegenheids outfit",
			Categorie: catVoorMoeder,
		},
		db.MateriaalObject{
			Naam:      "Winterjas",
			Categorie: catVoorMoeder,
		},
		db.MateriaalObject{
			Naam:      "Schoenen",
			Categorie: catVoorMoeder,
		},
		db.MateriaalObject{
			Naam:      "Kapsalon",
			Categorie: catVoorMoeder,
		},
		db.MateriaalObject{
			Naam:      "Buggy",
			Categorie: catBabymateriaal,
		},
	}

	for _, obj := range objectsToAdd {
		err = dbConn.Add(obj)
		if err != nil {
			return err
		}
	}

	return err
}
