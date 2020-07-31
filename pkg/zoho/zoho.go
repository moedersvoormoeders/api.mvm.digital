package zoho

import (
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	zoho "github.com/schmorrison/Zoho"
	"github.com/schmorrison/Zoho/crm"
)

var ErrNotFound = errors.New("Not found")

type CRM struct {
	zohoCRM *crm.API

	tokenMutex sync.Mutex
}

func NewCRM() *CRM {
	return &CRM{}
}

func (c *CRM) Connect(clientID, clientSecret string) error {
	z := zoho.New()
	z.SetZohoTLD("eu")

	scopes := []zoho.ScopeString{
		zoho.BuildScope(zoho.Crm, zoho.ModulesScope, zoho.AllMethod, zoho.NoOp),
	}

	storedToken, err := z.LoadAccessAndRefreshToken()
	if err == nil && storedToken.AccessToken != "" {
		// Found stored token!
		log.Println("Using stored Zoho token")
	} else if err != nil {
		log.Println("Error with stored Zoho token, requesting new one", err)
		if err := z.AuthorizationCodeRequest(clientID, clientSecret, scopes, "http://localhost:8080/oauthredirect"); err != nil {
			return err
		}
	} else {
		log.Println("Found no stored Zoho token, requesting new one")
		if err := z.AuthorizationCodeRequest(clientID, clientSecret, scopes, "http://localhost:8080/oauthredirect"); err != nil {
			return err
		}
	}

	c.zohoCRM = crm.New(z)
	err = z.RefreshTokenRequest()
	if err != nil {
		return fmt.Errorf("error while refreshing token: %w", err)
	}

	log.Println("Renewed Zoho token")

	go func() {
		for {
			time.Sleep(time.Minute)
			c.tokenMutex.Lock()
			err := z.RefreshTokenRequest()
			c.tokenMutex.Unlock()
			if err != nil {
				log.Println(err)
			} else {
				log.Println("Renewed Zoho token")
			}
		}
	}()

	return nil
}

func (c *CRM) Close() {
	log.Println("Saving new Zoho token to disk...")
	c.tokenMutex.Lock()
	c.zohoCRM.RefreshTokenRequest()
	c.tokenMutex.Unlock()
}
