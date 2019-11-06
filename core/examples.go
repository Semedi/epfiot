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

// Devices to be put in the database
var hdevices = []model.Hostdev{
	model.Hostdev{Bus: "4", Device: "1"},
	model.Hostdev{Bus: "3", Device: "1"},
}

var hdevices2 = []model.Hostdev{
	model.Hostdev{Bus: "1", Device: "1"},
	model.Hostdev{Bus: "2", Device: "1"},
}
