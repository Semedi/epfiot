package service

import (
	"context"

	graphql "github.com/graph-gophers/graphql-go"
	"github.com/jinzhu/gorm"
)

// User is the base user model to be used throughout the app
type User struct {
	gorm.Model
	Name string
	Vms  []Vm `gorm:"foreignkey:OwnerID"`
}

// UserResolver contains the database and the user model to resolve against
type UserResolver struct {
	db *DB
	m  User
}

// ID resolves the user ID
func (u *UserResolver) ID(ctx context.Context) *graphql.ID {
	return gqlIDP(u.m.ID)
}

// Name resolves the Name field for User, it is all caps to avoid name clashes
func (u *UserResolver) Name(ctx context.Context) *string {
	return &u.m.Name
}

// Vms resolves the Vms field for User
func (u *UserResolver) Vms(ctx context.Context) (*[]*VmResolver, error) {
	vms, err := u.db.GetUserVms(u.m.ID)
	if err != nil {
		return nil, err
	}

	r := make([]*VmResolver, len(vms))
	for i := range vms {
		r[i] = &VmResolver{
			db: u.db,
			m:  vms[i],
		}
	}

	return &r, nil
}
