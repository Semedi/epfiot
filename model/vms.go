package model

import (
	"context"

	graphql "github.com/graph-gophers/graphql-go"
	"github.com/jinzhu/gorm"
)

// Vm is the base type for Vms to be used by the db and gql
type Vm struct {
	gorm.Model
	OwnerID uint
	Name    string
	Tags    []Tag `gorm:"many2many:vms_tags"`
}

// VmResolver contains the DB and the model for resolving
type VmResolver struct {
	db *DB
	m  Vm
}

// ID resolves the ID field for Vm
func (p *VmResolver) ID(ctx context.Context) *graphql.ID {
	return gqlIDP(p.m.ID)
}

// Owner resolves the owner field for Vm
func (p *VmResolver) Owner() (*UserResolver, error) {
	user, err := p.db.getVmOwner(int32(p.m.OwnerID))
	if err != nil {
		return nil, err
	}

	r := UserResolver{
		db: p.db,
		m:  *user,
	}

	return &r, nil
}

// Name resolves the name field for Vm
func (p *VmResolver) Name(ctx context.Context) *string {
	return &p.m.Name
}

// Tags resolves the vm tags
func (p *VmResolver) Tags(ctx context.Context) (*[]*TagResolver, error) {
	tags, err := p.db.getVmTags(&p.m)
	if err != nil {
		return nil, err
	}

	r := make([]*TagResolver, len(tags))
	for i := range tags {
		r[i] = &TagResolver{
			db: p.db,
			m:  tags[i],
		}
	}

	return &r, nil
}
