package service

import (
	"context"

	graphql "github.com/graph-gophers/graphql-go"
	"github.com/jinzhu/gorm"
)

// Tag is the base type for a vm tag to be used by the db and gql
type Tag struct {
	gorm.Model
	Title string
	Vms   []Vm `gorm:"many2many:vm_tags"`
}

// TagResolver contains the db and the Tag model for resolving
type TagResolver struct {
	db *DB
	m  Tag
}

// ID resolves the ID for Tag
func (t *TagResolver) ID(ctx context.Context) *graphql.ID {
	return gqlIDP(t.m.ID)
}

// Title resolves the title field
func (t *TagResolver) Title(ctx context.Context) *string {
	return &t.m.Title
}

// Vms resolves the vmsvnoremap field
func (t *TagResolver) Vms(ctx context.Context) (*[]*VmResolver, error) {
	vms, err := t.db.getTagVms(ctx, &t.m)
	if err != nil {
		return nil, err
	}

	r := make([]*VmResolver, len(vms))
	for i := range vms {
		r[i] = &VmResolver{
			db: t.db,
			m:  vms[i],
		}
	}

	return &r, nil
}
