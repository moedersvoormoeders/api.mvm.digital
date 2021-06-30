package db

import (
	"time"

	"gorm.io/gorm/clause"

	"gorm.io/gorm"
)

type Afdeling string

const (
	AfdelingMateriaal Afdeling = "materiaal"
	AfdelingVoeding   Afdeling = "voeding"
)

type PrintOpties struct {
	gorm.Model
	PrintKindInfo bool `json:"printKindInfo"`
}

type Categorie struct {
	gorm.Model
	Naam          string      `json:"naam"`
	PerPersoon    bool        `json:"perPersoon"`
	OpMaat        bool        `json:"opMaat"`
	PrintOpties   PrintOpties `json:"printOpties"`
	PrintOptiesID int         `json:"printOptiesID"`
	Order         int         `json:"order"`
	Afdeling      Afdeling    `json:"afdeling"`
}

type Maat struct {
	gorm.Model
	Naam     string `json:"naam"`
	ObjectID uint
}

type Object struct {
	gorm.Model
	Naam        string    `json:"naam"`
	CategorieID int       `json:"categorieID"`
	Categorie   Categorie `json:"categorie"`
	Hidden      bool      `json:"hidden"`
	Prijs       float64   `json:"prijs"` // most of the times this is 0.0
	Maten       []Maat    `json:"maten"`
}

type OntvangEntry struct {
	gorm.Model
	Datum  time.Time `json:"datum"`
	Aantal int       `json:"aantal"`

	Object   Object `json:"object"`
	ObjectID int    `json:"objectID"`

	Maat   *Maat `json:"maat"`
	MaatID int   `json:"maatID"`

	Opmerking string `json:"opmerking"`

	MateriaalID int `json:"materiaalID"`
	VoedingID   int `json:"voedingID"`

	// to be deprecated once we have a copy of family members in our DB
	Ontvanger string `json:"ontvanger"`
}

// FillGekregen loads all objects under an entry
func (c *Connection) FillGekregen(gekregen []OntvangEntry) error {
	if gekregen == nil {
		return nil
	}

	cache := map[int]Object{}

	for id, entry := range gekregen {
		object := Object{}
		if cachedObject, hasCache := cache[entry.ObjectID]; hasCache {
			object = cachedObject
		} else {
			res := c.Preload(clause.Associations).Where("id = ?", entry.ObjectID).First(&object)
			if res.Error != nil {
				return res.Error
			}
			cache[entry.ObjectID] = object
		}

		gekregen[id].Object = object
	}

	maatCache := map[int]Maat{}
	for id, entry := range gekregen {
		if entry.MaatID == 0 {
			continue
		}
		maat := Maat{}
		if cachedObject, hasCache := maatCache[entry.MaatID]; hasCache {
			maat = cachedObject
		} else {
			res := c.Preload(clause.Associations).Where("id = ?", entry.MaatID).First(&maat)
			if res.Error != nil {
				return res.Error
			}
			maatCache[entry.MaatID] = maat
		}

		gekregen[id].Maat = &maat
	}

	return nil
}
