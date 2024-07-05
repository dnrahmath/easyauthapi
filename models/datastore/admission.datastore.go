package datastore

import (
	"time"

	pl "easyauthapi/models/payload"
	"easyauthapi/utils"
	"github.com/google/uuid"
)

//===================================================================

type AdmissionDatastore pl.Admission

//===================================================================

func (DataNew *AdmissionDatastore) Post() *pl.Admission {
	data := &pl.Admission{
		UuId:                DataNew.UuId,
		UuIdUser:            DataNew.UuIdUser,
		Username:            DataNew.Username,
		Password:            DataNew.Password,
		Email:               DataNew.Email,
		EmailVerified:       DataNew.EmailVerified,
		PhoneNumber:         DataNew.PhoneNumber,
		PhoneNumberVerified: DataNew.PhoneNumberVerified,
		UpdatedAt:           time.Now(),
		CreatedAt:           time.Now(),
	}
	return data
}

//==================================

func (DataPrev *AdmissionDatastore) Put(DataNew *pl.Admission) *pl.Admission {
	// Menggunakan kondisi terpisah untuk setiap field
	UuIdUser := utils.Ternary(DataNew.UuIdUser != uuid.Nil, DataNew.UuIdUser, DataPrev.UuIdUser).(uuid.UUID)
	Username := utils.Ternary(DataNew.Username != "", DataNew.Username, DataPrev.Username).(string)
	Password := utils.Ternary(DataNew.Password != "", DataNew.Password, DataPrev.Password).(string)
	Email := utils.Ternary(DataNew.Email != "", DataNew.Email, DataPrev.Email).(string)
	EmailVerified := utils.Ternary(!DataNew.EmailVerified, DataNew.EmailVerified, DataPrev.EmailVerified).(bool)
	PhoneNumber := utils.Ternary(DataNew.PhoneNumber != "", DataNew.PhoneNumber, DataPrev.PhoneNumber).(string)
	PhoneNumberVerified := utils.Ternary(!DataNew.PhoneNumberVerified, DataNew.PhoneNumberVerified, DataPrev.PhoneNumberVerified).(bool)

	// Membuat objek baru dengan nilai-nilai yang telah ditetapkan
	data := &pl.Admission{
		UuIdUser:            UuIdUser,
		Username:            Username,
		Password:            Password,
		Email:               Email,
		EmailVerified:       EmailVerified,
		PhoneNumber:         PhoneNumber,
		PhoneNumberVerified: PhoneNumberVerified,
		UpdatedAt:           DataPrev.CreatedAt,
		CreatedAt:           time.Now(),
	}
	return data
}

//===================================================================
