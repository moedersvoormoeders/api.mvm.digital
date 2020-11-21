package main

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/moedersvoormoeders/api.mvm.digital/pkg/zoho"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(NewCheckSintCmd())
}

type checkSintCmdOptions struct {
	zohoClientID     string
	zohoClientSecret string
	zohoCRM          *zoho.CRM
}

// NewCheckSintCmd generates the `check-sint` command
func NewCheckSintCmd() *cobra.Command {
	a := checkSintCmdOptions{}
	c := &cobra.Command{
		Use:     "check-sint",
		Short:   "Check sinterklaas inschrinvingen",
		PreRunE: a.Validate,
		RunE:    a.RunE,
	}

	c.Flags().StringVar(&a.zohoClientID, "zoho-clientid", "", "Zoho client ID, you can get this at https://accounts.zoho.eu/developerconsole")
	c.Flags().StringVar(&a.zohoClientSecret, "zoho-clientsecret", "", "Zoho client Secret, you can get this at https://accounts.zoho.eu/developerconsole")
	viper.BindPFlags(c.Flags())

	return c
}

func (c *checkSintCmdOptions) Validate(cmd *cobra.Command, args []string) error {
	if c.zohoClientID == "" || c.zohoClientSecret == "" {
		return errors.New("need to have Zoho tokens set" +
			"")
	}
	return nil
}

func (c *checkSintCmdOptions) RunE(cmd *cobra.Command, args []string) error {

	c.zohoCRM = zoho.NewCRM()
	err := c.zohoCRM.Connect(c.zohoClientID, c.zohoClientSecret)
	if err != nil {
		return fmt.Errorf("error connecting to Zoho: %w", err)
	}

	defer c.zohoCRM.Close()

	klanten, err := c.zohoCRM.GetAllKlanten()
	if err != nil {
		return err
	}
	log.Printf("Fetched all clients")

	for _, klant := range klanten {
		c.GetSint(klant.MVMNummer, 0)
	}

	return nil
}

func (c *checkSintCmdOptions) GetSint(mvmNummer string, attempt int) {
	time.Sleep(time.Millisecond * 200)
	log.Printf("Getting %s\n", mvmNummer)
	sint, err := c.zohoCRM.GetSinterklaasForMVMNummer(mvmNummer)
	if err != nil && attempt > 3 {
		log.Println(err, mvmNummer)
		return
	} else if err != nil {
		attempt++
		c.GetSint(mvmNummer, attempt)
		return
	}

	comes := false
	hasChildWithNoCheckmak := false
	for _, pakket := range sint.Pakketten {
		if pakket.Komt {
			comes = true
		} else {
			hasChildWithNoCheckmak = true
		}
	}

	if comes && hasChildWithNoCheckmak {
		log.Printf("Incorrect registration for %s\n", sint.MVMNummer)
	}
}
