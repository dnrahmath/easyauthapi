package datastore

import (
	"time"

	pl "easyauthapi/models/payload"
	"easyauthapi/utils"
	"github.com/google/uuid"
)

//===================================================================

type UserDatastore pl.User

//===================================================================

const (
	RoleGuest    = "Guest"
	RoleUser     = "User"
	RoleMerchant = "Merchant"
	RoleAdmin    = "Admin"
)

var (
	RoleStructDefault    = []pl.Role{{Name: RoleUser}}
	RoleStructMerchant   = []pl.Role{{Name: RoleUser}, {Name: RoleMerchant}}
	RoleStructFirstAdmin = []pl.Role{{Name: RoleUser}, {Name: RoleAdmin}}
)

//==================================

func (DataNew *UserDatastore) Post() *pl.User {
	data := &pl.User{
		UuId:          DataNew.UuId,
		UuIdAddr:      DataNew.UuIdAddr,
		UuIdAdmission: DataNew.UuIdAdmission,
		Gender:        DataNew.Gender,
		Name:          DataNew.Name,
		Noid:          DataNew.Noid,
		Religion:      DataNew.Religion,
		Roles:         utils.Ternary(DataNew.Roles != nil, DataNew.Roles, RoleStructDefault).([]pl.Role),
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
	return data
}

//==================================

func (DataPrev *UserDatastore) Put(DataNew *pl.User) *pl.User {
	// Menggunakan kondisi terpisah untuk setiap field
	UuId := utils.Ternary(DataNew.UuId != uuid.Nil, DataNew.UuId, DataPrev.UuId).(uuid.UUID)
	UuIdAddr := utils.Ternary(DataNew.UuIdAddr != uuid.Nil, DataNew.UuIdAddr, DataPrev.UuIdAddr).(uuid.UUID)
	UuIdAdmission := utils.Ternary(DataNew.UuIdAdmission != uuid.Nil, DataNew.UuIdAdmission, DataPrev.UuIdAdmission).(uuid.UUID)
	Gender := utils.Ternary(DataNew.Gender != "", DataNew.Gender, DataPrev.Gender).(string)
	Name := utils.Ternary(DataNew.Name != "", DataNew.Name, DataPrev.Name).(string)
	Noid := utils.Ternary(DataNew.Noid != "", DataNew.Noid, DataPrev.Noid).(string)
	Religion := utils.Ternary(DataNew.Religion != "", DataNew.Religion, DataPrev.Religion).(string)
	Roles := utils.Ternary(DataNew.Roles != nil, DataNew.Roles, DataPrev.Roles).([]pl.Role)

	// Membuat objek baru dengan nilai-nilai yang telah ditetapkan
	data := &pl.User{
		UuId:          UuId,
		UuIdAddr:      UuIdAddr,
		UuIdAdmission: UuIdAdmission,
		Gender:        Gender,
		Name:          Name,
		Noid:          Noid,
		Religion:      Religion,
		Roles:         Roles,
		CreatedAt:     DataPrev.CreatedAt,
		UpdatedAt:     time.Now(),
	}
	return data
}

//===================================================================
