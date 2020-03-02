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
	Memory  int
	Vcpu    int
	Ip      string
	State   string
	Dev     []Hostdev `gorm:"foreignkey:VmID"`
	Things  []Thing   `gorm:"foreignkey:VmID"`
}

type ConfigInput struct {
	User     *string
	Password *string
}
