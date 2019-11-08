package core

import (
	"context"
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"

	graphql "github.com/graph-gophers/graphql-go"
	"github.com/semedi/epfiot/core/model"
	"github.com/semedi/epfiot/driver"
)

type Resolver struct {
	Db         *DB
	Controller *driver.Controller
}

// GetUser resolves the getUser query
func (r *Resolver) GetUser(ctx context.Context, args struct{ ID graphql.ID }) (*UserResolver, error) {

	r.Controller.Handler.Init()

	id, err := gqlIDToUint(args.ID)
	if err != nil {
		return nil, err
	}

	user, err := r.Db.getUser(id)
	if err != nil {
		return nil, err
	}

	s := UserResolver{
		db: r.Db,
		m:  *user,
	}

	return &s, nil
}

func (r *Resolver) GetVms(ctx context.Context) (*[]*VmResolver, error) {
	id := ctx.Value("userid").(uint)

	vms, err := r.Db.GetUserVms(id)
	if err != nil {
		return nil, err
	}

	v := make([]*VmResolver, len(vms))
	for i := range vms {
		v[i] = &VmResolver{
			db: r.Db,
			m:  vms[i],
		}
	}

	return &v, nil
}

func (r *Resolver) GetUsb(ctx context.Context) (*[]*DevResolver, error) {

	devs, err := r.Db.getHostdevices()
	if err != nil {
		return nil, err
	}

	d := make([]*DevResolver, len(devs))
	for i := range devs {
		d[i] = &DevResolver{
			db: r.Db,
			m:  devs[i],
		}
	}

	return &d, nil
}

func (r *Resolver) GetUsers(ctx context.Context) (*[]*UserResolver, error) {
	users, err := r.Db.getUsers()
	if err != nil {
		return nil, err
	}

	u := make([]*UserResolver, len(users))
	for i := range users {
		u[i] = &UserResolver{
			db: r.Db,
			m:  users[i],
		}
	}

	return &u, nil
}

// GetVm resolves the getVm query
func (r *Resolver) GetVm(args struct{ ID graphql.ID }) (*VmResolver, error) {
	id, err := gqlIDToUint(args.ID)
	if err != nil {
		return nil, err
	}

	vm, err := r.Db.getVm(id)
	if err != nil {
		return nil, err
	}

	s := VmResolver{
		db: r.Db,
		m:  *vm,
	}

	return &s, nil
}

// GetDev resolves the getDev query
func (r *Resolver) GetDev(ctx context.Context, args struct{ ID graphql.ID }) (*DevResolver, error) {
	id, err := gqlIDToUint(args.ID)
	if err != nil {
		return nil, err
	}

	dev, err := r.Db.getDev(id)
	if err != nil {
		return nil, err
	}

	s := DevResolver{
		db: r.Db,
		m:  *dev,
	}

	return &s, nil
}

func (r *Resolver) GetThing(ctx context.Context, args struct{ ID graphql.ID }) (*ThingResolver, error) {
	id, err := gqlIDToUint(args.ID)
	if err != nil {
		return nil, err
	}

	thing, err := r.Db.getThing(id)
	if err != nil {
		return nil, err
	}

	s := ThingResolver{
		db: r.Db,
		m:  *thing,
	}

	return &s, nil
}

// vmInput has everything needed to do adds and updates on a vm
type vmInput struct {
	ID       *graphql.ID
	OwnerID  int32
	Base     string
	Name     string
	Memory   int32
	Vcpu     int32
	DevIDs   *[]*int32
	ThingIDs *[]*int32
}

type thingInput struct {
	ID   *graphql.ID
	Name string
	Info string
}

// ddVm Resolves the createvm mutation
func (r *Resolver) CreateVm(ctx context.Context, args struct{ Vm vmInput }) (*VmResolver, error) {
	id := ctx.Value("userid").(uint)

	driver.Copy_base(args.Vm.Base, id, args.Vm.Name)

	vm, err := r.Db.addVm(args.Vm, id)
	if err != nil {
		return nil, err
	}
	r.Controller.Handler.Create(*vm, id)

	s := VmResolver{
		db: r.Db,
		m:  *vm,
	}

	return &s, nil
}

// UpdateVm takes care of updating any field on the vm
func (r *Resolver) UpdateVm(args struct{ Vm vmInput }) (*VmResolver, error) {
	vm, err := r.Db.updateVm(&args.Vm)
	if err != nil {
		return nil, err
	}

	s := VmResolver{
		db: r.Db,
		m:  *vm,
	}

	return &s, nil
}

func (r *Resolver) PowerON(args struct{ VmID graphql.ID }) (*VmResolver, error) {
	vmID, err := gqlIDToUint(args.VmID)
	if err != nil {
		return nil, err
	}

	vm, err := r.Db.getVm(vmID)
	if err != nil {
		return nil, err
	}

	query := vm.Name

	err = r.Controller.Handler.PowerOn(query)
	if err != nil {
		return nil, err
	}

	return &VmResolver{
		db: r.Db,
		m:  *vm,
	}, nil

}

func (r *Resolver) PowerOFF(args struct{ VmID graphql.ID }) (*VmResolver, error) {
	vmID, err := gqlIDToUint(args.VmID)
	if err != nil {
		return nil, err
	}

	vm, err := r.Db.getVm(vmID)
	if err != nil {
		return nil, err
	}

	query := vm.Name
	err = r.Controller.Handler.Shutdown(query)
	if err != nil {
		return nil, err
	}

	return &VmResolver{
		db: r.Db,
		m:  *vm,
	}, nil

}

// TODO:
// send udp request to bootstrap only if IP
func (r *Resolver) AttachThing(args struct{ ThingID, VmID graphql.ID }) (*bool, error) {
	vmID, err := gqlIDToUint(args.VmID)
	if err != nil {
		return nil, err
	}

	thingID, err := gqlIDToUint(args.ThingID)
	if err != nil {
		return nil, err
	}

	thing, err := r.Db.getThing(thingID)
	if err != nil {
		return nil, err
	}

	thing.VmID = vmID

	r.Db.SaveThing(thing)

	b := true

	return &b, nil
}

// DeleteVm takes care of deleting a vm record
func (r *Resolver) DeleteVm(args struct{ UserID, VmID graphql.ID }) (*bool, error) {
	vmID, err := gqlIDToUint(args.VmID)
	if err != nil {
		return nil, err
	}

	userID, err := gqlIDToUint(args.UserID)
	if err != nil {
		return nil, err
	}

	return r.Db.deleteVm(userID, vmID)
}

func (r *Resolver) CreateThing(ctx context.Context, args struct{ Thing thingInput }) (*ThingResolver, error) {
	t, err := r.Db.AddThing(args.Thing)

	if err != nil {
		return nil, err
	}

	s := ThingResolver{
		db: r.Db,
		m:  *t,
	}

	return &s, nil
}

// TODO:
// send udp request to bootstrap if have ip
func (r *Resolver) CreateThingVm(ctx context.Context, args struct {
	Thing thingInput
	VmID  graphql.ID
}) (*ThingResolver, error) {
	vmID, err := gqlIDToUint(args.VmID)
	if err != nil {
		return nil, err
	}

	t, err := r.Db.AddThing(args.Thing)

	if err != nil {
		return nil, err
	}

	vm, err := r.Db.getVm(vmID)
	if err != nil {
		return nil, err
	}

	vm.Things = append(vm.Things, *t)
	r.Db.Savevm(vm)

	s := ThingResolver{
		db: r.Db,
		m:  *t,
	}

	return &s, nil
}

// encode cursor encodes the cursot position in base64
func encodeCursor(i int) graphql.ID {
	return graphql.ID(base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("cursor%d", i))))
}

// decode cursor decodes the base 64 encoded cursor and resturns the integer
func decodeCursor(s string) (int, error) {
	b, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return 0, err
	}

	i, err := strconv.Atoi(strings.TrimPrefix(string(b), "cursor"))
	if err != nil {
		return 0, err
	}

	return i, nil
}

// UserResolver contains the database and the user model to resolve against
type UserResolver struct {
	db *DB
	m  model.User
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

func gqlIDToUint(i graphql.ID) (uint, error) {
	r, err := strconv.ParseInt(string(i), 10, 32)
	if err != nil {
		return 0, err
	}

	return uint(r), nil
}

func gqlIDP(id uint) *graphql.ID {
	r := graphql.ID(fmt.Sprint(id))
	return &r
}
