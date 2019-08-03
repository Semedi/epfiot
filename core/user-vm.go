package core

import (
	"context"
	"encoding/base64"
	"strconv"
	"strings"

	graphql "github.com/graph-gophers/graphql-go"
)

type vmsConnArgs struct {
	First *int32
	After *graphql.ID
}

// VmsConnection returns nodes (vms) connected by edges (relationships)
func (u *UserResolver) VmsConnection(ctx context.Context, args vmsConnArgs) (*UserVmsConnectionResolver, error) {
	// query only the ID fields from the vms otherwise it would be wasteful
	ids, err := u.db.getUserVmIDs(u.m.ID)
	if err != nil {
		return nil, err
	}

	from := 0
	if args.After != nil {
		b, err := base64.StdEncoding.DecodeString(string(*args.After))
		if err != nil {
			return nil, err
		}
		i, err := strconv.Atoi(strings.TrimPrefix(string(b), "cursor"))
		if err != nil {
			return nil, err
		}
		from = i + 1
	}

	to := len(ids)
	if args.First != nil {
		to = from + int(*args.First)
		if to > len(ids) {
			to = len(ids)
		}
	}

	upc := UserVmsConnectionResolver{
		ids:  ids,
		from: from,
		to:   to,
	}
	return &upc, nil
}

// UserVmEdge is an edge (related node) that is returned in pagination
type UserVmEdge struct {
	cursor graphql.ID
	node   VmResolver
}

// Cursor resolves the cursor for pagination
func (u *UserVmEdge) Cursor(ctx context.Context) graphql.ID {
	return u.cursor
}

// Node resolves the node for pagination
func (u *UserVmEdge) Node(ctx context.Context) *VmResolver {
	return &u.node
}

// PageInfo gives page info for pagination
type PageInfo struct {
	startCursor     graphql.ID
	endCursor       graphql.ID
	hasNextPage     bool
	hasPreviousPage bool
}

// StartCursor ...
func (u *PageInfo) StartCursor(ctx context.Context) *graphql.ID {
	return &u.startCursor
}

// EndCursor ...
func (u *PageInfo) EndCursor(ctx context.Context) *graphql.ID {
	return &u.endCursor
}

// HasNextPage returns true if there are more results to show
func (u *PageInfo) HasNextPage(ctx context.Context) bool {
	return u.hasNextPage
}

// HasPreviousPage returns true if there are results behind the current cursor position
func (u *PageInfo) HasPreviousPage(ctx context.Context) bool {
	return u.hasPreviousPage
}

// UserVmsConnectionResolver is all the vms that are connected to a certain user
type UserVmsConnectionResolver struct {
	db   *DB
	ids  []int
	from int
	to   int
}

// TotalCount gives the total amount of vms in UserVmsConnection
func (u UserVmsConnectionResolver) TotalCount(ctx context.Context) int32 {
	return int32(len(u.ids))
}

// Edges gives a list of all the edges (related vms) that belong to a user
func (u *UserVmsConnectionResolver) Edges(ctx context.Context) (*[]*UserVmEdge, error) {
	// query goes here because I know all of the ids that are needed. If I queried in the
	// UserVmEdge resolver method, it would run multiple single queries
	vms, err := u.db.getVmsByID(u.ids, u.from, u.to)
	if err != nil {
		return nil, err
	}

	l := make([]*UserVmEdge, u.to-u.from)
	for i := range l {
		l[i] = &UserVmEdge{
			cursor: encodeCursor(u.from + i),
			node: VmResolver{
				db: u.db,
				m:  vms[i],
			},
		}
	}

	return &l, nil
}

// PageInfo resolves page info
func (u *UserVmsConnectionResolver) PageInfo(ctx context.Context) (*PageInfo, error) {
	p := PageInfo{
		startCursor:     encodeCursor(u.from),
		endCursor:       encodeCursor(u.to - 1),
		hasNextPage:     u.to < len(u.ids),
		hasPreviousPage: u.from > 0,
	}
	return &p, nil
}
