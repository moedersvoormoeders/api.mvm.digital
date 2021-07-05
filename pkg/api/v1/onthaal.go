package v1

import (
	"log"
	"net/http"
	"sort"

	"github.com/labstack/echo/v4"
	"github.com/moedersvoormoeders/api.mvm.digital/pkg/db"
	"gorm.io/gorm"
)

func init() {
	registers = append(registers, func(e *echo.Echo, h *HTTPHandler) {
		e.GET("/v1/onthaal/voeding/:mvmnummer", h.getOnthaalVoedingForKlant)
	})
}

func (h *HTTPHandler) getOnthaalVoedingForKlant(c echo.Context) error {
	mvmNummer := c.Param("mvmnummer")
	if mvmNummer == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "mvmnummer not set"})
	}
	voeding := db.Voeding{}
	res := h.db.Where(&db.Voeding{MVMNummer: mvmNummer}).First(&voeding)

	if res.Error != nil && res.Error != gorm.ErrRecordNotFound {
		log.Println("DB error: ", res.Error)
		return c.JSON(http.StatusInternalServerError, res.Error.Error())
	}

	var gekregen []db.OntvangEntry
	h.db.Limit(10).Order("datum DESC").Where(&db.OntvangEntry{
		VoedingID: int(voeding.ID),
	}).Find(&gekregen)

	if res.Error != nil && res.Error != gorm.ErrRecordNotFound {
		log.Println("DB error: ", res.Error)
		return c.JSON(http.StatusInternalServerError, res.Error.Error())
	}

	err := h.db.FillGekregen(gekregen)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, res.Error.Error())
	}

	// sort by Go till i figure out how to do it in Gorm
	sort.Slice(gekregen, func(i, j int) bool {
		return gekregen[i].Datum.Unix() > gekregen[j].Datum.Unix()
	})

	voeding.Gekregen = gekregen

	return c.JSON(http.StatusOK, voeding)
}
