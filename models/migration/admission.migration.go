package migration

import (
	initi "easyauthapi/configs"
	pl "easyauthapi/models/payload"
)

//===================================================================

type AdmissionMigration struct {
	Data       pl.Admission
	MigrateRes error
}

//===================================================================

func (model *AdmissionMigration) CollectionName() string {
	return "*admissions"
}

//==================================

func (model *AdmissionMigration) Migrate() *AdmissionMigration {
	err := initi.DB.AutoMigrate(model.Data)
	return &AdmissionMigration{
		MigrateRes: err,
	}
}

//====================================================================
