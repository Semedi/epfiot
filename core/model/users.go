package model

import (
	"github.com/jinzhu/gorm"
)

// User is the base user model to be used throughout the app
type User struct {
	gorm.Model
	Name string
	Vms  []Vm `gorm:"foreignkey:OwnerID"`
}
