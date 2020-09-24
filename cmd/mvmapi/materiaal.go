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

	dbConn, err := db.NewConnection(db.ConnectionDetails{
		Host:     a.postgresHost,
		Port:     a.postgresPort,
		User:     a.postgresUsername,
		Database: a.postgresDatabase,
		Password: a.postgresPassword,
	})

	dbConn.DoMigrate()

	badbyMaten := []db.MateriaalMaat{
		{Naam: "prematuur"},
		{Naam: "0 ma - 56"},
		{Naam: "3 ma - 62"},
		{Naam: "6 ma - 68"},
		{Naam: "9 ma - 74"},
		{Naam: "12 ma - 80"},
		{Naam: "18 ma - 86"},
	}

	defaultMaten := []db.MateriaalMaat{
		{Naam: "prematuur"},
		{Naam: "0 ma - 56"},
		{Naam: "3 ma - 62"},
		{Naam: "6 ma - 68"},
		{Naam: "9 ma - 74"},
		{Naam: "12 ma - 80"},
		{Naam: "18 ma - 86"},
		{Naam: "2 jr - 92"},
		{Naam: "3 jr - 98"},
		{Naam: "4 jr - 104"},
		{Naam: "5 jr - 110"},
		{Naam: "6 jr - 116"},
		{Naam: "7 jr - 122"},
		{Naam: "8 jr - 128"},
		{Naam: "9 jr - 134"},
		{Naam: "10 jr - 140"},
		{Naam: "11 jr - 146"},
		{Naam: "12 jr - 152"},
		{Naam: "14 jr - 164"},
	}

	schoenMaten := []db.MateriaalMaat{
		{Naam: "baby"},
		{Naam: "22"},
		{Naam: "23"},
		{Naam: "24"},
		{Naam: "25"},
		{Naam: "26"},
		{Naam: "27"},
		{Naam: "28"},
		{Naam: "29"},
		{Naam: "30"},
		{Naam: "31"},
		{Naam: "32"},
		{Naam: "33"},
		{Naam: "34"},
		{Naam: "35"},
		{Naam: "36"},
		{Naam: "37"},
		{Naam: "38"},
		{Naam: "39"},
		{Naam: "40"},
		{Naam: "41"},
		{Naam: "42"},
		{Naam: "43"},
		{Naam: "44"},
		{Naam: "45"},
		{Naam: ">45"},
	}

	if err != nil {
		return fmt.Errorf("error opening database: %w", err)
	}

	toAddCategories := []db.MateriaalCategory{
		db.MateriaalCategory{
			Naam:    "Kinderkleding",
			OpMaat:  true,
			PerKind: true,
			Order:   1,
		},
	}

	for _, obj := range toAddCategories {
		fmt.Println(obj)
		err = dbConn.Add(&obj)
		if err != nil {
			return err
		}
	}

	catKleding := db.MateriaalCategory{}
	dbConn.GetWhereIs(&catKleding, "naam", "Kinderkleding")
	/*catSpeelgoed := db.MateriaalCategory{}
	dbConn.GetWhereIs(&catSpeelgoed, "naam", "Speelgoed")
	catBabymateriaal := db.MateriaalCategory{}
	dbConn.GetWhereIs(&catBabymateriaal, "naam", "Babymateriaal")
	catVoorMoeder := db.MateriaalCategory{}
	dbConn.GetWhereIs(&catVoorMoeder, "naam", "Voor Moeder")*/

	objectsToAdd := []db.MateriaalObject{
		db.MateriaalObject{
			Naam:      "Pakket Zomer",
			Categorie: catKleding,
			Maten:     copySlice(defaultMaten),
		},
		db.MateriaalObject{
			Naam:      "Pakket Winter",
			Categorie: catKleding,
			Maten:     copySlice(defaultMaten),
		},
		db.MateriaalObject{
			Naam:      "Pakket Lente",
			Categorie: catKleding,
			Maten:     copySlice(badbyMaten),
		},
		db.MateriaalObject{
			Naam:      "Pakket Herfst",
			Categorie: catKleding,
			Maten:     copySlice(badbyMaten),
		},
		db.MateriaalObject{
			Naam:      "Schoenen Zomer",
			Categorie: catKleding,
			Maten:     copySlice(schoenMaten),
		},
		db.MateriaalObject{
			Naam:      "Schoenen Winter",
			Categorie: catKleding,
			Maten:     copySlice(schoenMaten),
		},
		db.MateriaalObject{
			Naam:      "Pantoffels winter",
			Categorie: catKleding,
			Maten:     copySlice(schoenMaten),
		},
		db.MateriaalObject{
			Naam:      "Laarzen winter",
			Categorie: catKleding,
			Maten:     copySlice(schoenMaten),
		},
		db.MateriaalObject{
			Naam:      "Laarzen regen",
			Categorie: catKleding,
			Maten:     copySlice(schoenMaten),
		},
		db.MateriaalObject{
			Naam:      "Turnpantoffels",
			Categorie: catKleding,
			Maten:     copySlice(schoenMaten),
		},
		db.MateriaalObject{
			Naam:      "Uniform",
			Categorie: catKleding,
			Maten:     copySlice(defaultMaten),
		},
		db.MateriaalObject{
			Naam:      "Feestkleding",
			Categorie: catKleding,
			Maten:     copySlice(defaultMaten),
		},
		db.MateriaalObject{
			Naam:      "Geboortepakket",
			Categorie: catKleding,
			Maten:     []db.MateriaalMaat{{Naam: "baby"}},
		},
		db.MateriaalObject{
			Naam:      "Extra",
			Categorie: catKleding,
			Maten:     copySlice(defaultMaten),
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

func copySlice(in []db.MateriaalMaat) []db.MateriaalMaat {
	out := []db.MateriaalMaat(nil)
	for _, maat := range in {
		out = append(out, maat)
	}
	return out
}
