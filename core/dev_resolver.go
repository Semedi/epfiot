package core

import (
	"context"
	graphql "github.com/graph-gophers/graphql-go"
	"github.com/semedi/epfiot/core/model"
)

// DevResolver contains the db and the Hostdevice model for resolving
type DevResolver struct {
	db *DB
	m  model.Hostdev
}

// ID resolves the ID for device
func (d *DevResolver) ID(ctx context.Context) *graphql.ID {
	return gqlIDP(d.m.ID)
}

// Bus resolves the bus field
func (d *DevResolver) Bus(ctx context.Context) *string {
	return &d.m.Bus
}

// Device resolves the device field
func (d *DevResolver) Device(ctx context.Context) *string {
	return &d.m.Device
}

// Info resolves the info field
func (d *DevResolver) Info(ctx context.Context) *string {
	return &d.m.Info
}

//// Vms resolves the vmsvnoremap field
//func (t *DevResolver) Vms(ctx context.Context) (*[]*VmResolver, error) {
//	vms, err := t.db.getTagVms(&t.m)
//	if err != nil {
//		return nil, err
//	}
//
//	r := make([]*VmResolver, len(vms))
//	for i := range vms {
//		r[i] = &VmResolver{
//			db: t.db,
//			m:  vms[i],
//		}
//	}
//
//	return &r, nil
//}

// VmResolver contains the DB and the model for resolving
