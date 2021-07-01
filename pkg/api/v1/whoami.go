package v1

import (
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/moedersvoormoeders/api.mvm.digital/pkg/api/auth"

	"github.com/labstack/echo/v4"
	"github.com/moedersvoormoeders/api.mvm.digital/pkg/db"
)

func init() {
	registers = append(registers, func(e *echo.Echo, h *HTTPHandler) {
		e.GET("/v1/auth/check", h.checkAuth)
		e.GET("/v1/auth/whoami/roles", h.getRoles)
	})
}

func (h *HTTPHandler) getRoles(c echo.Context) error {
	roles := []db.Role{}

	user, ok := c.Get("user").(*jwt.Token)
	if !ok {
		return c.JSON(http.StatusUnauthorized, echo.Map{"status": "JWT incorrect"})
	}
	claims, ok := user.Claims.(*auth.Claim)
	if !ok || claims.Name == "" {
		return c.JSON(http.StatusUnauthorized, echo.Map{"status": "JWT incorrect"})
	}

	roleBindings := []db.RoleBinding{}
	res := h.db.Preload("Role.Verbs").Preload("Role.Endpoints").Where("user_id = ?", claims.ID).Find(&roleBindings)
	if res.Error != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": res.Error.Error()})
	}

	for _, rb := range roleBindings {
		roles = append(roles, rb.Role)
	}

	return c.JSON(http.StatusOK, roles)
}

func (h *HTTPHandler) checkAuth(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*auth.Claim)
	if claims.Name == "" {
		return c.JSON(http.StatusUnauthorized, echo.Map{"status": "JWT incorrect"})
	}
	return c.JSON(http.StatusOK, echo.Map{"status": "ok"})
}
