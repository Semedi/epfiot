package model

// TEST DATA TO BE PUT INTO THE DB
var users = []User{
	User{Name: "Alice"},
	User{Name: "Bob"},
	User{Name: "Charlie"},
}

// Since the db is torn down and created on every run, I know the above users will have
// ID's 1, 2, 3
var vms = []Vm{
	Vm{Name: "debian7", OwnerID: 1},
	Vm{Name: "debian8", OwnerID: 1},
	Vm{Name: "ubuntu", OwnerID: 1},
	Vm{Name: "archlinux", OwnerID: 1},
	Vm{Name: "coreOs", OwnerID: 1},
	Vm{Name: "Centos", OwnerID: 1},
	Vm{Name: "Manjaro", OwnerID: 1},
	Vm{Name: "linuxmint", OwnerID: 2},
	Vm{Name: "TinyCore", OwnerID: 3},
	Vm{Name: "void", OwnerID: 3},
}

// Tags to be put in the database
var tags = []Tag{
	Tag{Title: "amd64"},
	Tag{Title: "x86"},
	Tag{Title: "arm"},
	Tag{Title: "gpu"},
}
