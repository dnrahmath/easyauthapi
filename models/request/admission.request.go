package request

import (
	pl "easyauthapi/models/payload"
	validation "github.com/go-ozzo/ozzo-validation"
)

/*====================================================================*/

type AdmissionPostReq pl.Admission

func (a AdmissionPostReq) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.UuIdUser, validation.Required),
		validation.Field(&a.Username, validation.Required),
		validation.Field(&a.Password, validation.Required),
		validation.Field(&a.Email, validation.Required),
		validation.Field(&a.PhoneNumber, validation.Required),
	)
}

/*====================================================================*/

type AdmissionPutReq pl.Admission

func (a AdmissionPutReq) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.UuIdUser, validation.Required),
		validation.Field(&a.Username, validation.Required),
		validation.Field(&a.Password, validation.Required),
		validation.Field(&a.Email, validation.Required),
		validation.Field(&a.EmailVerified, validation.Required),
		validation.Field(&a.PhoneNumber, validation.Required),
		validation.Field(&a.PhoneNumberVerified, validation.Required),
	)
}

/*====================================================================*/
