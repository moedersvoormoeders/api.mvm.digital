package v1

import (
	"net/http"

	"github.com/moedersvoormoeders/api.mvm.digital/pkg/db"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"

	"github.com/moedersvoormoeders/api.mvm.digital/pkg/api/auth"
)

type HTTPHandler struct {
	db db.Connection
}

func NewHTTPHandler() *HTTPHandler {
	return &HTTPHandler{}
}

func (h *HTTPHandler) Register(e *echo.Echo) {
	e.GET("/v1/auth/check", h.checkAuth)

	// materiaal
	e.GET("/v1/materiaal/objects", h.getMateriaalObjects)
	e.GET("/v1/materiaal/klant/:mvmnummer", h.getMateriaalForKlant)
}

func (h *HTTPHandler) checkAuth(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*auth.Claim)
	if claims.Name == "" {
		return c.JSON(http.StatusUnauthorized, echo.Map{"status": "JWT incorrect"})
	}
	return c.JSON(http.StatusOK, echo.Map{"status": "ok"})
}
