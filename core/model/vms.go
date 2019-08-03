package model

import (
	"github.com/jinzhu/gorm"
)

// Vm is the base type for Vms to be used by the db and gql
type Vm struct {
	gorm.Model
	OwnerID uint
    Base    string
	Name    string
	Tags    []Tag `gorm:"many2many:vms_tags"`
}
