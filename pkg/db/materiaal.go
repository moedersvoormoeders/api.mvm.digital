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

type MateriaalMaat struct {
	gorm.Model
	Naam              string `json:"naam"`
	MateriaalObjectID uint
}

type MateriaalObject struct {
	gorm.Model
	Naam        string `json:"naam"`
	CategorieID int
	Categorie   MateriaalCategory `json:"categorie"`
	Hidden      bool              `json:"hidden"`
	Prijs       float64           `json:"prijs"` // most of the times this is 0.0
	Maten       []MateriaalMaat   `json:"maten"`
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
	Maat      MateriaalMaat   `json:"maat"`
	Opmerking string          `json:"opmerking"`
	// to be deprecated once we have a copy of family members in our DB
	Ontvanger   string `json:"ontvanger"`
	MateriaalID uint
	MaatID      uint
}

// CleanMateriaalGekregen cleans out too deep references not to trigger UPDATE on these fields
func (c *Connection) CleanMateriaalGekregen(o *Materiaal) {
	for i := range o.Gekregen {
		if o.Gekregen[i].Maat.ID > 0 {
			o.Gekregen[i].MaatID = o.Gekregen[i].Maat.ID
			o.Gekregen[i].Maat = MateriaalMaat{}
		}
		if o.Gekregen[i].Object.ID > 0 {
			o.Gekregen[i].ObjectID = int(o.Gekregen[i].Object.ID)
			o.Gekregen[i].Object = MateriaalObject{}
		}
	}
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

	maatCache := map[uint]MateriaalMaat{}
	for id, entry := range o.Gekregen {
		maat := MateriaalMaat{}
		if cachedObject, hasCache := maatCache[entry.MaatID]; hasCache {
			maat = cachedObject
		} else {
			res := c.Preload(clause.Associations).Where("id = ?", entry.MaatID).First(&maat)
			if res.Error != nil {
				return res.Error
			}
			maatCache[entry.MaatID] = maat
		}

		o.Gekregen[id].Maat = maat
	}

	return nil
}
