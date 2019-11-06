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
	//model.Vm{Name: "debian8", Base: "null", OwnerID: 1},
	//model.Vm{Name: "archlinux", Base: "null", OwnerID: 1},
	//model.Vm{Name: "coreOs", Base: "null", OwnerID: 1},
	//model.Vm{Name: "Centos", Base: "null", OwnerID: 1},
	//model.Vm{Name: "Manjaro", Base: "null", OwnerID: 1},
	//model.Vm{Name: "linuxmint", Base: "null", OwnerID: 2},
	//model.Vm{Name: "TinyCore", Base: "null", OwnerID: 3},
	//model.Vm{Name: "void", Base: "null", OwnerID: 3},
}

// Tags to be put in the database
var tags = []model.Tag{
	model.Tag{Title: "arm"},
	model.Tag{Title: "amd64"},
}

var tags2 = []model.Tag{
	model.Tag{Title: "x86"},
}
