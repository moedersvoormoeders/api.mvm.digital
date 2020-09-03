package v1

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/moedersvoormoeders/api.mvm.digital/pkg/db"
)

func (h *HTTPHandler) getMateriaalObjects(c echo.Context) error {
	materiaalObjects := []db.MateriaalObject{}
	res := h.db.Preload("Categorie").Find(&materiaalObjects)

	if res.Error != nil {
		// TODO: look into how JS handles this
		return res.Error
	}

	return c.JSON(http.StatusOK, materiaalObjects)
}

func (h *HTTPHandler) getMateriaalForKlant(c echo.Context) error {
	mvmNummer := c.Param("mvmnummer")
	if mvmNummer == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "mvmnummer not set"})
	}
	materiaal := db.Materiaal{}
	err := h.db.GetWhereIs(&materiaal, "mvm_nummer", mvmNummer)

	if err != nil && err != db.ErrorNotFound {
		log.Println("DB error: ", err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, materiaal)
}
