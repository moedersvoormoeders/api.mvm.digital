package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/davecgh/go-spew/spew"

	zoho "github.com/schmorrison/Zoho"
	"github.com/schmorrison/Zoho/crm"

	"github.com/moedersvoormoeders/api.mvm.digital/pkg/db"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(NewVoedingImportCmd())
}

type voedingImportCmdOptions struct {
	postgresHost     string
	postgresPort     int
	postgresUsername string
	postgresDatabase string
	postgresPassword string

	zohoClientID     string
	zohoClientSecret string
}

func NewVoedingImportCmd() *cobra.Command {
	a := voedingImportCmdOptions{}
	c := &cobra.Command{
		Use:     "import-voeding",
		Short:   "imports voeding data fronm ZOHO",
		PreRunE: a.Validate,
		RunE:    a.RunE,
	}
	c.Flags().StringVar(&a.postgresHost, "postgres-host", "", "PostgreSQL hostname")
	c.Flags().IntVar(&a.postgresPort, "postgres-port", 5432, "PostgreSQL hostname")
	c.Flags().StringVar(&a.postgresUsername, "postgres-username", "", "PostgreSQL hostname")
	c.Flags().StringVar(&a.postgresPassword, "postgres-password", "", "PostgreSQL hostname")
	c.Flags().StringVar(&a.postgresDatabase, "postgres-database", "", "PostgreSQL hostname")

	// needed for Zoho connector
	c.Flags().StringVar(&a.zohoClientID, "zoho-clientid", "", "Zoho client ID, you can get this at https://accounts.zoho.eu/developerconsole")
	c.Flags().StringVar(&a.zohoClientSecret, "zoho-clientsecret", "", "Zoho client Secret, you can get this at https://accounts.zoho.eu/developerconsole")

	viper.BindPFlags(c.Flags())

	return c
}

func (a *voedingImportCmdOptions) Validate(cmd *cobra.Command, args []string) error {
	return nil
}

// The code you are about to see is single use throw away code that we only run once ever!
// it is on the lowest possible standards to save natural resources
func (a *voedingImportCmdOptions) RunE(cmd *cobra.Command, args []string) error {
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

	z := zoho.New()
	z.SetZohoTLD("eu")

	// to start oAuth2 flow
	scopes := []zoho.ScopeString{
		zoho.BuildScope(zoho.Crm, zoho.ModulesScope, zoho.AllMethod, zoho.NoOp),
	}

	if err := z.AuthorizationCodeRequest(a.zohoClientID, a.zohoClientSecret, scopes, "http://localhost:8080/oauthredirect"); err != nil {
		log.Fatal(err)
	}

	c := crm.New(z)

	go func() {
		for {
			time.Sleep(time.Minute)
			err := z.RefreshTokenRequest()
			if err != nil {
				fmt.Println(err)
			}
		}
	}()

	// The API for getting module records is bound to change once the returned data types can be defined.
	// The returned JSON values are subject to change given that custom fields are an instrinsic part of zoho. (see brainstorm above)

	continueScan := true
	page := 0
	for continueScan {
		out, err := c.ListRecords(&zohoKlant{}, crm.AccountsModule, map[string]zoho.Parameter{
			"page": zoho.Parameter(fmt.Sprintf("%d", page)),
		})
		if err != nil {
			log.Fatal(err)
		}
		data := out.(*zohoKlant)

		for _, entry := range data.Data {
			out, err := c.SearchRecords(&zohoVoeding{}, "Voeding", map[string]zoho.Parameter{
				"word":     zoho.Parameter(entry.DoelgroepNummer),
				"per_page": "200",
			})

			if err != nil {
				log.Println(entry.DoelgroepNummer, err)
				continue
			}

			data := out.(*zohoVoeding)
			dataID := ""
			for _, vEntry := range data.Data {
				if vEntry.DoelgroepNummer == entry.DoelgroepNummer {
					dataID = vEntry.ID
				}
			}

			if dataID == "" {
				log.Println(entry.DoelgroepNummer, "No dataID")
				continue
			}

			// now the same again but with the ID so we get all data
			out, err = c.GetRecord(&zohoVoeding{}, "Voeding", dataID)

			if err != nil {
				spew.Dump(out)
				log.Println(entry.DoelgroepNummer, err)
				panic(err)
				continue
			}

			data = out.(*zohoVoeding)
			for _, ventry := range data.Data {

				if ventry.DoelgroepNummer == entry.DoelgroepNummer {
					voeding := db.Voeding{}
					dbConn.Where(&db.Voeding{MVMNummer: entry.DoelgroepNummer}).First(&voeding)

					voeding.MVMNummer = ventry.DoelgroepNummer
					voeding.SpecialeVoeding = ventry.SpecialeVoeding
					voeding.Opmerking = ventry.AlgemeneOpmerkingen

					if voeding.ID == 0 {
						res := dbConn.Create(&voeding)
						if res.Error != nil {
							log.Println(err)
							continue
						}
					} else {
						res := dbConn.Updates(&voeding)
						if res.Error != nil {
							log.Println(err)
							continue
						}
					}

					for _, pakket := range ventry.Geschiedenis {
						for _, gekregen := range pakket.Gekregen {
							if strings.ToLower(gekregen) == "snoep sinterklaas" {
								gekregen = "Snoep Sinterklaas"
							}
							if gekregen == "Melkpoeder Nutrilon" {
								gekregen = "Melkpoeder Nutrilon 1"
							}
							if gekregen == "Melkpoeder Nan" {
								gekregen = "Melkpoeder Nan 1"
							}
							obj := db.Object{}
							res := dbConn.Where(&db.Object{Naam: strings.Trim(gekregen, " ")}).First(&obj)
							if res.Error != nil {
								fmt.Println(gekregen)
								panic(res.Error)
							}

							time, _ := time.Parse("2006-01-02", pakket.Datum)
							res = dbConn.Create(&db.OntvangEntry{
								Datum:     time,
								Aantal:    1,
								ObjectID:  int(obj.ID),
								Opmerking: pakket.Opmerking,
								VoedingID: int(voeding.ID),
							})
							if res.Error != nil {
								log.Println(err)
								continue
							}
						}
					}
				}
			}
			fmt.Printf("Moved %s\n", entry.DoelgroepNummer)
			time.Sleep(1 * time.Second)
		}
		if len(data.Data) < 200 {
			continueScan = false
		}
		page++
	}

	return err
}

type zohoVoeding struct {
	Data []struct {
		Owner struct {
			Name  string `json:"name"`
			ID    string `json:"id"`
			Email string `json:"email"`
		} `json:"Owner"`
		CurrencySymbol string      `json:"$currency_symbol"`
		PhotoID        interface{} `json:"$photo_id"`
		Voornaam       string      `json:"Voornaam"`
		ReviewProcess  struct {
			Approve  bool `json:"approve"`
			Reject   bool `json:"reject"`
			Resubmit bool `json:"resubmit"`
		} `json:"$review_process"`
		UpcomingActivity interface{} `json:"$upcoming_activity"`
		Name             string      `json:"Name"`
		LastActivityTime time.Time   `json:"Last_Activity_Time"`
		ModifiedBy       struct {
			Name  string `json:"name"`
			ID    string `json:"id"`
			Email string `json:"email"`
		} `json:"Modified_By"`
		Review            interface{} `json:"$review"`
		AantalVolwassenen int         `json:"Aantal_Volwassenen"`
		State             string      `json:"$state"`
		ProcessFlow       bool        `json:"$process_flow"`
		Klant             struct {
			Module string `json:"module"`
			Name   string `json:"name"`
			ID     string `json:"id"`
		} `json:"Klant"`
		ID       string `json:"id"`
		Naam     string `json:"Naam"`
		Approved bool   `json:"$approved"`
		Approval struct {
			Delegate bool `json:"delegate"`
			Approve  bool `json:"approve"`
			Reject   bool `json:"reject"`
			Resubmit bool `json:"resubmit"`
		} `json:"$approval"`
		Eenmalige           interface{}   `json:"Eenmalige"`
		ModifiedTime        time.Time     `json:"Modified_Time"`
		CreatedTime         time.Time     `json:"Created_Time"`
		Editable            bool          `json:"$editable"`
		Geboortedata        string        `json:"Geboortedata"`
		Postcode            string        `json:"Postcode"`
		Code                string        `json:"Code"`
		Orchestration       bool          `json:"$orchestration"`
		DoelgroepNummer     string        `json:"Doelgroep_Nummer"`
		InMerge             bool          `json:"$in_merge"`
		Status              string        `json:"$status"`
		AantalKinderen      int           `json:"Aantal_Kinderen"`
		SoortVoeding        string        `json:"Soort_voeding"`
		Tag                 []interface{} `json:"Tag"`
		SpecialeVoeding     string        `json:"Speciale_voeding"`
		ApprovalState       string        `json:"$approval_state"`
		AlgemeneOpmerkingen string        `json:"Algemene_Opmerkingen"`
		KoelzakGekregen     bool          `json:"Koelzak_gekregen"`
		Geschiedenis        []struct {
			Datum    string   `json:"Datum"`
			Gekregen []string `json:"Gekregen"`
			InMerge  bool     `json:"$in_merge"`
			ParentID struct {
				Module string `json:"module"`
				Name   string `json:"name"`
				ID     string `json:"id"`
			} `json:"Parent_Id"`
			ID            string `json:"id"`
			Opmerking     string `json:"Opmerking"`
			Orchestration bool   `json:"$orchestration"`
		} `json:"Geschiedenis"`
	} `json:"data"`
	Info crm.PageInfo `json:"info,omitempty"`
}

// this is klant as is defined in the Zoho data
type zohoKlant struct {
	Data []zohoKlantData `json:"data"`

	Info crm.PageInfo `json:"info,omitempty"`
}

type zohoKlantData struct {
	Doorverwijzingsbrief bool `json:"doorverwijzingsbrief"`
	Owner                struct {
		Name  string `json:"name"`
		ID    string `json:"id"`
		Email string `json:"email"`
	} `json:"Owner"`
	CurrencySymbol      string      `json:"$currency_symbol"`
	Datum1              string      `json:"Datum_1"`
	Gezinstatus         string      `json:"Gezinstatus"`
	LastActivityTime    time.Time   `json:"Last_Activity_Time"`
	State               string      `json:"$state"`
	ProcessFlow         bool        `json:"$process_flow"`
	Aantal65            int         `json:"Aantal_65"`
	Straat              string      `json:"Straat"`
	ID                  string      `json:"id"`
	Nederlands          string      `json:"Nederlands"`
	Approved            bool        `json:"$approved"`
	Datum1EInschrijving string      `json:"Datum_1e_inschrijving"`
	Geslacht            string      `json:"Geslacht"`
	MaandelijkseKosten  interface{} `json:"Maandelijkse_kosten"`
	Kindergeld          float64     `json:"Kindergeld"`
	Approval            struct {
		Delegate bool `json:"delegate"`
		Approve  bool `json:"approve"`
		Reject   bool `json:"reject"`
		Resubmit bool `json:"resubmit"`
	} `json:"$approval"`
	AttestGezinssamenstelling bool        `json:"Attest_gezinssamenstelling"`
	CreatedTime               time.Time   `json:"Created_Time"`
	DatumHerinschrijving      interface{} `json:"Datum_herinschrijving"`
	Nationaliteit1            []string    `json:"Nationaliteit1"`
	Huur                      float64     `json:"Huur"`
	Editable                  bool        `json:"$editable"`
	Postcode                  string      `json:"Postcode"`
	Code                      string      `json:"Code"`
	LeeftijdMoeder            int         `json:"Leeftijd_Moeder"`
	Soort                     string      `json:"Soort"`
	AantalVrouwen             int         `json:"Aantal_Vrouwen"`
	Instantie                 string      `json:"Instantie"`
	Aantal151                 int         `json:"Aantal_151"`
	Status                    string      `json:"$status"`
	Schulden                  float64     `json:"Schulden"`
	Description               string      `json:"Description"`
	BeschBudgetbeheer         float64     `json:"Besch_budgetbeheer"`
	PhotoID                   interface{} `json:"$photo_id"`
	Voornaam                  string      `json:"Voornaam"`
	Rating                    string      `json:"Rating"`
	DoelgroepNummer           string      `json:"Doelgroep_nummer"`
	ReviewProcess             struct {
		Approve  bool `json:"approve"`
		Reject   bool `json:"reject"`
		Resubmit bool `json:"resubmit"`
	} `json:"$review_process"`
	EMail1            string      `json:"E_mail_1"`
	Aantal152         int         `json:"Aantal_152"`
	GeboorteLandNieuw string      `json:"Geboorte_Land_nieuw"`
	NieuweEvaluatie   string      `json:"Nieuwe_evaluatie"`
	RecordImage       interface{} `json:"Record_Image"`
	ModifiedBy        struct {
		Name  string `json:"name"`
		ID    string `json:"id"`
		Email string `json:"email"`
	} `json:"Modified_By"`
	Review                      interface{}   `json:"$review"`
	RedenControle               string        `json:"Reden_Controle"`
	Armoedefactor               float64       `json:"Armoedefactor"`
	Phone                       string        `json:"Phone"`
	MaandelijksInkomen          float64       `json:"Maandelijks_inkomen"`
	ControleArmoededrempel      float64       `json:"Controle_armoededrempel"`
	AccountName                 string        `json:"Account_Name"`
	OverigeInformatie           interface{}   `json:"Overige_informatie"`
	ModifiedTime                time.Time     `json:"Modified_Time"`
	Huisnummer                  string        `json:"Huisnummer"`
	AflossingMaand              interface{}   `json:"Aflossing_maand"`
	Periode                     string        `json:"Periode"`
	Reden                       string        `json:"Reden"`
	Aantal122                   int           `json:"Aantal_122"`
	Aantal15                    int           `json:"Aantal_15"`
	Aantal121                   int           `json:"Aantal_121"`
	GasElectriciteit            interface{}   `json:"Gas_Electriciteit"`
	Rijksregisternummer         string        `json:"Rijksregisternummer"`
	SILC                        float64       `json:"SILC"`
	Dag                         string        `json:"Dag"`
	Akkoordverklaring           bool          `json:"Akkoordverklaring"`
	DoelgroepVestiging          string        `json:"Doelgroep_Vestiging"`
	Orchestration               bool          `json:"$orchestration"`
	ParentAccount               interface{}   `json:"Parent_Account"`
	Arbeidsstatus               string        `json:"Arbeidsstatus"`
	Geloof                      string        `json:"Geloof"`
	Drempel                     float64       `json:"Drempel"`
	AantalPersonen              int           `json:"Aantal_Personen"`
	Aantal12                    int           `json:"Aantal_12"`
	InMerge                     bool          `json:"$in_merge"`
	Gemeente                    string        `json:"Gemeente"`
	BijkomendeInfoDoorverwijzer interface{}   `json:"Bijkomende_info_doorverwijzer"`
	AantalBegunstigden          int           `json:"aantal_begunstigden"`
	OudNummer                   interface{}   `json:"Oud_nummer"`
	Tag                         []interface{} `json:"Tag"`
	ApprovalState               string        `json:"$approval_state"`
}
