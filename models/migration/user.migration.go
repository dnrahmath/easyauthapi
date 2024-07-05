package migration

import (
	initi "easyauthapi/configs"
	pl "easyauthapi/models/payload"
)

//===================================================================

type UserMigration struct {
	Data       pl.User
	MigrateRes error
}

//===================================================================

func (model *UserMigration) CollectionName() string {
	return "users"
}

//==================================

func (model *UserMigration) Migrate() *UserMigration {
	err := initi.DB.Preload("roles").AutoMigrate(&pl.User{}, &pl.Role{})
	return &UserMigration{
		MigrateRes: err,
	}
}

//(&User{}, &Role{}, &UserRole{})
// MigrateRes: initi.DB.AutoMigrate(model.Data),
//====================================================================
