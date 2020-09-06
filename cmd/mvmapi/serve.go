package main

import (
	"context"
	"errors"
	"fmt"

	v1 "github.com/moedersvoormoeders/api.mvm.digital/pkg/api/v1"

	"github.com/dgrijalva/jwt-go"

	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/bcrypt"

	"github.com/moedersvoormoeders/api.mvm.digital/pkg/api/auth"
	zohohttp "github.com/moedersvoormoeders/api.mvm.digital/pkg/api/zoho"
	"github.com/moedersvoormoeders/api.mvm.digital/pkg/db"
	"github.com/moedersvoormoeders/api.mvm.digital/pkg/zoho"
)

// this is used to compare to in case of no user found to keep the response time the same
const dummyHash = `$2a$10$8KqKzq6uHCL72Qhshj9L.uGUz/0lmkjupqYQKCy6th9Rv91k53g82`

func init() {
	rootCmd.AddCommand(NewServeCmd())
}

var protectedEntryPoints = []string{"/zoho", "/v1"}

type serveCmdOptions struct {
	BindAddr string
	Port     int

	jwtSecret string

	db *db.Connection

	zohoClientID     string
	zohoClientSecret string
	zohoCRM          *zoho.CRM

	postgresHost     string
	postgresPort     int
	postgresUsername string
	postgresDatabase string
	postgresPassword string
}

// NewServeCmd generates the `serve` command
func NewServeCmd() *cobra.Command {
	s := serveCmdOptions{}
	c := &cobra.Command{
		Use:     "serve",
		Short:   "Serves the HTTP REST endpoint",
		Long:    `Serves the HTTP REST endpoint on the given bind address and port`,
		PreRunE: s.Validate,
		RunE:    s.RunE,
	}
	c.Flags().StringVarP(&s.BindAddr, "bind-address", "b", "0.0.0.0", "address to bind port to")
	c.Flags().IntVarP(&s.Port, "port", "p", 8080, "Port to listen on")

	c.Flags().StringVar(&s.jwtSecret, "jwt-secret", "", "JWT siging key, please do not set in flags")

	// needed for Zoho connector
	c.Flags().StringVar(&s.zohoClientID, "zoho-clientid", "", "Zoho client ID, you can get this at https://accounts.zoho.eu/developerconsole")
	c.Flags().StringVar(&s.zohoClientSecret, "zoho-clientsecret", "", "Zoho client Secret, you can get this at https://accounts.zoho.eu/developerconsole")

	c.Flags().StringVar(&s.postgresHost, "postgres-host", "", "PostgreSQL hostname")
	c.Flags().IntVar(&s.postgresPort, "postgres-port", 5432, "PostgreSQL hostname")
	c.Flags().StringVar(&s.postgresUsername, "postgres-username", "", "PostgreSQL hostname")
	c.Flags().StringVar(&s.postgresPassword, "postgres-password", "", "PostgreSQL hostname")
	c.Flags().StringVar(&s.postgresDatabase, "postgres-database", "", "PostgreSQL hostname")
	return c
}

func (s *serveCmdOptions) Validate(cmd *cobra.Command, args []string) error {
	if s.jwtSecret == "" {
		return errors.New("jwt-secret not set")
	}

	if s.postgresUsername == "" || s.postgresPassword == "" || s.postgresDatabase == "" || s.postgresHost == "" {
		return errors.New("PostgreSQL credentials not set")
	}

	return nil
}

func (s *serveCmdOptions) RunE(cmd *cobra.Command, args []string) error {
	ctx, cancel := context.WithCancel(context.Background())
	printLogo()

	if s.zohoClientID == "" {
		fmt.Println("Zoho client ID not set, not loading Zoho integration")
	} else {
		s.zohoCRM = zoho.NewCRM()
		err := s.zohoCRM.Connect(s.zohoClientID, s.zohoClientSecret)
		if err != nil {
			return fmt.Errorf("error connecting to Zoho: %w", err)
		}

		defer s.zohoCRM.Close()
	}

	var err error
	s.db, err = db.NewConnection(db.ConnectionDetails{
		Host:     s.postgresHost,
		Port:     s.postgresPort,
		User:     s.postgresUsername,
		Database: s.postgresDatabase,
		Password: s.postgresPassword,
	})
	if err != nil {
		return fmt.Errorf("error opening database: %w", err)
	}

	err = s.db.DoMigrate()
	if err != nil {
		return fmt.Errorf("error migrating database: %w", err)
	}

	e := echo.New()
	e.HideBanner = true
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey: []byte(s.jwtSecret),
		Claims:     &auth.Claim{},
		Skipper: func(c echo.Context) bool {
			// always skip JWT unless path is a protectedPrefix
			for _, protectedPrefix := range protectedEntryPoints {
				if strings.HasPrefix(c.Path(), protectedPrefix) {
					return false
				}
			}
			return true
		},
	}))

	// handlers
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "mvm.digital API endpoint")
	})
	e.POST("/login", s.login)
	if s.zohoClientID != "" {
		zohohttp.NewHTTPHandler().Register(e, s.zohoCRM)
	}
	v1.NewHTTPHandler(s.db).Register(e)

	go func() {
		e.Start(fmt.Sprintf("%s:%d", s.BindAddr, s.Port))
		cancel() // server ended, stop the world
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	for {
		select {
		case <-c:
			cancel()
		case <-ctx.Done():
			return nil
		}
	}
}

type AuthData struct {
	Username string `form:"username" json:"username"`
	Password string `form:"password" json:"password"`
}

// TODO: move this!
func (s *serveCmdOptions) login(c echo.Context) error {
	data := new(AuthData)
	err := c.Bind(data)
	if err != nil {
		log.Println(err)
	}

	if data.Username == "" || data.Password == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "username or password not specified"})
	}

	data.Username = strings.ToLower(data.Username)

	// TODO: check login
	user := db.User{}
	err = s.db.GetWhereIs(&user, "username", data.Username)
	if errors.Is(err, db.ErrorNotFound) {
		_ = bcrypt.CompareHashAndPassword([]byte(dummyHash), []byte(data.Password))
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "username or password incorrect"})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data.Password)); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "username or password incorrect"})
	}

	// Set custom claims
	claims := &auth.Claim{
		user.Name,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"token": t,
	})
}
