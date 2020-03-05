package core

import (
	"context"
	graphql "github.com/graph-gophers/graphql-go"
	"github.com/semedi/epfiot/core/model"
	"github.com/semedi/epfiot/driver"
)

type Resolver struct {
	Db         *DB
	Controller driver.Provider
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
	Config   *model.ConfigInput
}

type thingInput struct {
	ID   *graphql.ID
	Name string
	Info string
}

// GetUser resolves the getUser query
func (r *Resolver) GetUser(ctx context.Context, args struct{ ID graphql.ID }) (*UserResolver, error) {

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

func (r *Resolver) updateVms() error {

	vms, err := r.Db.getVms()
	if err != nil {
		return err
	}

	for i := range vms {
		err := r.Controller.Update(&vms[i])
		if err != nil {
			return err
		}

		// lock
		r.Db.Savevm(&vms[i])
		// endlock

	}

	return nil
}

func (r *Resolver) GetVms(ctx context.Context) (*[]*VmResolver, error) {
	id := ctx.Value("userid").(uint)

	vms, err := r.Db.GetUserVms(id)
	if err != nil {
		return nil, err
	}

	v := make([]*VmResolver, len(vms))
	// concurrency point
	for i := range vms {

		err := r.Controller.Update(&vms[i])
		if err != nil {
			return nil, err
		}

		// lock
		r.Db.Savevm(&vms[i])
		// endlock

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

// ddVm Resolves the createvm mutation
func (r *Resolver) CreateVm(ctx context.Context, args struct{ Vm vmInput }) (*VmResolver, error) {
	id := ctx.Value("userid").(uint)

	// concurrency point:
	driver.Copy_base(args.Vm.Base, id, args.Vm.Name)

	vm, err := r.Db.addVm(args.Vm, id)
	if err != nil {
		return nil, err
	}
	// end concurrency

	config_path, err := driver.Create_config(id, args.Vm.Name, args.Vm.Config)
	if err != nil {
		return nil, err
	}

	r.Controller.Create(*vm, id, config_path)

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

	err = r.Controller.PowerOn(query)
	if err != nil {
		return nil, err
	}

	r.Controller.Update(vm)
	r.Db.Savevm(vm)

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
	err = r.Controller.Shutdown(query)
	if err != nil {
		return nil, err
	}

	vm.State = "POWEROFF"

	return &VmResolver{
		db: r.Db,
		m:  *vm,
	}, nil

}

//refactor
func (r *Resolver) ForceOFF(args struct{ VmID graphql.ID }) (*bool, error) {
	vmID, err := gqlIDToUint(args.VmID)
	b := false

	if err != nil {
		return &b, err
	}

	vm, err := r.Db.getVm(vmID)
	if err != nil {
		return &b, err
	}

	query := vm.Name
	err = r.Controller.ForceOFF(query)
	if err != nil {
		return &b, err
	}

	b = true
	return &b, nil
}

//refactor
func (r *Resolver) ForceDestroyVM(ctx context.Context, args struct{ VmID graphql.ID }) (*bool, error) {
	id := ctx.Value("userid").(uint)

	vmID, err := gqlIDToUint(args.VmID)
	b := false

	if err != nil {
		return &b, err
	}

	vm, err := r.Db.getVm(vmID)
	if err != nil {
		return &b, err
	}

	query := vm.Name
	err = r.Controller.ForceDestroy(query)
	if err != nil {
		return &b, err
	}

	driver.Erasefiles(id, query)

	b = true
	return &b, nil

}

// TODO:
//	ERASE VM FROM DATABASE
//  DELETE DISK IN DRIVER OPERATION
func (r *Resolver) DestroyVM(ctx context.Context, args struct{ VmID graphql.ID }) (*bool, error) {
	id := ctx.Value("userid").(uint)

	vmID, err := gqlIDToUint(args.VmID)
	b := false

	if err != nil {
		return &b, err
	}

	vm, err := r.Db.getVm(vmID)
	if err != nil {
		return &b, err
	}

	query := vm.Name
	err = r.Controller.Destroy(query)
	if err != nil {
		return &b, err
	}

	driver.Erasefiles(id, query)

	b = true
	return &b, nil

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

func (r *Resolver) get_vm_device(DevID, VmID graphql.ID) (*model.Vm, *model.Hostdev, error) {

	vmID, err := gqlIDToUint(VmID)
	if err != nil {
		return nil, nil, err
	}

	devID, err := gqlIDToUint(DevID)
	if err != nil {
		return nil, nil, err
	}

	dev, err := r.Db.getDev(devID)
	if err != nil {
		return nil, nil, err
	}

	vm, err := r.Db.getVm(vmID)
	if err != nil {
		return nil, nil, err
	}

	return vm, dev, nil
}

func (r *Resolver) DetachDevice(args struct{ DevID, VmID graphql.ID }) (*bool, error) {

	vm, dev, err := r.get_vm_device(args.DevID, args.VmID)
	if err != nil {
		return nil, err
	}

	err = r.Controller.DetachDevice(*vm, *dev)
	if err != nil {
		return nil, err
	}

	dev.VmID = vm.Model.ID
	r.Db.SaveDev(dev)

	b := true
	return &b, nil
}

// TODO:
// send udp request to bootstrap only if IP
func (r *Resolver) AttachDevice(args struct{ DevID, VmID graphql.ID }) (*bool, error) {
	vm, dev, err := r.get_vm_device(args.DevID, args.VmID)
	if err != nil {
		return nil, err
	}

	err = r.Controller.AttachDevice(*vm, *dev)
	if err != nil {
		return nil, err
	}

	dev.VmID = vm.Model.ID
	r.Db.SaveDev(dev)

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
