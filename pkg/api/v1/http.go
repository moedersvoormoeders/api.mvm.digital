package v1

import (
	"github.com/moedersvoormoeders/api.mvm.digital/pkg/db"

	"github.com/labstack/echo/v4"
)

type HTTPHandler struct {
	db *db.Connection
}

func NewHTTPHandler(db *db.Connection) *HTTPHandler {
	return &HTTPHandler{
		db: db,
	}
}

func (h *HTTPHandler) Register(e *echo.Echo) {

	// materiaal
	e.GET("/v1/materiaal/objects", h.getMateriaalObjects)
	e.GET("/v1/materiaal/klant/:mvmnummer", h.getMateriaalForKlant)
	e.POST("/v1/materiaal/klant/:mvmnummer", h.postMateriaalForKlant)
	e.POST("/v1/sinterklaas/klant/:mvmnummer", h.postSinterklaasForKlant)

	//whoami
	e.GET("/v1/auth/check", h.checkAuth)
	e.GET("/v1/whoami/roles", h.getRoles)
}
