package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"

	graphql "github.com/graph-gophers/graphql-go"
)

type Resolver struct {
	db *DB
}

// GetUser resolves the getUser query
func (r *Resolver) GetUser(ctx context.Context, args struct{ ID graphql.ID }) (*UserResolver, error) {
	id, err := gqlIDToUint(args.ID)
	if err != nil {
		return nil, err
	}

	user, err := r.db.getUser(ctx, id)
	if err != nil {
		return nil, err
	}

	s := UserResolver{
		db: r.db,
		m:  *user,
	}

	return &s, nil
}

func (r *Resolver) GetUsers(ctx context.Context) (*[]*UserResolver, error) {
	users, err := r.db.getUsers(ctx)
	if err != nil {
		return nil, err
	}

	u := make([]*UserResolver, len(users))
	for i := range users {
		u[i] = &UserResolver{
			db: r.db,
			m:  users[i],
		}
	}

	return &u, nil
}

// GetVm resolves the getVm query
func (r *Resolver) GetVm(ctx context.Context, args struct{ ID graphql.ID }) (*VmResolver, error) {
	id, err := gqlIDToUint(args.ID)
	if err != nil {
		return nil, err
	}

	vm, err := r.db.getVm(ctx, id)
	if err != nil {
		return nil, err
	}

	s := VmResolver{
		db: r.db,
		m:  *vm,
	}

	return &s, nil
}

// GetTag resolves the getTag query
func (r *Resolver) GetTag(ctx context.Context, args struct{ Title string }) (*TagResolver, error) {
	tag, err := r.db.getTagBytTitle(ctx, args.Title)
	if err != nil {
		return nil, err
	}

	s := TagResolver{
		db: r.db,
		m:  *tag,
	}

	return &s, nil
}

// vmInput has everything needed to do adds and updates on a vm
type vmInput struct {
	ID      *graphql.ID
	OwnerID int32
	Name    string
	TagIDs  *[]*int32
}

// AddVm Resolves the addVm mutation
func (r *Resolver) AddVm(ctx context.Context, args struct{ Vm vmInput }) (*VmResolver, error) {
	vm, err := r.db.addVm(ctx, args.Vm)
	if err != nil {
		return nil, err
	}

	s := VmResolver{
		db: r.db,
		m:  *vm,
	}

	return &s, nil
}

// UpdateVm takes care of updating any field on the vm
func (r *Resolver) UpdateVm(ctx context.Context, args struct{ Vm vmInput }) (*VmResolver, error) {
	vm, err := r.db.updateVm(ctx, &args.Vm)
	if err != nil {
		return nil, err
	}

	s := VmResolver{
		db: r.db,
		m:  *vm,
	}

	return &s, nil
}

// DeleteVm takes care of deleting a vm record
func (r *Resolver) DeleteVm(ctx context.Context, args struct{ UserID, VmID graphql.ID }) (*bool, error) {
	vmID, err := gqlIDToUint(args.VmID)
	if err != nil {
		return nil, err
	}

	userID, err := gqlIDToUint(args.UserID)
	if err != nil {
		return nil, err
	}

	return r.db.deleteVm(ctx, userID, vmID)
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
