package zoho

import (
	"time"

	zoho "github.com/schmorrison/Zoho"
	"github.com/schmorrison/Zoho/crm"
)

// Klant defines a niceified data struct of a client
// This is to be moved from Zoho to our own DB in a later stage
type Klant struct {
	// TODO add ID
	ZohoID              string `json:"zohoID""`
	MVMNummer           string `json:"mvmNummer"`
	Naam                string `json:"naam""`
	Voornaam            string `json:"voornaam"`
	Code                string `json:"code"`
	Dag                 string `json:"dag"`
	Classificatie       string `json:"classificatie"`
	RedenControle       string `json:"redenControle"`
	Einddatum           string `json:"einddatum"` // TODO: make this time.Time post-Zoho
	TypeVoeding         string `json:"typeVoeding"`
	AantalOnder12Jaar   int    `json:"aantalOnder12Jaar"`   // TODO: make this automated post-Zoho
	AantalBovenOf12Jaar int    `json:"aantalBovenOf12Jaar"` // TODO: make this automated post-Zoho
}

// this is klant as is defined in the Zoho data
type zohoKlant struct {
	Data []struct {
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
		Kindergeld          int         `json:"Kindergeld"`
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
		Huur                      interface{} `json:"Huur"`
		Editable                  bool        `json:"$editable"`
		Postcode                  string      `json:"Postcode"`
		Code                      string      `json:"Code"`
		LeeftijdMoeder            int         `json:"Leeftijd_Moeder"`
		Soort                     string      `json:"Soort"`
		AantalVrouwen             int         `json:"Aantal_Vrouwen"`
		Instantie                 string      `json:"Instantie"`
		Aantal151                 int         `json:"Aantal_151"`
		Status                    string      `json:"$status"`
		Schulden                  int         `json:"Schulden"`
		Description               string      `json:"Description"`
		BeschBudgetbeheer         int         `json:"Besch_budgetbeheer"`
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
		MaandelijksInkomen          int           `json:"Maandelijks_inkomen"`
		ControleArmoededrempel      int           `json:"Controle_armoededrempel"`
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
		SILC                        int           `json:"SILC"`
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
	} `json:"data"`

	Info crm.PageInfo `json:"info,omitempty"`
}

func (c *CRM) GetKlantForMVMNummer(mvmNummer string) (Klant, error) {
	out, err := c.zohoCRM.SearchRecords(&zohoKlant{}, "Accounts", map[string]zoho.Parameter{
		"word":     zoho.Parameter(mvmNummer),
		"per_page": "200",
	})

	if err != nil {
		return Klant{}, err
	}

	data := out.(*zohoKlant)

	for _, entry := range data.Data {
		if entry.DoelgroepNummer == mvmNummer {
			return Klant{
				ZohoID:              entry.ID,
				MVMNummer:           entry.DoelgroepNummer,
				Naam:                entry.AccountName,
				Voornaam:            entry.Voornaam,
				Code:                entry.Code,
				Dag:                 entry.Dag,
				Classificatie:       entry.Rating,
				RedenControle:       entry.RedenControle,
				Einddatum:           entry.NieuweEvaluatie,
				TypeVoeding:         entry.Geloof, // a horrible mistake in early zoho...
				AantalOnder12Jaar:   entry.Aantal12,
				AantalBovenOf12Jaar: entry.Aantal121,
			}, nil
		}
	}

	return Klant{}, ErrNotFound
}
