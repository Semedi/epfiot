package model

import (
	"github.com/jinzhu/gorm"
)

// Hostdev is the base type for a vm host device to be used by the db and gql
type Hostdev struct {
	gorm.Model
	Bus    string
	Device string
	VmID   uint
}
