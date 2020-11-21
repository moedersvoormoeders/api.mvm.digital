package zoho

import (
	"time"

	zoho "github.com/schmorrison/Zoho"
)

// Sinterklaas defines a niceified and incomplete data struct of a sikterklaas inscrhijving
// This is to be moved from Zoho to our own DB in a later stage
type Sinterklaas struct {
	// TODO add ID
	ZohoID    string              `json:"zohoID""`
	Naam      string              `json:"naam"`
	MVMNummer string              `json:"mvmNummer"`
	Tijdslot  string              `json:"tijdslot"`
	Pakketten []SinterklaasPakket `json:"paketten"`
}

type SinterklaasPakket struct {
	ID        string  `json:"id"`
	Komt      bool    `json:"komt"`
	Opmerking string  `json:"opmerking"`
	Geslacht  string  `json:"geslacht"`
	Leeftijd  float64 `json:"leeftijd"`
	Naam      string  `json:"naam"`
}

type zohoSinterklaas struct {
	Data []struct {
		Owner struct {
			Name  string `json:"name"`
			ID    string `json:"id"`
			Email string `json:"email"`
		} `json:"Owner"`
		CurrencySymbol string `json:"$currency_symbol"`
		Kinderen       []struct {
			Geslacht string  `json:"Geslacht"`
			Leeftijd float64 `json:"Leeftijd"`
			InMerge  bool    `json:"$in_merge"`
			Komt     bool    `json:"Komt"`
			ParentID struct {
				Module string `json:"module"`
				Name   string `json:"name"`
				ID     string `json:"id"`
			} `json:"Parent_Id"`
			ID            string `json:"id"`
			Naam          string `json:"Naam"`
			Opmerking     string `json:"Opmerking"`
			Orchestration bool   `json:"$orchestration"`
		} `json:"Kinderen"`
		PhotoID       interface{} `json:"$photo_id"`
		Komt          bool        `json:"Komt"`
		ReviewProcess struct {
			Approve  bool `json:"approve"`
			Reject   bool `json:"reject"`
			Resubmit bool `json:"resubmit"`
		} `json:"$review_process"`
		UpcomingActivity interface{} `json:"$upcoming_activity"`
		Name             string      `json:"Name"`
		LastActivityTime time.Time   `json:"Last_Activity_Time"`
		Review           interface{} `json:"$review"`
		State            string      `json:"$state"`
		UnsubscribedMode interface{} `json:"Unsubscribed_Mode"`
		ProcessFlow      bool        `json:"$process_flow"`
		Klant            struct {
			Module string `json:"module"`
			Name   string `json:"name"`
			ID     string `json:"id"`
		} `json:"Klant"`
		Straat        string `json:"Straat"`
		Classificatie string `json:"Classificatie"`
		ID            string `json:"id"`
		Naam          string `json:"Naam"`
		Approved      bool   `json:"$approved"`
		Approval      struct {
			Delegate bool `json:"delegate"`
			Approve  bool `json:"approve"`
			Reject   bool `json:"reject"`
			Resubmit bool `json:"resubmit"`
		} `json:"$approval"`
		ModifiedTime     time.Time     `json:"Modified_Time"`
		Huisnummer       string        `json:"Huisnummer"`
		GSMNummer        string        `json:"GSM_Nummer"`
		Tijdslot         string        `json:"Tijdslot"`
		CreatedTime      time.Time     `json:"Created_Time"`
		UnsubscribedTime interface{}   `json:"Unsubscribed_Time"`
		Editable         bool          `json:"$editable"`
		Opmerking        interface{}   `json:"Opmerking"`
		Code             string        `json:"Code"`
		Postcode         string        `json:"Postcode"`
		Orchestration    bool          `json:"$orchestration"`
		DoelgroepNummer  string        `json:"Doelgroep_Nummer"`
		InMerge          bool          `json:"$in_merge"`
		Gemeente         string        `json:"Gemeente"`
		Status           string        `json:"$status"`
		Tag              []interface{} `json:"Tag"`
		ApprovalState    string        `json:"$approval_state"`
	} `json:"data"`
}

func (c *CRM) GetSinterklaasForMVMNummer(mvmNummer string) (Sinterklaas, error) {
	out, err := c.zohoCRM.SearchRecords(&zohoSinterklaas{}, "Inschrijvingen", map[string]zoho.Parameter{
		"word":     zoho.Parameter(mvmNummer),
		"per_page": "200",
	})

	if err != nil {
		return Sinterklaas{}, err
	}

	data := out.(*zohoSinterklaas)
	dataID := ""
	for _, entry := range data.Data {
		if entry.DoelgroepNummer == mvmNummer {
			dataID = entry.ID
		}
	}

	if dataID == "" {
		return Sinterklaas{}, ErrNotFound
	}

	// now the same again but with the ID so we get all data
	out, err = c.zohoCRM.GetRecord(&zohoSinterklaas{}, "Inschrijvingen", dataID)

	if err != nil {
		return Sinterklaas{}, err
	}

	data = out.(*zohoSinterklaas)
	for _, entry := range data.Data {
		if entry.DoelgroepNummer == mvmNummer {
			sinterklaas := Sinterklaas{
				ZohoID: entry.ID,

				MVMNummer: entry.DoelgroepNummer,
				Tijdslot:  entry.Tijdslot,
				Naam:      entry.Naam,
				Pakketten: []SinterklaasPakket{},
			}
			for _, pakket := range entry.Kinderen {
				sinterklaas.Pakketten = append(sinterklaas.Pakketten, SinterklaasPakket{
					ID:        pakket.ID,
					Komt:      pakket.Komt,
					Opmerking: pakket.Opmerking,
					Geslacht:  pakket.Geslacht,
					Leeftijd:  pakket.Leeftijd,
					Naam:      pakket.Naam,
				})
			}
			return sinterklaas, nil
		}
	}

	return Sinterklaas{}, ErrNotFound
}
