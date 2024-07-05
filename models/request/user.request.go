package request

import (
	pl "easyauthapi/models/payload"
	"errors"

	validation "github.com/go-ozzo/ozzo-validation"
)

// var passwordRule = []validation.Rule{
// 	validation.Required,
// 	validation.Length(8, 32),
// 	validation.Match(regexp.MustCompile(`^\S+$`)).Error("cannot contain whitespaces"),
// }

// ===================================================================
// type UserPostReq pl.UserPostPut

// func (a UserPostReq) Validate() error {
// 	return validation.ValidateStruct(&a,
// 		validation.Field(&a.Username, validation.Required),
// 		validation.Field(&a.Password, passwordRule...),
// 		validation.Field(&a.Name, validation.Required, validation.Length(3, 64)),
// 		validation.Field(&a.Role, validation.By(func(value interface{}) error {
// 			if roles, ok := value.([]string); !ok || len(roles) > 0 {
// 				return errors.New("role should not be filled")
// 			}
// 			return nil
// 		})),
// 	)
// }

//===================================================================

type UserPutReq pl.UserPostPut

func (a UserPutReq) Validate() error {
	return validation.ValidateStruct(&a,
		// validation.Field(&a.Username, validation.Required),
		// validation.Field(&a.Password, passwordRule...),
		// validation.Field(&a.Email, validation.Required),
		// validation.Field(&a.PhoneNumber, validation.Required),
		validation.Field(&a.Gender, validation.Required),
		validation.Field(&a.Name, validation.Required, validation.Length(3, 64)),
		validation.Field(&a.Noid, validation.Required),
		validation.Field(&a.Religion, validation.Required),
		validation.Field(&a.Role, validation.By(func(value interface{}) error {
			if roles, ok := value.([]string); !ok || len(roles) > 0 {
				return errors.New("role should not be filled")
			}
			return nil
		})),
	)
}

//===================================================================

type UserPostReqByAdm pl.UserPostPut

func (a UserPostReqByAdm) Validate() error {
	return validation.ValidateStruct(&a,
		// validation.Field(&a.Username, validation.Required),
		// validation.Field(&a.Password, passwordRule...),
		// validation.Field(&a.Email, validation.Required),
		// validation.Field(&a.PhoneNumber, validation.Required),
		validation.Field(&a.Gender, validation.Required),
		validation.Field(&a.Name, validation.Required, validation.Length(3, 64)),
		validation.Field(&a.Noid, validation.Required),
		validation.Field(&a.Religion, validation.Required),
		validation.Field(&a.Role, validation.Required),
	)
}

//===================================================================

type UserPutReqByAdm pl.UserPostPut

func (a UserPutReqByAdm) Validate() error {
	return validation.ValidateStruct(&a,
		// validation.Field(&a.Username, validation.Required),
		// validation.Field(&a.Password, passwordRule...),
		// validation.Field(&a.Email, validation.Required),
		// validation.Field(&a.PhoneNumber, validation.Required),
		validation.Field(&a.Gender, validation.Required),
		validation.Field(&a.Name, validation.Required, validation.Length(3, 64)),
		validation.Field(&a.Noid, validation.Required),
		validation.Field(&a.Religion, validation.Required),
		validation.Field(&a.Role, validation.Required),
	)
}

//===================================================================

type LoginOrRegisReq pl.LoginAndRegis

func (a LoginOrRegisReq) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.Value, validation.Required),
		validation.Field(&a.Password, validation.Required),
	)
}

//===================================================================
