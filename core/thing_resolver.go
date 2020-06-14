package core

import (
	"context"

	graphql "github.com/graph-gophers/graphql-go"
	"github.com/semedi/epfiot/core/model"
)

type ThingResolver struct {
	db *DB
	m  model.Thing
}

// ID resolves the ID for Thing
func (t *ThingResolver) ID(ctx context.Context) *graphql.ID {
	return gqlIDP(t.m.ID)
}

// Name resolves the name field
func (t *ThingResolver) Name(ctx context.Context) *string {
	return &t.m.Name
}

// Info resolves the info field
func (t *ThingResolver) Info(ctx context.Context) *string {
	return &t.m.Info
}

func (t *ThingResolver) Server(ctx context.Context) *string {
	return &t.m.Server
}
