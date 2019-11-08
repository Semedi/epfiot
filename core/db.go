package core

import (
	"github.com/jinzhu/gorm"
	// nolint: gotype
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/semedi/epfiot/core/model"
)

type DB struct {
	DB *gorm.DB
}

// NewDB returns a new DB connection
func NewDB(path string) (*DB, error) {
	// connect to the example db, create if it doesn't exist
	db, err := gorm.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}

	// drop tables and all data, and recreate them fresh for this run
	db.DropTableIfExists(&model.User{}, &model.Vm{}, &model.Hostdev{}, &model.Thing{})
	db.AutoMigrate(&model.User{}, &model.Vm{}, &model.Hostdev{}, &model.Thing{})

	// put all the users into the db
	for _, u := range users {
		if err := db.Create(&u).Error; err != nil {
			return nil, err
		}
	}

	for _, t := range things {
		if err := db.Create(&t).Error; err != nil {
			return nil, err
		}
	}

	// put all the vms into the db
	for _, p := range vms {
		if err := db.Create(&p).Error; err != nil {
			return nil, err
		}
	}

	return &DB{db}, nil
}

func (db *DB) Savevm(vm *model.Vm) {
	db.DB.Save(vm)
}

func (db *DB) SaveThing(thing *model.Thing) {
	db.DB.Save(thing)
}

func (db *DB) CreateHostdevs(devices_str [][]string) error {
	BUS := 0
	DEV := 1
	INFO := 2

	for _, d := range devices_str {
		hdev := model.Hostdev{
			Bus:    d[BUS],
			Device: d[DEV],
			Info:   d[INFO],
		}

		if err := db.DB.Create(&hdev).Error; err != nil {
			return err
		}
	}

	return nil
}

// ###########################################################
// USERS:
// ###########################################################

func (db *DB) getUserVmIDs(userID uint) ([]int, error) {
	var ids []int
	err := db.DB.Where("owner_id = ?", userID).Find(&[]model.Vm{}).Pluck("id", &ids).Error
	if err != nil {
		return nil, err
	}
	return ids, nil
}

func (db *DB) Find_user(username string) (*model.User, error) {
	var user model.User

	err := db.DB.First(&user, "name = ?", username).Error
	if err != nil {
		return nil, err
	}

	if user.Name == "" {
		return nil, nil
	}

	return &user, nil
}

func (db *DB) getUser(id uint) (*model.User, error) {
	var user model.User
	err := db.DB.First(&user, id).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (db *DB) getUsers() ([]model.User, error) {
	var users []model.User
	err := db.DB.Find(&users).Error
	if err != nil {
		return nil, err
	}

	return users, nil
}

// GetUserVms gets vms associated with the user
func (db *DB) GetUserVms(id uint) ([]model.Vm, error) {
	var u model.User
	u.ID = id

	var p []model.Vm
	err := db.DB.Model(&u).Association("Vms").Find(&p).Error
	if err != nil {
		return nil, err
	}

	return p, nil
}

func (db *DB) getHostdevices() ([]model.Hostdev, error) {
	var devices []model.Hostdev
	err := db.DB.Find(&devices).Error
	if err != nil {
		return nil, err
	}

	return devices, nil
}

// ###########################################################
// VMS:
// ###########################################################

// GetVm should authorize the user and  return a vm or error
func (db *DB) getVm(id uint) (*model.Vm, error) {
	var p model.Vm
	err := db.DB.First(&p, id).Error
	if err != nil {
		return nil, err
	}

	return &p, nil
}

func (db *DB) getVmOwner(id int32) (*model.User, error) {
	var u model.User
	err := db.DB.First(&u, id).Error
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (db *DB) getVmDev(p *model.Vm) ([]model.Hostdev, error) {
	var devices []model.Hostdev
	err := db.DB.Model(p).Related(&devices).Error
	if err != nil {
		return nil, err
	}

	return devices, nil
}

func (db *DB) getVmThings(p *model.Vm) ([]model.Thing, error) {
	var things []model.Thing
	err := db.DB.Model(p).Related(&things).Error
	if err != nil {
		return nil, err
	}

	return things, nil
}

func (db *DB) getVmsByID(ids []int, from, to int) ([]model.Vm, error) {
	var p []model.Vm
	err := db.DB.Where("id in (?)", ids[from:to]).Find(&p).Error
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (db *DB) updateVm(args *vmInput) (*model.Vm, error) {
	// get the vm to be updated from the db
	var p model.Vm
	err := db.DB.First(&p, args.ID).Error
	if err != nil {
		return nil, err
	}

	// so the pointer dereference is safe
	if args.DevIDs == nil {
		return nil, err
	}

	// if there are devices to be updated, go through that process
	var newDevices []model.Hostdev
	if len(*args.DevIDs) > 0 {
		err = db.DB.Where("id in (?)", *args.DevIDs).Find(&newDevices).Error
		if err != nil {
			return nil, err
		}

		// replace the devices set with the new one
		err = db.DB.Model(&p).Association("Dev").Replace(newDevices).Error
		if err != nil {
			return nil, err
		}
	}

	// so the pointer dereference is safe
	if args.ThingIDs == nil {
		return nil, err
	}

	// if there are things to be updated, go through that process
	var newThings []model.Thing
	if len(*args.ThingIDs) > 0 {
		err = db.DB.Where("id in (?)", *args.ThingIDs).Find(&newThings).Error
		if err != nil {
			return nil, err
		}

		// replace the things set with the new one
		err = db.DB.Model(&p).Association("Things").Replace(newThings).Error
		if err != nil {
			return nil, err
		}
	}

	updated := model.Vm{
		Name:    args.Name,
		OwnerID: uint(args.OwnerID),
	}

	err = db.DB.Model(&p).Updates(updated).Error
	if err != nil {
		return nil, err
	}

	err = db.DB.First(&p, args.ID).Error
	if err != nil {
		return nil, err
	}

	return &p, nil
}

func (db *DB) deleteVm(userID, VmID uint) (*bool, error) {
	// make sure the record exist
	var p model.Vm
	err := db.DB.First(&p, VmID).Error
	if err != nil {
		return nil, err
	}

	// delete devices
	err = db.DB.Model(&p).Association("Dev").Clear().Error
	if err != nil {
		return nil, err
	}

	// delete Things
	err = db.DB.Model(&p).Association("Things").Clear().Error
	if err != nil {
		return nil, err
	}

	// delete record
	err = db.DB.Delete(&p).Error
	if err != nil {
		return nil, err
	}

	res := true

	return &res, err
}

func (db *DB) addVm(input vmInput, userid uint) (*model.Vm, error) {
	// get relationed devices from the DB and put them in the vm to be saved
	var devices []model.Hostdev

	if input.DevIDs != nil {
		err := db.DB.Where("id in (?)", *input.DevIDs).Find(&devices).Error
		if err != nil {
			return nil, err
		}
	}

	var things []model.Thing
	if input.ThingIDs != nil {
		err := db.DB.Where("id in (?)", *input.ThingIDs).Find(&things).Error
		if err != nil {
			return nil, err
		}
	}

	vm := model.Vm{
		Name:    input.Name,
		Base:    input.Base,
		Memory:  (int)(input.Memory),
		Vcpu:    (int)(input.Vcpu),
		OwnerID: userid,
		Dev:     devices,
		Things:  things,
	}

	err := db.DB.Create(&vm).Error
	if err != nil {
		return nil, err
	}

	return &vm, nil
}

// ###########################################################
// DEVICES:
// ###########################################################

//func (db *DB) getTagVms(t *model.Tag) ([]model.Vm, error) {
//	var p []model.Vm
//	err := db.DB.Model(t).Related(&p, "Vms").Error
//	if err != nil {
//		return nil, err
//	}
//
//	return p, nil
//}

func (db *DB) getDev(id uint) (*model.Hostdev, error) {
	var d model.Hostdev

	err := db.DB.First(&d, id).Error
	if err != nil {
		return nil, err
	}

	return &d, nil
}

// ###########################################################
// THINGS:
// ###########################################################
func (db *DB) getThing(id uint) (*model.Thing, error) {
	var t model.Thing

	err := db.DB.First(&t, id).Error
	if err != nil {
		return nil, err
	}

	return &t, nil
}

func (db *DB) AddThing(input thingInput) (*model.Thing, error) {
	thing := model.Thing{
		Name: input.Name,
		Info: input.Info,
	}

	err := db.DB.Create(&thing).Error
	if err != nil {
		return nil, err
	}

	return &thing, nil
}
