package db

import (
	"time"

	"gorm.io/gorm/clause"

	"gorm.io/gorm"
)

type MateriaalCategory struct {
	gorm.Model
	Naam    string `json:"naam"`
	PerKind bool   `json:"perKind"`
	OpMaat  bool   `json:"opMaat"`
	Order   int    `json:"order"`
}

type MateriaalObject struct {
	gorm.Model
	Naam        string `json:"naam"`
	CategorieID int
	Categorie   MateriaalCategory `json:"categorie"`
	Hidden      bool              `json:"hidden"`
	Prijs       float64           `json:"prijs"` // most of the times this is 0.0
}

type Materiaal struct {
	gorm.Model
	MVMNummer string           `json:"mvmNummer" gorm:"column:mvm_nummer"`
	Opmerking string           `json:"opmerking"`
	Gekregen  []MateriaalEntry `json:"gekregen"`
}

type MateriaalEntry struct {
	gorm.Model
	Datum     time.Time       `json:"datum"`
	Aantal    int             `json:"aantal"`
	ObjectID  int             `json:"objectId"`
	Object    MateriaalObject `json:"object"`
	Maat      string          `json:"maat"`
	Opmerking string          `json:"opmerking"`
	// to be deprecated once we have a copy of family members in our DB
	Ontvanger   string `json:"ontvanger"`
	MateriaalID uint
}

func (c *Connection) FillMateriaalGekregen(o *Materiaal) error {
	if o.Gekregen == nil {
		return nil
	}

	cache := map[int]MateriaalObject{}

	for id, entry := range o.Gekregen {
		object := MateriaalObject{}
		if cachedObject, hasCache := cache[entry.ObjectID]; hasCache {
			object = cachedObject
		} else {
			res := c.Preload(clause.Associations).Where("id = ?", entry.ObjectID).First(&object)
			if res.Error != nil {
				return res.Error
			}
			cache[entry.ObjectID] = object
		}

		o.Gekregen[id].Object = object
	}

	return nil
}
