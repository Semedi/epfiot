package core

import "github.com/semedi/epfiot/core/model"

// TEST DATA TO BE PUT INTO THE DB
var users = []model.User{
	model.User{Name: "Alice@gmail.com"},
	model.User{Name: "Bob@gmail.com"},
	model.User{Name: "Charlie@gmail.com"},
}

// Since the db is torn down and created on every run, I know the above users will have
// ID's 1, 2, 3

var vms = []model.Vm{
	model.Vm{Name: "debian7", Base: "null", OwnerID: 1, Memory: 512, Vcpu: 1},
	model.Vm{Name: "ubuntu", Base: "null", OwnerID: 2},
}

var things = []model.Thing{
	model.Thing{Name: "test", Info: "sensor"},
}
