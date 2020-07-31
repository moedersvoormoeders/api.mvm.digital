package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(NewServeCmd())
}

var protectedEntryPoints = []string{"/zoho", "/v1"}

type serveCmdOptions struct {
	BindAddr string
	Port     int
}

// NewServeCmd generates the `serve` command
func NewServeCmd() *cobra.Command {
	s := serveCmdOptions{}
	c := &cobra.Command{
		Use:   "serve",
		Short: "Serves the HTTP REST endpoint",
		Long:  `Serves the HTTP REST endpoint on the given bind address and port`,
		RunE:  s.RunE,
	}
	c.Flags().StringVarP(&s.BindAddr, "bind-address", "b", "0.0.0.0", "address to bind port to")
	c.Flags().IntVarP(&s.Port, "port", "p", 8080, "Port to listen on")

	viper.BindPFlags(c.Flags())

	return c
}

func (s *serveCmdOptions) RunE(cmd *cobra.Command, args []string) error {

	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	e.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey: []byte("secret"),
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

	return e.Start(fmt.Sprintf("%s:%d", s.BindAddr, s.Port))
}
