package model

import (
	"github.com/jinzhu/gorm"
)

// Tag is the base type for a vm tag to be used by the db and gql
type Tag struct {
	gorm.Model
	Title string
	Vms   []Vm `gorm:"many2many:vm_tags"`
}

