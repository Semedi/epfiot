package core

import (
	"math/rand"
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
	db.DropTableIfExists(&model.User{}, &model.Vm{}, &model.Tag{})
	db.AutoMigrate(&model.User{}, &model.Vm{}, &model.Tag{})

	// put all the users into the db
	for _, u := range users {
		if err := db.Create(&u).Error; err != nil {
			return nil, err
		}
	}

	var tg = []model.Tag{}
	for _, t := range tags {
		if err := db.Create(&t).Error; err != nil {
			return nil, err
		}

		tg = append(tg, t)
	}

	// put all the vms into the db
	for _, p := range vms {
		p.Tags = tg[:rand.Intn(5)]
		if err := db.Create(&p).Error; err != nil {
			return nil, err
		}
	}

	return &DB{db}, nil
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

func (db *DB) getVmTags(p *model.Vm) ([]model.Tag, error) {
	var t []model.Tag
	err := db.DB.Model(p).Related(&t, "Tags").Error
	if err != nil {
		return nil, err
	}

	return t, nil
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
	if args.TagIDs == nil {
		return nil, err
	}

	// if there are tags to be updated, go through that process
	var newTags []model.Tag
	if len(*args.TagIDs) > 0 {
		err = db.DB.Where("id in (?)", args.TagIDs).Find(&newTags).Error
		if err != nil {
			return nil, err
		}

		// replace the old tag set with the new one
		err = db.DB.Model(&p).Association("Tags").Replace(newTags).Error
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

	// delete tags
	err = db.DB.Model(&p).Association("Tags").Clear().Error
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
	// get the M2M relation tags from the DB and put them in the vm to be saved
	var t []model.Tag
	err := db.DB.Where("id in (?)", input.TagIDs).Find(&t).Error
	if err != nil {
		return nil, err
	}

	vm := model.Vm{
		Name:    input.Name,
        Base:    input.Base,
        Memory:  (int)(input.Memory),
        Vcpu:    (int)(input.Vcpu),
		OwnerID: userid,
		Tags:    t,
	}

	err = db.DB.Create(&vm).Error
	if err != nil {
		return nil, err
	}

	return &vm, nil
}

// ###########################################################
// TAGS:
// ###########################################################
func (db *DB) getTagVms(t *model.Tag) ([]model.Vm, error) {
	var p []model.Vm
	err := db.DB.Model(t).Related(&p, "Vms").Error
	if err != nil {
		return nil, err
	}

	return p, nil
}

func (db *DB) getTagBytTitle(title string) (*model.Tag, error) {
	var t model.Tag
	err := db.DB.Where("title = ?", title).First(&t).Error
	if err != nil {
		return nil, err
	}

	return &t, nil
}
