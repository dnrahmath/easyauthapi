package migration

import (
	initi "easyauthapi/configs"
	pl "easyauthapi/models/payload"
)

//===================================================================

type TokenMigration struct {
	Data       pl.Token
	MigrateRes error
}

//===================================================================

func (model *TokenMigration) CollectionName() string {
	return "tokens"
}

//==================================

func (model *TokenMigration) Migrate() *TokenMigration {
	err := initi.DB.AutoMigrate(model.Data)
	return &TokenMigration{
		MigrateRes: err,
	}
}

//====================================================================
