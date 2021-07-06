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

	e.GET("/zoho/search/", h.getKlantSearch)
	e.GET("/zoho/klant/", h.getKlantForMVMNummer)
	e.GET("/zoho/contacten/", h.getContactenForMVMNummer)
	e.GET("/zoho/sinterklaas/", h.getSinterklaasForMVMNummer)
}

func (h *HTTPHandler) getKlantSearch(c echo.Context) error {
	query := c.QueryParam("query")

	if query == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "missing search query query parameter"})
	}

	klant, err := h.zohoCRM.GetKlantenForQuery(query)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, klant)
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

func (h *HTTPHandler) getSinterklaasForMVMNummer(c echo.Context) error {
	mvmNummer := c.QueryParam("mvmNummer")

	if mvmNummer == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "missing mvmNummer query parameter"})
	}

	sint, err := h.zohoCRM.GetSinterklaasForMVMNummer(mvmNummer)
	if err == zoho.ErrNotFound {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "geen sinterklaas gevonden"})
	}
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, sint)
}

func (h *HTTPHandler) getContactenForMVMNummer(c echo.Context) error {
	mvmNummer := c.QueryParam("mvmNummer")

	if mvmNummer == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "missing mvmNummer query parameter"})
	}

	contacten, err := h.zohoCRM.GetContactenForMVMNummer(mvmNummer)
	if err == zoho.ErrNotFound {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "geen klant gevonden"})
	}
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err})
	}

	return c.JSON(http.StatusOK, contacten)
}
