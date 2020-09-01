package zoho

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/moedersvoormoeders/api.mvm.digital/pkg/zoho"
)

type HTTPHandler struct {
	zohoCRM *zoho.CRM
}

func NewHTTPHandler() *HTTPHandler {
	return &HTTPHandler{}
}

func (h *HTTPHandler) Register(e *echo.Echo, z *zoho.CRM) {
	h.zohoCRM = z

	e.GET("/zoho/klant/", h.getKlantForMVMNummer)
	e.GET("/zoho/voeding/", h.getVoedingForMVMNummer)
}

func (h *HTTPHandler) getKlantForMVMNummer(c echo.Context) error {
	mvmNummer := c.QueryParam("mvmNummer")

	if mvmNummer == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "missing mvmNummer query parameter"})
	}

	klant, err := h.zohoCRM.GetKlantForMVMNummer(mvmNummer)
	if err == zoho.ErrNotFound {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "geen klant gevonden"})
	}
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, klant)
}

func (h *HTTPHandler) getVoedingForMVMNummer(c echo.Context) error {
	mvmNummer := c.QueryParam("mvmNummer")

	if mvmNummer == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "missing mvmNummer query parameter"})
	}

	voeding, err := h.zohoCRM.GetVoedingForMVMNummer(mvmNummer)
	if err == zoho.ErrNotFound {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "geen klant gevonden"})
	}
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err})
	}

	return c.JSON(http.StatusOK, voeding)
}
