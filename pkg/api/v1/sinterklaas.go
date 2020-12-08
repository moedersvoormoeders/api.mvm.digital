package v1

import (
	"log"
	"net/http"
	"time"

	"gorm.io/gorm/clause"

	"github.com/labstack/echo/v4"
	"github.com/moedersvoormoeders/api.mvm.digital/pkg/db"
)

type sinterklaasRequest struct {
	Kinderen []string `json:"kinderen"`
}

func (h *HTTPHandler) postSinterklaasForKlant(c echo.Context) error {
	mvmNummer := c.Param("mvmnummer")
	if mvmNummer == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "mvmnummer not set"})
	}
	data := sinterklaasRequest{}
	err := c.Bind(&data)
	if err != nil {
		log.Println("Body parse error", err)
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	materiaalFromDB := db.Materiaal{}
	res := h.db.Preload(clause.Associations).Where(&db.Materiaal{MVMNummer: mvmNummer}).First(&materiaalFromDB)
	if res.Error != nil {
		log.Println("DB error: ", err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	} else if res.RowsAffected == 0 {
		h.db.Create(&db.Materiaal{
			MVMNummer: mvmNummer,
		})
		err = h.db.GetWhereIs(&materiaalFromDB, "mvm_nummer", mvmNummer)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
	}

	sintPakket := db.MateriaalObject{}
	err = h.db.GetWhereIs(&sintPakket, "naam", "Sinterklaas")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	for _, kind := range data.Kinderen {
		materiaalFromDB.Gekregen = append(materiaalFromDB.Gekregen, db.MateriaalEntry{
			Datum:       time.Now(),
			Aantal:      1,
			ObjectID:    int(sintPakket.ID),
			Object:      sintPakket,
			Maat:        db.MateriaalMaat{},
			Opmerking:   "Automatische registratie Sinterklaas",
			Ontvanger:   kind,
			MateriaalID: materiaalFromDB.ID,
		})
	}

	res = h.db.Updates(&materiaalFromDB)
	if res.Error != nil {
		return c.JSON(http.StatusInternalServerError, res.Error.Error())
	}

	err = h.db.Model(&materiaalFromDB).Association("Gekregen").Replace(&materiaalFromDB.Gekregen)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	for _, mat := range materiaalFromDB.Gekregen {
		res := h.db.Updates(&mat)
		if res.Error != nil {
			return c.JSON(http.StatusInternalServerError, res.Error.Error())
		}
	}

	return c.JSON(http.StatusOK, echo.Map{"status": "ok", "message": "Sinterklaas is opgeslagen"})
}
