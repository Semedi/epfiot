package model

import (
	"context"
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"

	graphql "github.com/graph-gophers/graphql-go"
	"github.com/semedi/epfiot/driver"
)


type Resolver struct {
	Db *DB
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

// GetTag resolves the getTag query
func (r *Resolver) GetTag(ctx context.Context, args struct{ Title string }) (*TagResolver, error) {
	tag, err := r.Db.getTagBytTitle(args.Title)
	if err != nil {
		return nil, err
	}

	s := TagResolver{
		db: r.Db,
		m:  *tag,
	}

	return &s, nil
}

// vmInput has everything needed to do adds and updates on a vm
type vmInput struct {
	ID      *graphql.ID
	OwnerID int32
    Base    string
	Name    string
	TagIDs  *[]*int32
}

// ddVm Resolves the createvm mutation
func (r *Resolver) CreateVm(ctx context.Context, args struct{ Vm vmInput }) (*VmResolver, error) {
    id := ctx.Value("userid").(uint)

    driver.Copy_base(args.Vm.Base, id, args.Vm.Name)

	vm, err := r.Db.addVm(args.Vm, id)
	if err != nil {
		return nil, err
	}

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
