package v1

import (
	"log"
	"net/http"

	"gorm.io/gorm/clause"

	"github.com/labstack/echo/v4"
	"github.com/moedersvoormoeders/api.mvm.digital/pkg/db"
	"gorm.io/gorm"
)

func (h *HTTPHandler) getMateriaalObjects(c echo.Context) error {
	materiaalObjects := []db.MateriaalObject{}
	res := h.db.Preload(clause.Associations).Find(&materiaalObjects)

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
	res := h.db.Preload(clause.Associations).Where(&db.Materiaal{MVMNummer: mvmNummer}).First(&materiaal)

	if res.Error != nil && res.Error != gorm.ErrRecordNotFound {
		log.Println("DB error: ", res.Error)
		return c.JSON(http.StatusInternalServerError, res.Error.Error())
	}

	err := h.db.FillMateriaalGekregen(&materiaal)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, res.Error.Error())
	}

	return c.JSON(http.StatusOK, materiaal)
}

func (h *HTTPHandler) postMateriaalForKlant(c echo.Context) error {
	mvmNummer := c.Param("mvmnummer")
	if mvmNummer == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "mvmnummer not set"})
	}
	materiaal := db.Materiaal{}
	err := c.Bind(&materiaal)
	if err != nil {
		log.Println("Body parse error", err)
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	materiaal.MVMNummer = mvmNummer

	materiaalFromDB := db.Materiaal{}
	err = h.db.GetWhereIs(&materiaalFromDB, "mvm_nummer", mvmNummer)
	if err != nil && err != db.ErrorNotFound {
		log.Println("DB error: ", err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	} else if err == db.ErrorNotFound {
		h.db.Create(&materiaal)
		return c.JSON(http.StatusOK, echo.Map{"status": "ok", "message": "Materiaal is opgeslagen"})
	}

	materiaal.Model = materiaalFromDB.Model
	res := h.db.Updates(&materiaal)
	if res.Error != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	err = h.db.Model(&materiaal).Association("Gekregen").Replace(&materiaal.Gekregen)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	for _, mat := range materiaal.Gekregen {
		res := h.db.Updates(&mat)
		if res.Error != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
	}

	return c.JSON(http.StatusOK, echo.Map{"status": "ok", "message": "Materiaal is opgeslagen"})
}
