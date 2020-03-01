package core

import (
	"context"
	graphql "github.com/graph-gophers/graphql-go"
	"github.com/semedi/epfiot/core/model"
)

type VmResolver struct {
	db *DB
	m  model.Vm
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

// memory resolves the memory field for Vm
func (p *VmResolver) Memory(ctx context.Context) *int32 {
	r := int32(p.m.Memory)
	return &r
}

// vcpu resolves the vcpu field for Vm
func (p *VmResolver) Vcpu(ctx context.Context) *int32 {
	r := int32(p.m.Vcpu)
	return &r
}

// Base resolves the base field for Vm
func (p *VmResolver) Base(ctx context.Context) *string {
	return &p.m.Base
}

func (p *VmResolver) Ip(ctx context.Context) *string {
	return &p.m.Ip
}

// Dev resolves the vm devices
func (p *VmResolver) Dev(ctx context.Context) (*[]*DevResolver, error) {
	devices, err := p.db.getVmDev(&p.m)
	if err != nil {
		return nil, err
	}

	r := make([]*DevResolver, len(devices))
	for i := range devices {
		r[i] = &DevResolver{
			db: p.db,
			m:  devices[i],
		}
	}

	return &r, nil
}

func (p *VmResolver) Things(ctx context.Context) (*[]*ThingResolver, error) {
	things, err := p.db.getVmThings(&p.m)
	if err != nil {
		return nil, err
	}

	r := make([]*ThingResolver, len(things))
	for i := range things {
		r[i] = &ThingResolver{
			db: p.db,
			m:  things[i],
		}
	}

	return &r, nil
}
