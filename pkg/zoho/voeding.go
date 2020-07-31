package zoho

import (
	"time"

	zoho "github.com/schmorrison/Zoho"
	"github.com/schmorrison/Zoho/crm"
)

// Voeding defines a niceified and incomplete data struct of a client
// This is to be moved from Zoho to our own DB in a later stage
type Voeding struct {
	// TODO add ID
	ZohoID          string          `json:"zohoID""`
	MVMNummer       string          `json:"mvmNummer"`
	SpecialeVoeding string          `json:"specialeVoeding"`
	Opmerking       string          `json:"opmerking"`
	Pakketten       []VoedingPakket `json:"paketten"`
	Geboortedata    string          `json:"geboortedata"` // TODO: delete post-Zoho
}

type VoedingPakket struct {
	ID        string   `json:"id"`
	Datum     string   `json:"datum"` //TODO: make time.Time post-Zoho
	Gekregen  []string `json:"gekregen"`
	Opmerking string   `json:"opmerking"`
}

// this is klant as is defined in the Zoho data
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

func (c *CRM) GetVoedingForMVMNummer(mvmNummer string) (Voeding, error) {
	out, err := c.zohoCRM.SearchRecords(&zohoVoeding{}, "Voeding", map[string]zoho.Parameter{
		"word":     zoho.Parameter(mvmNummer),
		"per_page": "200",
	})

	if err != nil {
		return Voeding{}, err
	}

	data := out.(*zohoVoeding)
	dataID := ""
	for _, entry := range data.Data {
		if entry.DoelgroepNummer == mvmNummer {
			dataID = entry.ID
		}
	}

	if dataID == "" {
		return Voeding{}, ErrNotFound
	}

	// now the same again but with the ID so we get all data
	out, err = c.zohoCRM.GetRecord(&zohoVoeding{}, "Voeding", dataID)

	if err != nil {
		return Voeding{}, err
	}

	data = out.(*zohoVoeding)
	for _, entry := range data.Data {
		if entry.DoelgroepNummer == mvmNummer {
			voeding := Voeding{
				ZohoID:          entry.ID,
				MVMNummer:       entry.DoelgroepNummer,
				SpecialeVoeding: entry.SpecialeVoeding,
				Opmerking:       entry.AlgemeneOpmerkingen,
				Geboortedata:    entry.Geboortedata,
				Pakketten:       []VoedingPakket{},
			}
			for _, pakket := range entry.Geschiedenis {
				voeding.Pakketten = append(voeding.Pakketten, VoedingPakket{
					ID:        pakket.ID,
					Datum:     pakket.Datum,
					Gekregen:  pakket.Gekregen,
					Opmerking: pakket.Opmerking,
				})
			}
			return voeding, nil
		}
	}

	return Voeding{}, ErrNotFound
}
