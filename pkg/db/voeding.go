package db

import (
	"gorm.io/gorm"
)

type Voeding struct {
	gorm.Model
	MVMNummer       string         `json:"mvmNummer" gorm:"column:mvm_nummer"`
	Opmerking       string         `json:"opmerking"`
	SpecialeVoeding string         `json:"specialeVoeding"`
	Gekregen        []OntvangEntry `json:"gekregen"`
}
