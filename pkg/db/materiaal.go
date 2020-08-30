package db

import (
	"time"

	"github.com/jinzhu/gorm"
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
	Naam      string            `json:"naam"`
	Categorie MateriaalCategory `json:"categorie"`
}

type Materiaal struct {
	gorm.Model
	MVMNummer string           `json:"mvmNummer"`
	Opmerking string           `json:"opmerking"`
	Gekregen  []MateriaalEntry `json:"gekregen"`
}

type MateriaalEntry struct {
	gorm.Model
	Datum     time.Time       `json:"datum"`
	Aantal    int             `json:"aantal"`
	Object    MateriaalObject `json:"object"`
	Maat      string          `json:"maat"`
	Opmerking string          `json:"opmerking"`
	// to be deprecated once we have a copy of family members in our DB
	Ontvanger string `json:"ontvanger"`
}
