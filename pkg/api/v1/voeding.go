package v1

import (
	"fmt"
	"log"
	"net/http"
	"sort"
	"strconv"

	"gorm.io/gorm/clause"

	"github.com/labstack/echo/v4"
	"github.com/moedersvoormoeders/api.mvm.digital/pkg/db"
	"gorm.io/gorm"
)

func init() {
	registers = append(registers, func(e *echo.Echo, h *HTTPHandler) {
		e.GET("/v1/voeding/objects", h.getVoedingObjects)
		e.GET("/v1/voeding/klant/:mvmnummer", h.getVoedingForKlant)
		e.POST("/v1/voeding/klant/:mvmnummer", h.postVoedingForKlant)
		e.GET("/v1/voeding/klant/:mvmnummer/gekregen", h.getVoedingGekregen)
		e.POST("/v1/voeding/klant/:mvmnummer/gekregen/:id", h.postVoedingRij)
		e.DELETE("/v1/voeding/klant/:mvmnummer/gekregen/:id", h.deleteVoedingRij)
	})
}

func (h *HTTPHandler) getVoedingObjects(c echo.Context) error {
	cagegories := []db.Categorie{}
	res := h.db.Preload(clause.Associations).Where("afdeling", "Voeding").Order("naam").Find(&cagegories)
	if res.Error != nil {
		return res.Error
	}

	var objects []db.Object
	for _, cat := range cagegories {
		var catObj []db.Object
		res := h.db.Preload(clause.Associations).Where("categorie_id", cat.ID).Order("naam").Find(&catObj)
		if res.Error != nil {
			return res.Error
		}

		objects = append(objects, catObj...)
	}

	return c.JSON(http.StatusOK, objects)
}

func (h *HTTPHandler) getVoedingGekregen(c echo.Context) error {
	mvmNummer := c.Param("mvmnummer")
	if mvmNummer == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "mvmnummer not set"})
	}
	voeding := db.Voeding{}
	res := h.db.Where(&db.Voeding{MVMNummer: mvmNummer}).First(&voeding)

	if voeding.ID == 0 {
		c.Response().Header().Add("Num-Total-Entries", fmt.Sprintf("%d", 0))

		return c.JSON(http.StatusOK, []db.Voeding{})
	}

	var gekregen []db.OntvangEntry
	h.db.Order("datum DESC").Scopes(Paginate(c)).Where(&db.OntvangEntry{
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

	var count int64
	h.db.Model(&db.OntvangEntry{}).Where("voeding_id = ?", int(voeding.ID)).Count(&count)

	c.Response().Header().Add("Num-Total-Entries", fmt.Sprintf("%d", count))

	return c.JSON(http.StatusOK, gekregen)
}

func (h *HTTPHandler) getVoedingForKlant(c echo.Context) error {
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

	return c.JSON(http.StatusOK, voeding)
}

func (h *HTTPHandler) postVoedingForKlant(c echo.Context) error {
	mvmNummer := c.Param("mvmnummer")
	if mvmNummer == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "mvmnummer not set"})
	}

	voeding := db.Voeding{}
	err := c.Bind(&voeding)
	if err != nil {
		log.Println("Body parse error", err)
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	voeding.MVMNummer = mvmNummer
	voeding.Gekregen = nil // we don't want to update that in this call we should use a seperate API for this

	if voeding.ID < 1 {
		res := h.db.Create(&voeding)
		if res.Error != nil {
			return err
		}
		return c.JSON(http.StatusOK, echo.Map{"status": "ok", "message": "Voeding is opgeslagen", "data": voeding})
	}

	res := h.db.Select("SpecialeVoeding", "Opmerking").Updates(&voeding)
	if res.Error != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"status": "error", "message": res.Error.Error()})
	}

	return c.JSON(http.StatusOK, echo.Map{"status": "ok", "message": "Voeding is opgeslagen", "data": voeding})
}

func (h *HTTPHandler) postVoedingRij(c echo.Context) error {
	mvmNummer := c.Param("mvmnummer")
	if mvmNummer == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"status": "error", "message": "mvmnummer not set"})
	}

	idStr := c.Param("id")
	if idStr == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"status": "error", "message": "id not set"})
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"status": "error", "message": "invalid id"})
	}

	voeding := db.Voeding{}
	res := h.db.Preload(clause.Associations).Where(&db.Voeding{MVMNummer: mvmNummer}).First(&voeding)
	if res.Error != nil {
		return res.Error
	}

	entry := db.OntvangEntry{}
	err = c.Bind(&entry)
	if err != nil {
		log.Println("Body parse error", err)
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	entry.VoedingID = int(voeding.ID)
	entry.ID = uint(id)

	if id < 1 {
		res := h.db.Create(&entry)
		if res.Error != nil {
			return err
		}
		return c.JSON(http.StatusOK, echo.Map{"status": "ok", "message": "Voeding entry is opgeslagen", "data": entry})
	}

	res = h.db.Updates(&entry)
	if res.Error != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"status": "error", "message": res.Error.Error()})
	}

	return c.JSON(http.StatusOK, echo.Map{"status": "ok", "message": "Voeding entry is opgeslagen", "data": entry})
}

func (h *HTTPHandler) deleteVoedingRij(c echo.Context) error {
	mvmNummer := c.Param("mvmnummer")
	if mvmNummer == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"status": "error", "message": "mvmnummer not set"})
	}

	idStr := c.Param("id")
	if idStr == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"status": "error", "message": "id not set"})
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"status": "error", "message": "invalid id"})
	}

	voeding := db.Voeding{}
	res := h.db.Preload(clause.Associations).Where(&db.Voeding{MVMNummer: mvmNummer}).First(&voeding)
	if res.Error != nil {
		return res.Error
	}

	res = h.db.Delete(&db.OntvangEntry{
		VoedingID: int(voeding.ID),
		Model:     gorm.Model{ID: uint(id)},
	})
	if res.Error != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"status": "error", "message": res.Error.Error()})
	}

	return c.JSON(http.StatusOK, echo.Map{"status": "ok", "message": "Voeding entry is verwijderd"})
}
