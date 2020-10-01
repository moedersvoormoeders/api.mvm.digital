package zoho

import (
	"fmt"
	"strings"
	"time"

	"github.com/davecgh/go-spew/spew"

	zoho "github.com/schmorrison/Zoho"
	"github.com/schmorrison/Zoho/crm"
)

// Contact defines a niceified data struct of a contact
// this can be a child or a partner
// This is a structure that should be reformed in a family member
// This is to be moved from Zoho to our own DB in a later stage
type Contact struct {
	// TODO add ID
	ZohoID        string    `json:"zohoID"`
	Naam          string    `json:"naam"`
	Voornaam      string    `json:"voornaam"`
	Geslacht      string    `json:"geslacht"`
	GeboorteDatum time.Time `json:"geboorteDatum"`
	SoortRelatie  string    `json:"soortRelatie"`
}

// this is contact as is defined in the Zoho data
type zohoContact struct {
	Data []struct {
		Owner struct {
			Name  string `json:"name"`
			ID    string `json:"id"`
			Email string `json:"email"`
		} `json:"Owner"`
		Description    interface{} `json:"Description"`
		CurrencySymbol string      `json:"$currency_symbol"`
		PhotoID        interface{} `json:"$photo_id"`
		ReviewProcess  struct {
			Approve  bool `json:"approve"`
			Reject   bool `json:"reject"`
			Resubmit bool `json:"resubmit"`
		} `json:"$review_process"`
		Salutation       interface{} `json:"Salutation"`
		LastActivityTime interface{} `json:"Last_Activity_Time"`
		FirstName        string      `json:"First_Name"`
		FullName         string      `json:"Full_Name"`
		RecordImage      interface{} `json:"Record_Image"`
		ModifiedBy       struct {
			Name  string `json:"name"`
			ID    string `json:"id"`
			Email string `json:"email"`
		} `json:"Modified_By"`
		Review           interface{} `json:"$review"`
		State            string      `json:"$state"`
		UnsubscribedMode interface{} `json:"Unsubscribed_Mode"`
		ProcessFlow      bool        `json:"$process_flow"`
		AccountName      struct {
			Name string `json:"name"`
			ID   string `json:"id"`
		} `json:"Account_Name"`
		ID          string      `json:"id"`
		Approved    bool        `json:"$approved"`
		ReportingTo interface{} `json:"Reporting_To"`
		Geslacht    string      `json:"Geslacht"`
		Approval    struct {
			Delegate bool `json:"delegate"`
			Approve  bool `json:"approve"`
			Reject   bool `json:"reject"`
			Resubmit bool `json:"resubmit"`
		} `json:"$approval"`
		ModifiedTime     time.Time   `json:"Modified_Time"`
		DateOfBirth      string      `json:"Date_of_Birth"`
		CreatedTime      time.Time   `json:"Created_Time"`
		UnsubscribedTime interface{} `json:"Unsubscribed_Time"`
		Editable         bool        `json:"$editable"`
		Orchestration    bool        `json:"$orchestration"`
		SoortRelatie     string      `json:"Soort_relatie"`
		LastName         string      `json:"Last_Name"`
		InMerge          bool        `json:"$in_merge"`
		Status           string      `json:"$status"`
		CreatedBy        struct {
			Name  string `json:"name"`
			ID    string `json:"id"`
			Email string `json:"email"`
		} `json:"Created_By"`
		Tag           []interface{} `json:"Tag"`
		ApprovalState string        `json:"$approval_state"`
	} `json:"data"`
	Info struct {
		PerPage     int  `json:"per_page"`
		Count       int  `json:"count"`
		Page        int  `json:"page"`
		MoreRecords bool `json:"more_records"`
	} `json:"info"`
}

func (c *CRM) GetContactenForMVMNummer(mvmNummer string) ([]Contact, error) {
	klant, err := c.GetKlantForMVMNummer(mvmNummer)
	if err != nil {
		return nil, err
	}
	out, err := c.zohoCRM.ListRecords(&zohoContact{}, crm.CRMModule(fmt.Sprintf("Accounts/%s/Contacts", klant.ZohoID)), map[string]zoho.Parameter{
		"per_page": "200",
	})

	if err != nil {
		// no contacts error
		if strings.Contains(err.Error(), "There is no content available for the request") {
			return []Contact(nil), nil
		}
		spew.Dump(err)
		return nil, err
	}

	data := out.(*zohoContact)

	contacten := []Contact(nil)
	for _, entry := range data.Data {
		// ignoring errors as i fear we will get nils...
		date, _ := time.Parse("2006-01-02", entry.DateOfBirth)
		contacten = append(contacten, Contact{
			ZohoID:        entry.ID,
			Naam:          entry.LastName,
			Voornaam:      entry.FirstName,
			Geslacht:      entry.Geslacht,
			GeboorteDatum: date,
			SoortRelatie:  entry.SoortRelatie,
		})
	}

	return contacten, nil
}
